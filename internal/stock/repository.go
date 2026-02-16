package stock

import (
	"fmt"
	"sort"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
	"github.com/shopspring/decimal"
)

const (
	CollectionCompanies    = "companies"
	CollectionStockPrices  = "stock_prices"
	CollectionMarketEvents = "market_events"
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

	CreateStockPrice(price *models.StockPrice) error
	GetLatestPrice(companyID string) (*models.StockPrice, error)
	GetPriceHistory(companyID string, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error)
	GetPriceAtTime(companyID string, timestamp time.Time) (*models.StockPrice, error)

	GetTopGainers(limit int) ([]models.Company, error)
	GetTopLosers(limit int) ([]models.Company, error)
	GetMostActive(limit int) ([]models.Company, error)

	CreateMarketEvent(event *models.MarketEvent) error
	GetUpcomingEvents(companyID string, limit int) ([]models.MarketEvent, error)
}

type repository struct {
	client *appwrite.Client
}

func NewRepository(client *appwrite.Client) Repository {
	return &repository{client: client}
}

func (r *repository) CreateCompany(company *models.Company) error {
	data := map[string]interface{}{
		"symbol":       company.Symbol,
		"name":         company.Name,
		"sector":       company.Sector,
		"market_cap":   company.MarketCap.InexactFloat64(),
		"description":  company.Description,
		"founded_year": company.FoundedYear,
		"employees":    company.Employees,
		"is_active":    company.IsActive,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, company)
}

func (r *repository) GetCompanyByID(id string) (*models.Company, error) {
	doc, err := r.client.Databases.GetDocument(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		id,
	)
	if err != nil {
		return nil, err
	}

	var company models.Company
	if err := appwrite.Decode(doc, &company); err != nil {
		return nil, fmt.Errorf("failed to decode company: %w", err)
	}
	return &company, nil
}

func (r *repository) GetCompanyBySymbol(symbol string) (*models.Company, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		appwrite.WithListDocumentsQueries([]string{
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
	if err := appwrite.DecodeListItem(resp, 0, &company); err != nil {
		return nil, fmt.Errorf("failed to decode company: %w", err)
	}
	return &company, nil
}

func (r *repository) ListCompanies(limit, offset int) ([]models.Company, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("is_active", true),
			query.OrderDesc("market_cap"),
			query.Limit(limit),
			query.Offset(offset),
		}),
	)
	if err != nil {
		return nil, err
	}

	var companies []models.Company
	for i := range resp.Documents {
		var c models.Company
		if err := appwrite.DecodeListItem(resp, i, &c); err == nil {
			companies = append(companies, c)
		}
	}
	return companies, nil
}

func (r *repository) ListCompaniesBySector(sector string, limit, offset int) ([]models.Company, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("sector", sector),
			query.Equal("is_active", true),
			query.OrderDesc("market_cap"),
			query.Limit(limit),
			query.Offset(offset),
		}),
	)
	if err != nil {
		return nil, err
	}

	var companies []models.Company
	for i := range resp.Documents {
		var c models.Company
		if err := appwrite.DecodeListItem(resp, i, &c); err == nil {
			companies = append(companies, c)
		}
	}
	return companies, nil
}

func (r *repository) SearchCompanies(q string, limit int) ([]models.Company, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("is_active", true),
			query.Search("name", q),
			query.Limit(limit),
		}),
	)
	if err != nil {
		return nil, err
	}

	var companies []models.Company
	for i := range resp.Documents {
		var c models.Company
		if err := appwrite.DecodeListItem(resp, i, &c); err == nil {
			companies = append(companies, c)
		}
	}
	return companies, nil
}

func (r *repository) UpdateCompany(company *models.Company) error {
	data := map[string]interface{}{
		"name":         company.Name,
		"sector":       company.Sector,
		"market_cap":   company.MarketCap.InexactFloat64(),
		"description":  company.Description,
		"founded_year": company.FoundedYear,
		"employees":    company.Employees,
		"is_active":    company.IsActive,
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		company.ID,
		r.client.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, company)
}

func (r *repository) CreateStockPrice(price *models.StockPrice) error {
	data := map[string]interface{}{
		"company_id":  price.CompanyID,
		"open_price":  price.OpenPrice.InexactFloat64(),
		"high_price":  price.HighPrice.InexactFloat64(),
		"low_price":   price.LowPrice.InexactFloat64(),
		"close_price": price.ClosePrice.InexactFloat64(),
		"volume":      price.Volume,
		"timestamp":   price.Timestamp.Format(time.RFC3339),
		"timeframe":   price.Timeframe,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionStockPrices,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, price)
}

func (r *repository) GetLatestPrice(companyID string) (*models.StockPrice, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionStockPrices,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("company_id", companyID),
			query.Equal("timeframe", "1D"), // Changed from 1d to 1D to match data
			query.OrderDesc("timestamp"),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}
	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("no price data found")
	}

	var price models.StockPrice
	if err := appwrite.DecodeListItem(resp, 0, &price); err != nil {
		return nil, fmt.Errorf("failed to decode price: %w", err)
	}
	return &price, nil
}

func (r *repository) GetPriceHistory(companyID string, timeframe string, from, to time.Time, limit int) ([]models.StockPrice, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionStockPrices,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("company_id", companyID),
			query.Equal("timeframe", timeframe),
			query.Between("timestamp", from.Format(time.RFC3339), to.Format(time.RFC3339)),
			query.OrderDesc("timestamp"),
			query.Limit(limit),
		}),
	)
	if err != nil {
		return nil, err
	}

	var prices []models.StockPrice
	for i := range resp.Documents {
		var p models.StockPrice
		if err := appwrite.DecodeListItem(resp, i, &p); err == nil {
			prices = append(prices, p)
		}
	}
	return prices, nil
}

func (r *repository) GetPriceAtTime(companyID string, timestamp time.Time) (*models.StockPrice, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionStockPrices,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("company_id", companyID),
			query.LessThanEqual("timestamp", timestamp.Format(time.RFC3339)),
			query.OrderDesc("timestamp"),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}
	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("no price data found")
	}

	var price models.StockPrice
	if err := appwrite.DecodeListItem(resp, 0, &price); err != nil {
		return nil, fmt.Errorf("failed to decode price: %w", err)
	}
	return &price, nil
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
	data := map[string]interface{}{
		"company_id":        event.CompanyID,
		"event_type":        event.EventType,
		"title":             event.Title,
		"description":       event.Description,
		"impact_percentage": event.ImpactPercentage,
		"event_date":        event.EventDate.Format(time.RFC3339),
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionMarketEvents,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, event)
}

func (r *repository) GetUpcomingEvents(companyID string, limit int) ([]models.MarketEvent, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionMarketEvents,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("company_id", companyID),
			query.GreaterThanEqual("event_date", time.Now().Format(time.RFC3339)),
			query.OrderAsc("event_date"),
			query.Limit(limit),
		}),
	)
	if err != nil {
		return nil, err
	}

	var events []models.MarketEvent
	for i := range resp.Documents {
		var e models.MarketEvent
		if err := appwrite.DecodeListItem(resp, i, &e); err == nil {
			events = append(events, e)
		}
	}
	return events, nil
}

// GetAllSectors returns a unique list of all sectors
func (r *repository) GetAllSectors() ([]string, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionCompanies,
		appwrite.WithListDocumentsQueries([]string{
			query.Limit(10000),
		}),
	)
	if err != nil {
		return nil, err
	}

	sectorMap := make(map[string]bool)
	for i := range resp.Documents {
		var company models.Company
		if err := appwrite.DecodeListItem(resp, i, &company); err == nil {
			if company.Sector != "" {
				sectorMap[company.Sector] = true
			}
		}
	}

	sectors := make([]string, 0, len(sectorMap))
	for sector := range sectorMap {
		sectors = append(sectors, sector)
	}
	sort.Strings(sectors)

	return sectors, nil
}
