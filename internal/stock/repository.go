package stock

import (
	"fmt"
	"time"

	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Company operations
	CreateCompany(company *models.Company) error
	GetCompanyByID(id uint) (*models.Company, error)
	GetCompanyBySymbol(symbol string) (*models.Company, error)
	ListCompanies(limit, offset int) ([]models.Company, error)
	ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error)
	SearchCompanies(query string, limit int) ([]models.Company, error)
	UpdateCompany(company *models.Company) error

	// Stock price operations
	CreateStockPrice(price *models.StockPrice) error
	GetLatestPrice(companyID uint) (*models.StockPrice, error)
	GetPriceHistory(companyID uint, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error)
	GetPriceAtTime(companyID uint, timestamp time.Time) (*models.StockPrice, error)

	// Market overview
	GetTopGainers(limit int) ([]models.Company, error)
	GetTopLosers(limit int) ([]models.Company, error)
	GetMostActive(limit int) ([]models.Company, error)

	// Market events
	CreateMarketEvent(event *models.MarketEvent) error
	GetUpcomingEvents(companyID uint, limit int) ([]models.MarketEvent, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Company operations

func (r *repository) CreateCompany(company *models.Company) error {
	query := `
		INSERT INTO companies (symbol, name, sector, market_cap, description, founded_year, employees, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	return r.db.Raw(query, company.Symbol, company.Name, company.Sector, company.MarketCap,
		company.Description, company.FoundedYear, company.Employees, company.IsActive,
		time.Now(), time.Now()).Scan(&company.ID).Error
}

func (r *repository) GetCompanyByID(id uint) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&company).Error
	if err != nil {
		return nil, err
	}
	if company.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &company, nil
}

func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE symbol = ? LIMIT 1`
	err := r.db.Raw(query, symbol).Scan(&company).Error
	if err != nil {
		return nil, err
	}
	if company.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &company, nil
}

func (r *repository) ListCompanies(limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := `SELECT * FROM companies WHERE is_active = true ORDER BY market_cap DESC LIMIT ? OFFSET ?`
	err := r.db.Raw(query, limit, offset).Scan(&companies).Error
	return companies, err
}

func (r *repository) ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := `SELECT * FROM companies WHERE sector = ? AND is_active = true ORDER BY market_cap DESC LIMIT ? OFFSET ?`
	err := r.db.Raw(query, sector, limit, offset).Scan(&companies).Error
	return companies, err
}

func (r *repository) SearchCompanies(query string, limit int) ([]models.Company, error) {
	var companies []models.Company
	searchQuery := `
		SELECT * FROM companies 
		WHERE is_active = true AND (
			LOWER(name) LIKE LOWER(?) OR 
			LOWER(symbol) LIKE LOWER(?) OR 
			LOWER(sector) LIKE LOWER(?)
		)
		ORDER BY market_cap DESC
		LIMIT ?
	`
	searchPattern := "%" + query + "%"
	err := r.db.Raw(searchQuery, searchPattern, searchPattern, searchPattern, limit).Scan(&companies).Error
	return companies, err
}

func (r *repository) UpdateCompany(company *models.Company) error {
	query := `
		UPDATE companies 
		SET name = ?, sector = ?, market_cap = ?, description = ?, 
		    founded_year = ?, employees = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`
	return r.db.Exec(query, company.Name, company.Sector, company.MarketCap, company.Description,
		company.FoundedYear, company.Employees, company.IsActive, time.Now(), company.ID).Error
}

// Stock price operations

func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	query := `
		INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, timestamp, timeframe, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	return r.db.Exec(query, price.CompanyID, price.OpenPrice, price.HighPrice, price.LowPrice,
		price.ClosePrice, price.Volume, price.Timestamp, price.Timeframe, time.Now()).Error
}

func (r *repository) GetLatestPrice(companyID uint) (*models.StockPrice, error) {
	var price models.StockPrice
	query := `
		SELECT * FROM stock_prices 
		WHERE company_id = ? AND timeframe = '1d'
		ORDER BY timestamp DESC 
		LIMIT 1
	`
	err := r.db.Raw(query, companyID).Scan(&price).Error
	if err != nil {
		return nil, err
	}
	if price.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &price, nil
}

func (r *repository) GetPriceHistory(companyID uint, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error) {
	var prices []models.StockPrice
	query := `
		SELECT * FROM stock_prices 
		WHERE company_id = ? AND timeframe = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
		LIMIT ?
	`
	err := r.db.Raw(query, companyID, timeframe, from, to, limit).Scan(&prices).Error
	return prices, err
}

func (r *repository) GetPriceAtTime(companyID uint, timestamp time.Time) (*models.StockPrice, error) {
	var price models.StockPrice
	query := `
		SELECT * FROM stock_prices 
		WHERE company_id = ? AND timestamp <= ?
		ORDER BY timestamp DESC
		LIMIT 1
	`
	err := r.db.Raw(query, companyID, timestamp).Scan(&price).Error
	if err != nil {
		return nil, err
	}
	if price.ID == 0 {
		return nil, fmt.Errorf("no price data found")
	}
	return &price, nil
}

// Market overview

func (r *repository) GetTopGainers(limit int) ([]models.Company, error) {
	var companies []models.Company
	query := `
		SELECT c.*, 
		       (sp_latest.close_price - sp_prev.close_price) / sp_prev.close_price * 100 as change_pct
		FROM companies c
		INNER JOIN LATERAL (
			SELECT close_price FROM stock_prices 
			WHERE company_id = c.id AND timeframe = '1d'
			ORDER BY timestamp DESC LIMIT 1
		) sp_latest ON true
		INNER JOIN LATERAL (
			SELECT close_price FROM stock_prices 
			WHERE company_id = c.id AND timeframe = '1d'
			ORDER BY timestamp DESC LIMIT 1 OFFSET 1
		) sp_prev ON true
		WHERE c.is_active = true
		ORDER BY change_pct DESC
		LIMIT ?
	`
	err := r.db.Raw(query, limit).Scan(&companies).Error
	return companies, err
}

func (r *repository) GetTopLosers(limit int) ([]models.Company, error) {
	var companies []models.Company
	query := `
		SELECT c.*, 
		       (sp_latest.close_price - sp_prev.close_price) / sp_prev.close_price * 100 as change_pct
		FROM companies c
		INNER JOIN LATERAL (
			SELECT close_price FROM stock_prices 
			WHERE company_id = c.id AND timeframe = '1d'
			ORDER BY timestamp DESC LIMIT 1
		) sp_latest ON true
		INNER JOIN LATERAL (
			SELECT close_price FROM stock_prices 
			WHERE company_id = c.id AND timeframe = '1d'
			ORDER BY timestamp DESC LIMIT 1 OFFSET 1
		) sp_prev ON true
		WHERE c.is_active = true
		ORDER BY change_pct ASC
		LIMIT ?
	`
	err := r.db.Raw(query, limit).Scan(&companies).Error
	return companies, err
}

func (r *repository) GetMostActive(limit int) ([]models.Company, error) {
	var companies []models.Company
	query := `
		SELECT c.*, sp.volume
		FROM companies c
		INNER JOIN LATERAL (
			SELECT volume FROM stock_prices 
			WHERE company_id = c.id AND timeframe = '1d'
			ORDER BY timestamp DESC LIMIT 1
		) sp ON true
		WHERE c.is_active = true
		ORDER BY sp.volume DESC
		LIMIT ?
	`
	err := r.db.Raw(query, limit).Scan(&companies).Error
	return companies, err
}

// Market events

func (r *repository) CreateMarketEvent(event *models.MarketEvent) error {
	query := `
		INSERT INTO market_events (company_id, event_type, title, description, impact_percentage, event_date, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	return r.db.Raw(query, event.CompanyID, event.EventType, event.Title, event.Description,
		event.ImpactPercentage, event.EventDate, time.Now()).Scan(&event.ID).Error
}

func (r *repository) GetUpcomingEvents(companyID uint, limit int) ([]models.MarketEvent, error) {
	var events []models.MarketEvent
	query := `
		SELECT * FROM market_events 
		WHERE company_id = ? AND event_date >= NOW()
		ORDER BY event_date ASC
		LIMIT ?
	`
	err := r.db.Raw(query, companyID, limit).Scan(&events).Error
	return events, err
}
