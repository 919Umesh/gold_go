package stock

import (
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

type Repository interface {
	// Companies
	CreateCompany(company *models.Company) error
	GetCompanyByID(id string) (*models.Company, error)
	GetCompanyBySymbol(symbol string) (*models.Company, error)
	ListCompanies(limit, offset int) ([]models.Company, error)
	ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error)
	UpdateCompanyPrice(companyID string, price string, marketCap string) error

	// Stock Prices (OHLCV)
	CreateStockPrice(price *models.StockPrice) error
	GetLatestPrice(companyID string) (*models.StockPrice, error)
	GetPriceHistory(companyID, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error)
	UpsertDailyPrice(price *models.StockPrice) error
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

// ──────────────────── Companies ────────────────────

func (r *repository) CreateCompany(company *models.Company) error {
	query := `INSERT INTO companies (symbol, name, sector, description, total_supply, shares_outstanding,
			  current_price, market_cap, eps, pe_ratio, book_value, pbv,
			  week_52_high, week_52_low, avg_120_day, yield_1_year, listed_date, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING *`
	var listedDate interface{}
	if company.ListedDate == "" {
		listedDate = nil
	} else {
		listedDate = company.ListedDate
	}

	return r.client.ExecuteInsert(query, company,
		company.Symbol, company.Name, company.Sector, company.Description,
		company.TotalSupply, company.SharesOutstanding.String(),
		company.CurrentPrice.String(), company.MarketCap.String(),
		company.EPS.String(), company.PERatio.String(),
		company.BookValue.String(), company.PBV.String(),
		company.Week52High.String(), company.Week52Low.String(),
		company.Avg120Day.String(), company.Yield1Year.String(),
		listedDate, company.IsActive)
}

func (r *repository) GetCompanyByID(id string) (*models.Company, error) {
	var c models.Company
	err := r.client.ExecuteQueryRow("SELECT * FROM companies WHERE id = $1", &c, id)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &c, nil
}

func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	var c models.Company
	err := r.client.ExecuteQueryRow("SELECT * FROM companies WHERE symbol = $1", &c, symbol)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &c, nil
}

func (r *repository) ListCompanies(limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := "SELECT * FROM companies WHERE is_active = $1 ORDER BY symbol LIMIT $2 OFFSET $3"
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
	query := "SELECT * FROM companies WHERE sector = $1 AND is_active = $2 ORDER BY symbol LIMIT $3 OFFSET $4"
	err := r.client.ExecuteQuery(query, &companies, sector, true, limit, offset)
	if err != nil {
		return nil, err
	}
	if companies == nil {
		companies = []models.Company{}
	}
	return companies, nil
}

func (r *repository) UpdateCompanyPrice(companyID string, price string, marketCap string) error {
	query := `UPDATE companies SET current_price = $1, market_cap = $2 WHERE id = $3 RETURNING *`
	var c models.Company
	return r.client.ExecuteUpdate(query, &c, price, marketCap, companyID)
}

// ──────────────────── Stock Prices ────────────────────

func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	query := `INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, turnover, change_percent, timestamp, timeframe)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
	return r.client.ExecuteInsert(query, price,
		price.CompanyID, price.OpenPrice.String(), price.HighPrice.String(),
		price.LowPrice.String(), price.ClosePrice.String(), price.Volume,
		price.Turnover.String(), price.ChangePercent.String(),
		price.Timestamp.Format(time.RFC3339), price.Timeframe)
}

func (r *repository) GetLatestPrice(companyID string) (*models.StockPrice, error) {
	var price models.StockPrice
	query := "SELECT * FROM stock_prices WHERE company_id = $1 ORDER BY timestamp DESC LIMIT $2"
	err := r.client.ExecuteQueryRow(query, &price, companyID, 1)
	if err != nil {
		return nil, fmt.Errorf("no price found: %w", err)
	}
	return &price, nil
}

func (r *repository) GetPriceHistory(companyID, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error) {
	var prices []models.StockPrice
	query := `SELECT * FROM stock_prices 
			  WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 AND timestamp <= $4 
			  ORDER BY timestamp DESC LIMIT $5`
	err := r.client.ExecuteQuery(query, &prices, companyID, timeframe,
		from.Format(time.RFC3339), to.Format(time.RFC3339), limit)
	if err != nil {
		return nil, err
	}
	if prices == nil {
		prices = []models.StockPrice{}
	}
	return prices, nil
}

func (r *repository) UpsertDailyPrice(price *models.StockPrice) error {
	existing, err := r.getDailyPrice(price.CompanyID, price.Timestamp)
	if err != nil {
		return r.CreateStockPrice(price)
	}

	query := `UPDATE stock_prices SET high_price = $1, low_price = $2, close_price = $3, volume = $4, turnover = $5, change_percent = $6
			  WHERE id = $7 RETURNING *`
	return r.client.ExecuteUpdate(query, existing,
		price.HighPrice.String(), price.LowPrice.String(),
		price.ClosePrice.String(), price.Volume,
		price.Turnover.String(), price.ChangePercent.String(), existing.ID)
}

func (r *repository) getDailyPrice(companyID string, day time.Time) (*models.StockPrice, error) {
	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var price models.StockPrice
	query := `SELECT * FROM stock_prices 
			  WHERE company_id = $1 AND timeframe = $2 AND timestamp >= $3 AND timestamp < $4 
			  ORDER BY timestamp DESC LIMIT $5`
	err := r.client.ExecuteQueryRow(query, &price, companyID, "1D",
		dayStart.Format(time.RFC3339), dayEnd.Format(time.RFC3339), 1)
	if err != nil {
		return nil, err
	}
	return &price, nil
}
