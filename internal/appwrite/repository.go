package appwrite

import (
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

const (
	CollectionCompanies    = "companies"
	CollectionStockPrices  = "stock_prices"
	CollectionMarketEvents = "market_events"
)

type Repository interface {
	CreateCompany(company *models.Company) error
	GetCompanyBySymbol(symbol string) (*models.Company, error)
	CreateStockPrice(price *models.StockPrice) error
}

type repository struct {
	client *Client
}

func NewRepository(c *Client) Repository {
	return &repository{
		client: c,
	}
}

// CreateCompany creates a new company document
func (r *repository) CreateCompany(company *models.Company) error {
	data := map[string]interface{}{
		"symbol":       company.Symbol,
		"name":         company.Name,
		"sector":       company.Sector,
		"market_cap":   company.MarketCap,
		"description":  company.Description,
		"founded_year": company.FoundedYear,
		"employees":    company.Employees,
		"is_active":    company.IsActive,
	}

	doc, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}

	// Use Decode to populate system fields
	return Decode(doc, company)
}

// GetCompanyBySymbol fetches a company by its symbol
func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		WithListDocumentsQueries([]string{
			query.Equal("symbol", symbol),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("company not found")
	}

	var company models.Company
	if err := DecodeListItem(resp, 0, &company); err != nil {
		return nil, fmt.Errorf("failed to decode company: %w", err)
	}

	return &company, nil
}

// CreateStockPrice adds a new price record
func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	data := map[string]interface{}{
		"company_id": price.CompanyID,
		"open":       price.OpenPrice,
		"high":       price.HighPrice,
		"low":        price.LowPrice,
		"close":      price.ClosePrice,
		"volume":     price.Volume,
		"timestamp":  price.Timestamp.Format(time.RFC3339),
		"timeframe":  price.Timeframe,
	}

	doc, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionStockPrices,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}

	return Decode(doc, price)
}
