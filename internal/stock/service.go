package stock

import "github.com/919Umesh/stock_market_sim/models"

type Service interface {
	ListCompanies(limit, offset int) ([]models.Company, error)
	GetCompany(symbol string) (*models.Company, error)
	GetCompanyByID(id string) (*models.Company, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListCompanies(limit, offset int) ([]models.Company, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListCompanies(limit, offset)
}

func (s *service) GetCompany(symbol string) (*models.Company, error) {
	return s.repo.GetCompanyBySymbol(symbol)
}

func (s *service) GetCompanyByID(id string) (*models.Company, error) {
	return s.repo.GetCompanyByID(id)
}
