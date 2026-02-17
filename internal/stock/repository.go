package stock

import (
	"fmt"
	"sort"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

const (
	TableCompanies    = "companies"
	TableStockPrices  = "stock_prices"
	TableMarketEvents = "market_events"
)

type Repository interface {
	CreateCompany(company *models.Company) error
	GetCompanyByID(id string) (*models.Company, error)
	GetCompanyBySymbol(symbol string) (*models.Company, error)
	ListCompanies(limit, offset int) ([]models.Company, error)
	ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error)
	GetAllSectors() ([]string, error)
	SearchCompanies(query string, limit int) ([]models.Company, error)
	UpdateCompany(company *models.Company) error
	UpdateCompanyShares(companyID string, availableShares int64) error
	GetTotalCompaniesCount() (int, error)
	GetTotalCompaniesBySectorCount(sector string) (int, error)

	CreateStockPrice(price *models.StockPrice) error
	GetLatestPrice(companyID string) (*models.StockPrice, error)
	GetPriceHistory(companyID string, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error)
	GetPriceAtTime(companyID string, timestamp time.Time) (*models.StockPrice, error)
	UpdateStockPriceVolume(priceID string, volumeIncrease int64) error

	GetTopGainers(limit int) ([]models.Company, error)
	GetTopLosers(limit int) ([]models.Company, error)
	GetMostActive(limit int) ([]models.Company, error)

	CreateMarketEvent(event *models.MarketEvent) error
	GetUpcomingEvents(companyID string, limit int) ([]models.MarketEvent, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

// =============================================================================
// CreateCompany — INSERT INTO companies (symbol, name, sector, market_cap, description,
//
//	founded_year, employees, total_shares, available_shares, is_active)
//	VALUES (...) RETURNING *
//
// =============================================================================
func (r *repository) CreateCompany(company *models.Company) error {
	query := `INSERT INTO companies (symbol, name, sector, market_cap, description, founded_year, employees, total_shares, available_shares, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
	// If total_shares not set, default to 10 million
	if company.TotalShares == 0 {
		company.TotalShares = 10000000
	}
	// available_shares defaults to total_shares
	if company.AvailableShares == 0 {
		company.AvailableShares = company.TotalShares
	}
	return r.client.ExecuteInsert(query, company,
		company.Symbol, company.Name, company.Sector, company.MarketCap.InexactFloat64(),
		company.Description, company.FoundedYear, company.Employees,
		company.TotalShares, company.AvailableShares, company.IsActive)
}

// =============================================================================
// GetCompanyByID — SELECT * FROM companies WHERE id = $1
// =============================================================================
func (r *repository) GetCompanyByID(companyID string) (*models.Company, error) {
	var company models.Company
	query := "SELECT * FROM companies WHERE id = $1"
	err := r.client.ExecuteQueryRow(query, &company, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}
	return &company, nil
}

// =============================================================================
// GetCompanyBySymbol — SELECT * FROM companies WHERE symbol = $1
// =============================================================================
func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	var company models.Company
	query := "SELECT * FROM companies WHERE symbol = $1"
	err := r.client.ExecuteQueryRow(query, &company, symbol)
	if err != nil {
		return nil, fmt.Errorf("company with symbol %s not found", symbol)
	}
	return &company, nil
}

// =============================================================================
// ListCompanies — SELECT * FROM companies WHERE is_active = $1
//
//	ORDER BY market_cap DESC LIMIT $2 OFFSET $3
//
// =============================================================================
func (r *repository) ListCompanies(limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := "SELECT * FROM companies WHERE is_active = $1 ORDER BY market_cap DESC LIMIT $2 OFFSET $3"
	err := r.client.ExecuteQuery(query, &companies, true, limit, offset)
	if err != nil {
		return nil, err
	}
	if companies == nil {
		companies = []models.Company{}
	}
	return companies, nil
}

// =============================================================================
// ListCompaniesBySector — SELECT * FROM companies WHERE is_active = $1 AND sector = $2
//
//	ORDER BY market_cap DESC LIMIT $3 OFFSET $4
//
// =============================================================================
func (r *repository) ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := "SELECT * FROM companies WHERE is_active = $1 AND sector = $2 ORDER BY market_cap DESC LIMIT $3 OFFSET $4"
	err := r.client.ExecuteQuery(query, &companies, true, sector, limit, offset)
	if err != nil {
		return nil, err
	}
	if companies == nil {
		companies = []models.Company{}
	}
	return companies, nil
}

// =============================================================================
// SearchCompanies — SELECT * FROM companies WHERE is_active = $1
//
//	AND (symbol ILIKE $2 OR name ILIKE $3)
//	ORDER BY market_cap DESC LIMIT $4
//
// =============================================================================
func (r *repository) SearchCompanies(q string, limit int) ([]models.Company, error) {
	var companies []models.Company
	pattern := fmt.Sprintf("*%s*", q)
	query := "SELECT * FROM companies WHERE is_active = $1 AND (symbol ILIKE $2 OR name ILIKE $3) ORDER BY market_cap DESC LIMIT $4"
	err := r.client.ExecuteQuery(query, &companies, true, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	if companies == nil {
		companies = []models.Company{}
	}
	return companies, nil
}

// =============================================================================
// UpdateCompany — UPDATE companies SET name=$1, sector=$2, market_cap=$3,
//
//	description=$4, founded_year=$5, employees=$6, is_active=$7
//	WHERE id = $8 RETURNING *
//
// =============================================================================
func (r *repository) UpdateCompany(company *models.Company) error {
	query := `UPDATE companies SET name = $1, sector = $2, market_cap = $3, description = $4,
			  founded_year = $5, employees = $6, is_active = $7 WHERE id = $8 RETURNING *`
	return r.client.ExecuteUpdate(query, company,
		company.Name, company.Sector, company.MarketCap.InexactFloat64(), company.Description,
		company.FoundedYear, company.Employees, company.IsActive, company.ID)
}

// =============================================================================
// UpdateCompanyShares — UPDATE companies SET available_shares = $1 WHERE id = $2
// Used when buying/selling to track available shares in the market
// =============================================================================
func (r *repository) UpdateCompanyShares(companyID string, availableShares int64) error {
	query := `UPDATE companies SET available_shares = $1 WHERE id = $2`
	return r.client.ExecuteUpdate(query, nil, availableShares, companyID)
}

// =============================================================================
// CreateStockPrice — INSERT INTO stock_prices (company_id, open_price, high_price,
//
//	low_price, close_price, volume, timestamp, timeframe)
//	VALUES (...) RETURNING *
//
// =============================================================================
func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	query := `INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, timestamp, timeframe)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, price,
		price.CompanyID, price.OpenPrice.InexactFloat64(), price.HighPrice.InexactFloat64(),
		price.LowPrice.InexactFloat64(), price.ClosePrice.InexactFloat64(),
		price.Volume, price.Timestamp.Format(time.RFC3339), price.Timeframe)
}

// =============================================================================
// GetLatestPrice — SELECT * FROM stock_prices WHERE company_id = $1
//
//	ORDER BY timestamp DESC LIMIT 1
//
// =============================================================================
func (r *repository) GetLatestPrice(companyID string) (*models.StockPrice, error) {
	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 ORDER BY timestamp DESC LIMIT 1"
	err := r.client.ExecuteQuery(query, &prices, companyID)
	if err != nil {
		return nil, fmt.Errorf("no price data found for company %s", companyID)
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no price data found for company %s", companyID)
	}
	return &prices[0], nil
}

// =============================================================================
// UpdateStockPriceVolume — UPDATE stock_prices SET volume = volume + $1 WHERE id = $2
// Increases volume when a trade occurs (buy or sell)
// =============================================================================
func (r *repository) UpdateStockPriceVolume(priceID string, volumeIncrease int64) error {
	// First get current volume
	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE id = $1"
	err := r.client.ExecuteQuery(query, &prices, priceID)
	if err != nil || len(prices) == 0 {
		return fmt.Errorf("price record not found")
	}
	newVolume := prices[0].Volume + volumeIncrease
	updateQuery := `UPDATE stock_prices SET volume = $1 WHERE id = $2`
	return r.client.ExecuteUpdate(updateQuery, nil, newVolume, priceID)
}

// =============================================================================
// GetPriceHistory — SELECT * FROM stock_prices WHERE company_id = $1
//
//	AND timeframe = $2 AND timestamp >= $3
//	ORDER BY timestamp DESC LIMIT $4
func (r *repository) GetPriceHistory(companyID string, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error) {
	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 ORDER BY timestamp DESC LIMIT $4"
	err := r.client.ExecuteQuery(query, &prices,
		companyID, timeframe, from.Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}

	// Additional client-side filter for upper bound (PostgREST doesn't support
	// two filters on the same column easily via the SQL parser approach)
	filtered := make([]models.StockPrice, 0, len(prices))
	for _, p := range prices {
		if !p.Timestamp.After(to) {
			filtered = append(filtered, p)
		}
	}

	if filtered == nil {
		filtered = []models.StockPrice{}
	}
	return filtered, nil
}

// =============================================================================
// GetPriceAtTime — SELECT * FROM stock_prices WHERE company_id = $1
//
//	AND timestamp <= $2 ORDER BY timestamp DESC LIMIT 1
//
// =============================================================================
func (r *repository) GetPriceAtTime(companyID string, timestamp time.Time) (*models.StockPrice, error) {
	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 AND timestamp <= $2 ORDER BY timestamp DESC LIMIT 1"
	err := r.client.ExecuteQuery(query, &prices,
		companyID, timestamp.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("no price data found for company %s", companyID)
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no price data found for company %s", companyID)
	}
	return &prices[0], nil
}

// Helper: get price change percentage for a company
func (r *repository) getCompanyChangePct(companyID string) (decimal.Decimal, error) {
	latest, err := r.GetLatestPrice(companyID)
	if err != nil {
		return decimal.Zero, nil
	}
	return latest.CalculateChange(), nil
}

// =============================================================================
// GetTopGainers — SELECT companies sorted by highest price change
// =============================================================================
func (r *repository) GetTopGainers(limit int) ([]models.Company, error) {
	companies, err := r.ListCompanies(100, 0)
	if err != nil {
		return nil, err
	}

	type enrichedCompany struct {
		company models.Company
		change  decimal.Decimal
	}

	var enriched []enrichedCompany
	for _, c := range companies {
		change, _ := r.getCompanyChangePct(c.ID)
		enriched = append(enriched, enrichedCompany{c, change})
	}

	sort.Slice(enriched, func(i, j int) bool {
		return enriched[i].change.GreaterThan(enriched[j].change)
	})

	result := make([]models.Company, 0, len(companies))
	for i, ec := range enriched {
		if i >= limit {
			break
		}
		result = append(result, ec.company)
	}
	return result, nil
}

// =============================================================================
// GetTopLosers — SELECT companies sorted by lowest price change
// =============================================================================
func (r *repository) GetTopLosers(limit int) ([]models.Company, error) {
	companies, err := r.ListCompanies(100, 0)
	if err != nil {
		return nil, err
	}

	type enrichedCompany struct {
		company models.Company
		change  decimal.Decimal
	}

	var enriched []enrichedCompany
	for _, c := range companies {
		change, _ := r.getCompanyChangePct(c.ID)
		enriched = append(enriched, enrichedCompany{c, change})
	}

	sort.Slice(enriched, func(i, j int) bool {
		return enriched[i].change.LessThan(enriched[j].change)
	})

	result := make([]models.Company, 0, len(companies))
	for i, ec := range enriched {
		if i >= limit {
			break
		}
		result = append(result, ec.company)
	}
	return result, nil
}

// =============================================================================
// GetMostActive — SELECT companies sorted by highest trading volume
// =============================================================================
func (r *repository) GetMostActive(limit int) ([]models.Company, error) {
	companies, err := r.ListCompanies(100, 0)
	if err != nil {
		return nil, err
	}

	type enrichedCompany struct {
		company models.Company
		volume  int64
	}

	var enriched []enrichedCompany
	for _, c := range companies {
		latest, err := r.GetLatestPrice(c.ID)
		vol := int64(0)
		if err == nil {
			vol = latest.Volume
		}
		enriched = append(enriched, enrichedCompany{c, vol})
	}

	sort.Slice(enriched, func(i, j int) bool {
		return enriched[i].volume > enriched[j].volume
	})

	result := make([]models.Company, 0, len(companies))
	for i, ec := range enriched {
		if i >= limit {
			break
		}
		result = append(result, ec.company)
	}
	return result, nil
}

// =============================================================================
// CreateMarketEvent — INSERT INTO market_events (company_id, event_type, title,
//
//	description, impact_percentage, event_date) VALUES (...)
//	RETURNING *
//
// =============================================================================
func (r *repository) CreateMarketEvent(event *models.MarketEvent) error {
	query := `INSERT INTO market_events (company_id, event_type, title, description, impact_percentage, event_date)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	return r.client.ExecuteInsert(query, event,
		event.CompanyID, event.EventType, event.Title, event.Description,
		event.ImpactPercentage, event.EventDate.Format(time.RFC3339))
}

// =============================================================================
// GetUpcomingEvents — SELECT * FROM market_events WHERE company_id = $1
//
//	AND event_date >= $2 ORDER BY event_date ASC LIMIT $3
//
// =============================================================================
func (r *repository) GetUpcomingEvents(companyID string, limit int) ([]models.MarketEvent, error) {
	now := time.Now().Format(time.RFC3339)

	var events []models.MarketEvent
	query := "SELECT * FROM market_events WHERE company_id = $1 AND event_date >= $2 ORDER BY event_date ASC LIMIT $3"
	err := r.client.ExecuteQuery(query, &events, companyID, now, limit)
	if err != nil {
		return nil, err
	}
	if events == nil {
		events = []models.MarketEvent{}
	}
	return events, nil
}

// =============================================================================
// GetAllSectors — SELECT sector FROM companies WHERE is_active = $1
//
//	ORDER BY sector ASC
//
// =============================================================================
func (r *repository) GetAllSectors() ([]string, error) {
	var companies []models.Company
	query := "SELECT sector FROM companies WHERE is_active = $1 ORDER BY sector ASC"
	err := r.client.ExecuteQuery(query, &companies, true)
	if err != nil {
		return nil, err
	}

	// Get unique sectors
	sectorMap := make(map[string]bool)
	for _, c := range companies {
		if c.Sector != "" {
			sectorMap[c.Sector] = true
		}
	}

	sectors := make([]string, 0, len(sectorMap))
	for sector := range sectorMap {
		sectors = append(sectors, sector)
	}
	sort.Strings(sectors)
	return sectors, nil
}

// =============================================================================
// GetTotalCompaniesCount — SELECT id FROM companies WHERE is_active = $1
// =============================================================================
func (r *repository) GetTotalCompaniesCount() (int, error) {
	var companies []models.Company
	query := "SELECT id FROM companies WHERE is_active = $1"
	err := r.client.ExecuteQuery(query, &companies, true)
	if err != nil {
		return 0, err
	}
	return len(companies), nil
}

// =============================================================================
// GetTotalCompaniesBySectorCount — SELECT id FROM companies
//
//	WHERE is_active = $1 AND sector = $2
//
// =============================================================================
func (r *repository) GetTotalCompaniesBySectorCount(sector string) (int, error) {
	var companies []models.Company
	query := "SELECT id FROM companies WHERE is_active = $1 AND sector = $2"
	err := r.client.ExecuteQuery(query, &companies, true, sector)
	if err != nil {
		return 0, err
	}
	return len(companies), nil
}
