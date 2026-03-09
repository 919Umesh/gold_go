package ipo

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

var (
	ErrIPONotOpen          = errors.New("IPO is not open for applications")
	ErrIPOWindowClosed     = errors.New("IPO application window has closed")
	ErrIPOWindowNotOpen    = errors.New("IPO application window has not opened yet")
	ErrAlreadyApplied      = errors.New("you have already applied for this IPO")
	ErrExceedsMaxShares    = errors.New("requested shares exceed maximum per applicant")
	ErrIPONotClosed        = errors.New("IPO window must be closed before allocation")
	ErrIPOAlreadyAllocated = errors.New("IPO has already been allocated")
)

type Service interface {
	// Admin operations
	CreateCompany(symbol, name, sector string, totalSupply int64) (*models.Company, error)
	LaunchIPO(companyID string, pricePerShare decimal.Decimal, totalShares, maxPerApplicant int64, openAt, closeAt time.Time) (*models.IPO, error)
	AllocateIPO(ipoID string) (*AllocationResult, error)

	// User operations
	ApplyForIPO(userID, ipoID string, sharesRequested int64) (*models.IPOApplication, error)
	ListIPOs(limit int) ([]models.IPO, error)
	GetIPO(ipoID string) (*models.IPO, error)
}

type AllocationResult struct {
	TotalApplicants   int             `json:"total_applicants"`
	TotalRequested    int64           `json:"total_requested"`
	TotalAllocated    int64           `json:"total_allocated"`
	AllocatedCount    int             `json:"allocated_count"`
	NotAllocatedCount int             `json:"not_allocated_count"`
	RefundedAmount    decimal.Decimal `json:"refunded_amount"`
}

type service struct {
	repo      Repository
	stockRepo stock.Repository
	walletSvc wallet.Service
}

func NewService(repo Repository, stockRepo stock.Repository, walletSvc wallet.Service) Service {
	return &service{
		repo:      repo,
		stockRepo: stockRepo,
		walletSvc: walletSvc,
	}
}

// ──────────────────── Admin: Create Company ────────────────────

func (s *service) CreateCompany(symbol, name, sector string, totalSupply int64) (*models.Company, error) {
	company := &models.Company{
		Symbol:       symbol,
		Name:         name,
		Sector:       sector,
		TotalSupply:  totalSupply,
		CurrentPrice: decimal.Zero,
		MarketCap:    decimal.Zero,
		IsActive:     true,
	}

	if err := s.stockRepo.CreateCompany(company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	return company, nil
}

// ──────────────────── Admin: Launch IPO ────────────────────

func (s *service) LaunchIPO(companyID string, pricePerShare decimal.Decimal, totalShares, maxPerApplicant int64, openAt, closeAt time.Time) (*models.IPO, error) {
	// Verify company exists
	_, err := s.stockRepo.GetCompanyByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	ipo := &models.IPO{
		CompanyID:       companyID,
		PricePerShare:   pricePerShare,
		TotalShares:     totalShares,
		AllocatedShares: 0,
		MaxPerApplicant: maxPerApplicant,
		OpenAt:          openAt,
		CloseAt:         closeAt,
		Status:          models.IPOStatusOpen,
	}

	if err := s.repo.CreateIPO(ipo); err != nil {
		return nil, fmt.Errorf("failed to create IPO: %w", err)
	}

	return ipo, nil
}

// ──────────────────── User: Apply for IPO ────────────────────

func (s *service) ApplyForIPO(userID, ipoID string, sharesRequested int64) (*models.IPOApplication, error) {
	ipo, err := s.repo.GetIPOByID(ipoID)
	if err != nil {
		return nil, err
	}

	// Check status
	if ipo.Status != models.IPOStatusOpen {
		return nil, ErrIPONotOpen
	}

	// Check window
	now := time.Now()
	if now.Before(ipo.OpenAt) {
		return nil, ErrIPOWindowNotOpen
	}
	if now.After(ipo.CloseAt) {
		return nil, ErrIPOWindowClosed
	}

	// Check max per applicant
	if sharesRequested > ipo.MaxPerApplicant {
		return nil, ErrExceedsMaxShares
	}

	// Check duplicate
	_, err = s.repo.GetApplicationByUserAndIPO(userID, ipoID)
	if err == nil {
		return nil, ErrAlreadyApplied
	}

	// Calculate cost and deduct from trading wallet
	cost := ipo.PricePerShare.Mul(decimal.NewFromInt(sharesRequested))

	// Lock funds in trading wallet
	if err := s.walletSvc.LockFunds(userID, cost); err != nil {
		return nil, fmt.Errorf("insufficient trading wallet balance: %w", err)
	}

	app := &models.IPOApplication{
		IPOID:           ipoID,
		UserID:          userID,
		SharesRequested: sharesRequested,
		AmountPaid:      cost,
		Status:          models.IPOAppStatusPending,
	}

	if err := s.repo.CreateApplication(app); err != nil {
		// Release locked funds on failure
		_ = s.walletSvc.ReleaseFunds(userID, cost)
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

// ──────────────────── Admin: Allocate IPO (Lottery) ────────────────────

func (s *service) AllocateIPO(ipoID string) (*AllocationResult, error) {
	ipo, err := s.repo.GetIPOByID(ipoID)
	if err != nil {
		return nil, err
	}

	if ipo.Status == models.IPOStatusAllocated {
		return nil, ErrIPOAlreadyAllocated
	}

	// Close the IPO first
	if ipo.Status == models.IPOStatusOpen {
		if err := s.repo.UpdateIPOStatus(ipoID, models.IPOStatusClosed); err != nil {
			return nil, err
		}
	}

	applications, err := s.repo.GetApplicationsByIPO(ipoID)
	if err != nil {
		return nil, err
	}

	if len(applications) == 0 {
		_ = s.repo.UpdateIPOAllocated(ipoID, 0)
		return &AllocationResult{}, nil
	}

	result := &AllocationResult{
		TotalApplicants: len(applications),
		RefundedAmount:  decimal.Zero,
	}

	for _, app := range applications {
		result.TotalRequested += app.SharesRequested
	}

	// Determine allocation using skip-value / lottery algorithm
	allocatedIndices := s.lotteryAllocation(applications, ipo.TotalShares, ipo.MaxPerApplicant)

	var totalAllocated int64

	for i, app := range applications {
		if shares, ok := allocatedIndices[i]; ok {
			// Allocated
			totalAllocated += shares

			// Deduct the locked funds (they were locked during application)
			cost := ipo.PricePerShare.Mul(decimal.NewFromInt(shares))
			if err := s.walletSvc.DeductLockedFunds(app.UserID, cost); err != nil {
				slog.Error("failed to deduct locked funds", "user_id", app.UserID, "error", err)
			}

			// If partially allocated, refund the difference
			if shares < app.SharesRequested {
				refundShares := app.SharesRequested - shares
				refundAmt := ipo.PricePerShare.Mul(decimal.NewFromInt(refundShares))
				_ = s.walletSvc.ReleaseFunds(app.UserID, refundAmt)
				result.RefundedAmount = result.RefundedAmount.Add(refundAmt)
				_ = s.repo.UpdateApplicationStatus(app.ID, models.IPOAppStatusAllocated, shares, refundAmt.String())
			} else {
				_ = s.repo.UpdateApplicationStatus(app.ID, models.IPOAppStatusAllocated, shares, "0")
			}

			// Credit shares to portfolio
			s.creditPortfolio(app.UserID, ipo.CompanyID, shares, ipo.PricePerShare)
			result.AllocatedCount++
		} else {
			// Not allocated — full refund
			_ = s.walletSvc.ReleaseFunds(app.UserID, app.AmountPaid)
			_ = s.repo.UpdateApplicationStatus(app.ID, models.IPOAppStatusRefunded, 0, app.AmountPaid.String())
			result.RefundedAmount = result.RefundedAmount.Add(app.AmountPaid)
			result.NotAllocatedCount++
		}
	}

	result.TotalAllocated = totalAllocated

	// Update IPO
	_ = s.repo.UpdateIPOAllocated(ipoID, totalAllocated)

	// Update company price to IPO price
	marketCap := ipo.PricePerShare.Mul(decimal.NewFromInt(ipo.TotalShares))
	_ = s.stockRepo.UpdateCompanyPrice(ipo.CompanyID, ipo.PricePerShare.String(), marketCap.String())

	// Create initial stock price
	_ = s.stockRepo.CreateStockPrice(&models.StockPrice{
		CompanyID:  ipo.CompanyID,
		OpenPrice:  ipo.PricePerShare,
		HighPrice:  ipo.PricePerShare,
		LowPrice:   ipo.PricePerShare,
		ClosePrice: ipo.PricePerShare,
		Volume:     totalAllocated,
		Timestamp:  time.Now(),
		Timeframe:  "1D",
	})

	return result, nil
}

// lotteryAllocation implements skip-value / lottery-based allocation
// Returns a map of application index -> shares allocated
func (s *service) lotteryAllocation(apps []models.IPOApplication, availableShares, maxPerApplicant int64) map[int]int64 {
	allocated := make(map[int]int64)

	if len(apps) == 0 || availableShares <= 0 {
		return allocated
	}

	// Calculate total requested
	var totalRequested int64
	for _, app := range apps {
		totalRequested += app.SharesRequested
	}

	// If under-subscribed, give everyone what they asked
	if totalRequested <= availableShares {
		for i, app := range apps {
			allocated[i] = app.SharesRequested
		}
		return allocated
	}

	// Oversubscribed — use skip-value lottery
	// Each applicant gets a "lot" (shares they requested, capped by max)
	n := len(apps)
	sharesPerLot := maxPerApplicant
	availableLots := availableShares / sharesPerLot
	if availableLots <= 0 {
		availableLots = 1
	}

	if int64(n) <= availableLots {
		// Enough lots for everyone, just cap at available
		remaining := availableShares
		for i, app := range apps {
			give := app.SharesRequested
			if give > remaining {
				give = remaining
			}
			if give > 0 {
				allocated[i] = give
				remaining -= give
			}
			if remaining <= 0 {
				break
			}
		}
		return allocated
	}

	// Skip-value selection
	// skipInterval = totalApplicants / availableLots
	skipInterval := float64(n) / float64(availableLots)

	// Random starting index
	startIndex := rand.Intn(n)

	remaining := availableShares
	selected := int64(0)

	for lot := int64(0); lot < availableLots && remaining > 0; lot++ {
		idx := (startIndex + int(float64(lot)*skipInterval)) % n

		if _, alreadySelected := allocated[idx]; alreadySelected {
			// Skip if already allocated (handle wraparound collisions)
			continue
		}

		give := apps[idx].SharesRequested
		if give > sharesPerLot {
			give = sharesPerLot
		}
		if give > remaining {
			give = remaining
		}

		allocated[idx] = give
		remaining -= give
		selected += give
	}

	return allocated
}

func (s *service) creditPortfolio(userID, companyID string, quantity int64, price decimal.Decimal) {
	item, err := s.repo.GetPortfolioItem(userID, companyID)
	if err != nil {
		// Create new portfolio item
		newItem := &models.Portfolio{
			UserID:      userID,
			CompanyID:   companyID,
			Quantity:    quantity,
			AvgBuyPrice: price,
		}
		if createErr := s.repo.CreatePortfolioItem(newItem); createErr != nil {
			slog.Error("failed to create portfolio item", "user_id", userID, "error", createErr)
		}
		return
	}

	// Update existing
	totalInvested := item.AvgBuyPrice.Mul(decimal.NewFromInt(item.Quantity))
	newInvested := price.Mul(decimal.NewFromInt(quantity))
	newQty := item.Quantity + quantity
	newAvg := totalInvested.Add(newInvested).Div(decimal.NewFromInt(newQty))

	item.Quantity = newQty
	item.AvgBuyPrice = newAvg

	if err := s.repo.UpdatePortfolioItem(item); err != nil {
		slog.Error("failed to update portfolio item", "user_id", userID, "error", err)
	}
}

// ──────────────────── Queries ────────────────────

func (s *service) ListIPOs(limit int) ([]models.IPO, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListIPOs(limit)
}

func (s *service) GetIPO(ipoID string) (*models.IPO, error) {
	return s.repo.GetIPOByID(ipoID)
}
