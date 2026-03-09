package event

import (
	"github.com/919Umesh/stock_market_sim/models"
)

type Service interface {
	GetEventsByCompany(companyID string, limit int) ([]models.CompanyEvent, error)
	GetUpcomingEvents(limit int) ([]models.CompanyEvent, error)
	GetEventsByType(eventType string, limit int) ([]models.CompanyEvent, error)
	GetAllEvents(limit, offset int) ([]models.CompanyEvent, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetEventsByCompany(companyID string, limit int) ([]models.CompanyEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetEventsByCompanyID(companyID, limit)
}

func (s *service) GetUpcomingEvents(limit int) ([]models.CompanyEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetUpcomingEvents(limit)
}

func (s *service) GetEventsByType(eventType string, limit int) ([]models.CompanyEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetEventsByType(eventType, limit)
}

func (s *service) GetAllEvents(limit, offset int) ([]models.CompanyEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetAllEvents(limit, offset)
}
