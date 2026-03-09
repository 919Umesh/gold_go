package ipo

import (
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

type Repository interface {
	// IPO
	CreateIPO(ipo *models.IPO) error
	GetIPOByID(id string) (*models.IPO, error)
	ListIPOs(limit int) ([]models.IPO, error)
	UpdateIPOStatus(ipoID, status string) error
	UpdateIPOAllocated(ipoID string, allocatedShares int64) error

	// Applications
	CreateApplication(app *models.IPOApplication) error
	GetApplicationsByIPO(ipoID string) ([]models.IPOApplication, error)
	GetApplicationByUserAndIPO(userID, ipoID string) (*models.IPOApplication, error)
	UpdateApplicationStatus(appID, status string, sharesAllocated int64, amountRefunded string) error

	// Portfolio
	GetPortfolioItem(userID, companyID string) (*models.Portfolio, error)
	CreatePortfolioItem(item *models.Portfolio) error
	UpdatePortfolioItem(item *models.Portfolio) error
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

// ──────────────────── IPO ────────────────────

func (r *repository) CreateIPO(ipo *models.IPO) error {
	query := `INSERT INTO ipos (company_id, price_per_share, total_shares, max_per_applicant, open_at, close_at, status)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	return r.client.ExecuteInsert(query, ipo,
		ipo.CompanyID, ipo.PricePerShare.String(), ipo.TotalShares, ipo.MaxPerApplicant,
		ipo.OpenAt.Format(time.RFC3339), ipo.CloseAt.Format(time.RFC3339), ipo.Status)
}

func (r *repository) GetIPOByID(id string) (*models.IPO, error) {
	var ipo models.IPO
	err := r.client.ExecuteQueryRow("SELECT * FROM ipos WHERE id = $1", &ipo, id)
	if err != nil {
		return nil, fmt.Errorf("IPO not found: %w", err)
	}
	return &ipo, nil
}

func (r *repository) ListIPOs(limit int) ([]models.IPO, error) {
	var ipos []models.IPO
	query := "SELECT * FROM ipos ORDER BY created_at DESC LIMIT $1"
	err := r.client.ExecuteQuery(query, &ipos, limit)
	if err != nil {
		return nil, err
	}
	if ipos == nil {
		ipos = []models.IPO{}
	}
	return ipos, nil
}

func (r *repository) UpdateIPOStatus(ipoID, status string) error {
	query := `UPDATE ipos SET status = $1 WHERE id = $2 RETURNING *`
	var ipo models.IPO
	return r.client.ExecuteUpdate(query, &ipo, status, ipoID)
}

func (r *repository) UpdateIPOAllocated(ipoID string, allocatedShares int64) error {
	query := `UPDATE ipos SET allocated_shares = $1, status = $2 WHERE id = $3 RETURNING *`
	var ipo models.IPO
	return r.client.ExecuteUpdate(query, &ipo, allocatedShares, models.IPOStatusAllocated, ipoID)
}

// ──────────────────── Applications ────────────────────

func (r *repository) CreateApplication(app *models.IPOApplication) error {
	query := `INSERT INTO ipo_applications (ipo_id, user_id, shares_requested, amount_paid, status)
			  VALUES ($1, $2, $3, $4, $5) RETURNING *`
	return r.client.ExecuteInsert(query, app,
		app.IPOID, app.UserID, app.SharesRequested, app.AmountPaid.String(), app.Status)
}

func (r *repository) GetApplicationsByIPO(ipoID string) ([]models.IPOApplication, error) {
	var apps []models.IPOApplication
	query := "SELECT * FROM ipo_applications WHERE ipo_id = $1 ORDER BY created_at"
	err := r.client.ExecuteQuery(query, &apps, ipoID)
	if err != nil {
		return nil, err
	}
	if apps == nil {
		apps = []models.IPOApplication{}
	}
	return apps, nil
}

func (r *repository) GetApplicationByUserAndIPO(userID, ipoID string) (*models.IPOApplication, error) {
	var app models.IPOApplication
	query := "SELECT * FROM ipo_applications WHERE user_id = $1 AND ipo_id = $2"
	err := r.client.ExecuteQueryRow(query, &app, userID, ipoID)
	if err != nil {
		return nil, fmt.Errorf("application not found")
	}
	return &app, nil
}

func (r *repository) UpdateApplicationStatus(appID, status string, sharesAllocated int64, amountRefunded string) error {
	query := `UPDATE ipo_applications SET status = $1, shares_allocated = $2, amount_refunded = $3 WHERE id = $4 RETURNING *`
	var app models.IPOApplication
	return r.client.ExecuteUpdate(query, &app, status, sharesAllocated, amountRefunded, appID)
}

// ──────────────────── Portfolio ────────────────────

func (r *repository) GetPortfolioItem(userID, companyID string) (*models.Portfolio, error) {
	var item models.Portfolio
	query := "SELECT * FROM portfolios WHERE user_id = $1 AND company_id = $2"
	err := r.client.ExecuteQueryRow(query, &item, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("portfolio item not found")
	}
	return &item, nil
}

func (r *repository) CreatePortfolioItem(item *models.Portfolio) error {
	query := `INSERT INTO portfolios (user_id, company_id, quantity, avg_buy_price)
			  VALUES ($1, $2, $3, $4) RETURNING *`
	return r.client.ExecuteInsert(query, item,
		item.UserID, item.CompanyID, item.Quantity, item.AvgBuyPrice.String())
}

func (r *repository) UpdatePortfolioItem(item *models.Portfolio) error {
	query := `UPDATE portfolios SET quantity = $1, avg_buy_price = $2 WHERE id = $3 RETURNING *`
	return r.client.ExecuteUpdate(query, item,
		item.Quantity, item.AvgBuyPrice.String(), item.ID)
}
