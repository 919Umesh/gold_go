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

	UpsertDailyPrice(price *models.StockPrice) error
	UpdateCompanyMarketCap(companyID string, marketCap decimal.Decimal) error
	GetPreviousDayClose(companyID string) (*models.StockPrice, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

func (r *repository) CreateCompany(company *models.Company) error {
	query := `INSERT INTO companies (symbol, name, sector, market_cap, description, founded_year, employees, total_shares, available_shares, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
	if company.TotalShares == 0 {
		company.TotalShares = 10000000
	}
	if company.AvailableShares == 0 {
		company.AvailableShares = company.TotalShares
	}
	return r.client.ExecuteInsert(query, company,
		company.Symbol, company.Name, company.Sector, company.MarketCap.InexactFloat64(),
		company.Description, company.FoundedYear, company.Employees,
		company.TotalShares, company.AvailableShares, company.IsActive)
}

func (r *repository) GetCompanyByID(companyID string) (*models.Company, error) {
	var company models.Company
	query := "SELECT * FROM companies WHERE id = $1"
	err := r.client.ExecuteQueryRow(query, &company, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}
	return &company, nil
}

func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	var company models.Company
	query := "SELECT * FROM companies WHERE symbol = $1"
	err := r.client.ExecuteQueryRow(query, &company, symbol)
	if err != nil {
		return nil, fmt.Errorf("company with symbol %s not found", symbol)
	}
	return &company, nil
}

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

func (r *repository) UpdateCompany(company *models.Company) error {
	query := `UPDATE companies SET name = $1, sector = $2, market_cap = $3, description = $4,
			  founded_year = $5, employees = $6, is_active = $7 WHERE id = $8 RETURNING *`
	return r.client.ExecuteUpdate(query, company,
		company.Name, company.Sector, company.MarketCap.InexactFloat64(), company.Description,
		company.FoundedYear, company.Employees, company.IsActive, company.ID)
}

func (r *repository) UpdateCompanyShares(companyID string, availableShares int64) error {
	query := `UPDATE companies SET available_shares = $1 WHERE id = $2`
	return r.client.ExecuteUpdate(query, nil, availableShares, companyID)
}

func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	query := `INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, timestamp, timeframe)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, price,
		price.CompanyID, price.OpenPrice.InexactFloat64(), price.HighPrice.InexactFloat64(),
		price.LowPrice.InexactFloat64(), price.ClosePrice.InexactFloat64(),
		price.Volume, price.Timestamp.Format(time.RFC3339), price.Timeframe)
}

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

func (r *repository) GetPriceHistory(companyID string, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error) {
	var prices []models.StockPrice
	var query string
	var err error

	if timeframe == "all" || timeframe == "" {
		// "all" returns daily (1D) candles — the primary meaningful data
		query = "SELECT * FROM stock_prices WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 ORDER BY timestamp DESC LIMIT $4"
		err = r.client.ExecuteQuery(query, &prices, companyID, "1D", from.Format(time.RFC3339), limit)
	} else {
		query = "SELECT * FROM stock_prices WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 ORDER BY timestamp DESC LIMIT $4"
		err = r.client.ExecuteQuery(query, &prices, companyID, timeframe, from.Format(time.RFC3339), limit)
	}

	if err != nil {
		return nil, err
	}

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

func (r *repository) getCompanyChangePct(companyID string) (decimal.Decimal, error) {
	latest, err := r.GetLatestPrice(companyID)
	if err != nil {
		return decimal.Zero, nil
	}
	return latest.CalculateChange(), nil
}

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

func (r *repository) CreateMarketEvent(event *models.MarketEvent) error {
	query := `INSERT INTO market_events (company_id, event_type, title, description, impact_percentage, event_date)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	return r.client.ExecuteInsert(query, event,
		event.CompanyID, event.EventType, event.Title, event.Description,
		event.ImpactPercentage, event.EventDate.Format(time.RFC3339))
}

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

func (r *repository) GetAllSectors() ([]string, error) {
	var companies []models.Company
	query := "SELECT sector FROM companies WHERE is_active = $1 ORDER BY sector ASC"
	err := r.client.ExecuteQuery(query, &companies, true)
	if err != nil {
		return nil, err
	}

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

func (r *repository) GetTotalCompaniesCount() (int, error) {
	var companies []models.Company
	query := "SELECT id FROM companies WHERE is_active = $1"
	err := r.client.ExecuteQuery(query, &companies, true)
	if err != nil {
		return 0, err
	}
	return len(companies), nil
}

func (r *repository) GetTotalCompaniesBySectorCount(sector string) (int, error) {
	var companies []models.Company
	query := "SELECT id FROM companies WHERE is_active = $1 AND sector = $2"
	err := r.client.ExecuteQuery(query, &companies, true, sector)
	if err != nil {
		return 0, err
	}
	return len(companies), nil
}

// UpsertDailyPrice creates or updates the 1D candle for a given day.
// If a 1D record already exists for that company on the same day, update it.
// Otherwise create a new one.
func (r *repository) UpsertDailyPrice(price *models.StockPrice) error {
	// Look for existing 1D record for today
	dayStart := time.Date(price.Timestamp.Year(), price.Timestamp.Month(), price.Timestamp.Day(), 0, 0, 0, 0, price.Timestamp.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 AND timestamp < $4 LIMIT 1"
	err := r.client.ExecuteQuery(query, &prices, price.CompanyID, "1D", dayStart.Format(time.RFC3339), dayEnd.Format(time.RFC3339))

	if err == nil && len(prices) > 0 {
		// Update existing daily candle
		existing := prices[0]
		updateQuery := `UPDATE stock_prices SET high_price = $1, low_price = $2, close_price = $3, volume = $4 WHERE id = $5`
		return r.client.ExecuteUpdate(updateQuery, nil,
			price.HighPrice.InexactFloat64(), price.LowPrice.InexactFloat64(),
			price.ClosePrice.InexactFloat64(), price.Volume, existing.ID)
	}

	// Create new daily candle
	return r.CreateStockPrice(price)
}

// UpdateCompanyMarketCap updates market cap for a company after price changes
func (r *repository) UpdateCompanyMarketCap(companyID string, marketCap decimal.Decimal) error {
	query := `UPDATE companies SET market_cap = $1 WHERE id = $2`
	return r.client.ExecuteUpdate(query, nil, marketCap.InexactFloat64(), companyID)
}

// GetPreviousDayClose returns the last closing price from before today
func (r *repository) GetPreviousDayClose(companyID string) (*models.StockPrice, error) {
	today := time.Now().Truncate(24 * time.Hour)
	var prices []models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT 1"
	err := r.client.ExecuteQuery(query, &prices, companyID, today.Format(time.RFC3339))
	if err != nil || len(prices) == 0 {
		return nil, fmt.Errorf("no previous day close found for company %s", companyID)
	}
	return &prices[0], nil
}
