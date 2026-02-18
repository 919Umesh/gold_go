package stock

import (
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

type Service interface {
	GetCompany(symbol string) (*models.Company, error)
	ListCompanies(limit, offset int) ([]models.Company, error)
	ListCompaniesWithTotal(limit, offset int) ([]models.Company, int, error)
	SearchCompanies(query string) ([]models.Company, error)
	GetCompaniesBySector(sector string, limit, offset int) ([]models.Company, error)
	GetCompaniesBySectorWithTotal(sector string, limit, offset int) ([]models.Company, int, error)
	GetAllSectors() ([]string, error)

	GetCurrentPrice(symbol string) (*models.StockPrice, error)
	GetPriceHistory(symbol string, timeframe string, days int) ([]models.StockPrice, error)

	GetMarketOverview() (*MarketOverview, error)
	GetTopGainers(limit int) ([]CompanyWithChange, error)
	GetTopLosers(limit int) ([]CompanyWithChange, error)
	GetMostActive(limit int) ([]models.Company, error)

	GetUpcomingEvents(symbol string) ([]models.MarketEvent, error)
}

type MarketOverview struct {
	TotalCompanies int                 `json:"total_companies"`
	TopGainers     []CompanyWithChange `json:"top_gainers"`
	TopLosers      []CompanyWithChange `json:"top_losers"`
	MostActive     []models.Company    `json:"most_active"`
}

type CompanyWithChange struct {
	models.Company
	CurrentPrice  decimal.Decimal `json:"current_price"`
	PreviousPrice decimal.Decimal `json:"previous_price"`
	Change        decimal.Decimal `json:"change"`
	ChangePercent decimal.Decimal `json:"change_percent"`
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCompany(symbol string) (*models.Company, error) {
	return s.repo.GetCompanyBySymbol(symbol)
}

func (s *service) ListCompanies(limit, offset int) ([]models.Company, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListCompanies(limit, offset)
}

func (s *service) ListCompaniesWithTotal(limit, offset int) ([]models.Company, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	companies, err := s.repo.ListCompanies(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.GetTotalCompaniesCount()
	if err != nil {
		return companies, 0, nil 
	}
	return companies, total, nil
}

func (s *service) SearchCompanies(query string) ([]models.Company, error) {
	if query == "" {
		return []models.Company{}, nil
	}
	return s.repo.SearchCompanies(query, 20)
}

func (s *service) GetCompaniesBySector(sector string, limit, offset int) ([]models.Company, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListCompaniesBySector(sector, limit, offset)
}

func (s *service) GetCompaniesBySectorWithTotal(sector string, limit, offset int) ([]models.Company, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	companies, err := s.repo.ListCompaniesBySector(sector, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.GetTotalCompaniesBySectorCount(sector)
	if err != nil {
		return companies, 0, nil // Return companies even if count fails
	}
	return companies, total, nil
}

func (s *service) GetAllSectors() ([]string, error) {
	return s.repo.GetAllSectors()
}

func (s *service) GetCurrentPrice(symbol string) (*models.StockPrice, error) {
	company, err := s.repo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	return s.repo.GetLatestPrice(company.ID)
}

func (s *service) GetPriceHistory(symbol string, timeframe string, days int) ([]models.StockPrice, error) {
	company, err := s.repo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	to := time.Now()
	from := to.AddDate(0, 0, -days)

	return s.repo.GetPriceHistory(company.ID, timeframe, from, to, 1000)
}

func (s *service) GetMarketOverview() (*MarketOverview, error) {
	gainers, err := s.GetTopGainers(5)
	if err != nil {
		return nil, err
	}

	losers, err := s.GetTopLosers(5)
	if err != nil {
		return nil, err
	}

	active, err := s.GetMostActive(5)
	if err != nil {
		return nil, err
	}

	companies, _ := s.repo.ListCompanies(100, 0)
	totalCount := len(companies)

	return &MarketOverview{
		TotalCompanies: totalCount,
		TopGainers:     gainers,
		TopLosers:      losers,
		MostActive:     active,
	}, nil
}

func (s *service) GetTopGainers(limit int) ([]CompanyWithChange, error) {
	companies, err := s.repo.GetTopGainers(limit)
	if err != nil {
		return nil, err
	}

	return s.enrichWithPriceChange(companies)
}

func (s *service) GetTopLosers(limit int) ([]CompanyWithChange, error) {
	companies, err := s.repo.GetTopLosers(limit)
	if err != nil {
		return nil, err
	}

	return s.enrichWithPriceChange(companies)
}

func (s *service) GetMostActive(limit int) ([]models.Company, error) {
	return s.repo.GetMostActive(limit)
}

func (s *service) GetUpcomingEvents(symbol string) ([]models.MarketEvent, error) {
	company, err := s.repo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	return s.repo.GetUpcomingEvents(company.ID, 10)
}

func (s *service) enrichWithPriceChange(companies []models.Company) ([]CompanyWithChange, error) {
	result := make([]CompanyWithChange, 0, len(companies))

	for _, company := range companies {
		current, err := s.repo.GetLatestPrice(company.ID)
		if err != nil {
			continue
		}

		yesterday := time.Now().AddDate(0, 0, -1)
		previous, err := s.repo.GetPriceAtTime(company.ID, yesterday)
		if err != nil {
			continue
		}

		change := current.ClosePrice.Sub(previous.ClosePrice)
		changePercent := decimal.Zero
		if !previous.ClosePrice.IsZero() {
			changePercent = change.Div(previous.ClosePrice).Mul(decimal.NewFromInt(100))
		}

		result = append(result, CompanyWithChange{
			Company:       company,
			CurrentPrice:  current.ClosePrice,
			PreviousPrice: previous.ClosePrice,
			Change:        change,
			ChangePercent: changePercent,
		})
	}

	return result, nil
}
