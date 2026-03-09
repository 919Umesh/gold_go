package event

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

type Repository interface {
	CreateEvent(event *models.CompanyEvent) error
	GetEventsByCompanyID(companyID string, limit int) ([]models.CompanyEvent, error)
	GetUpcomingEvents(limit int) ([]models.CompanyEvent, error)
	GetEventsByType(eventType string, limit int) ([]models.CompanyEvent, error)
	GetAllEvents(limit, offset int) ([]models.CompanyEvent, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

func (r *repository) CreateEvent(event *models.CompanyEvent) error {
	query := `INSERT INTO company_events (company_id, event_type, title, description, event_date, fiscal_year, status)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	return r.client.ExecuteInsert(query, event,
		event.CompanyID, event.EventType, event.Title, event.Description,
		event.EventDate.Format("2006-01-02T15:04:05Z07:00"), event.FiscalYear, event.Status)
}

func (r *repository) GetEventsByCompanyID(companyID string, limit int) ([]models.CompanyEvent, error) {
	var events []models.CompanyEvent
	query := "SELECT * FROM company_events WHERE company_id = $1 ORDER BY event_date DESC LIMIT $2"
	err := r.client.ExecuteQuery(query, &events, companyID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	if events == nil {
		events = []models.CompanyEvent{}
	}
	return events, nil
}

func (r *repository) GetUpcomingEvents(limit int) ([]models.CompanyEvent, error) {
	var events []models.CompanyEvent
	query := "SELECT * FROM company_events WHERE status = $1 ORDER BY event_date LIMIT $2"
	err := r.client.ExecuteQuery(query, &events, models.EventStatusUpcoming, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}
	if events == nil {
		events = []models.CompanyEvent{}
	}
	return events, nil
}

func (r *repository) GetEventsByType(eventType string, limit int) ([]models.CompanyEvent, error) {
	var events []models.CompanyEvent
	query := "SELECT * FROM company_events WHERE event_type = $1 ORDER BY event_date DESC LIMIT $2"
	err := r.client.ExecuteQuery(query, &events, eventType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}
	if events == nil {
		events = []models.CompanyEvent{}
	}
	return events, nil
}

func (r *repository) GetAllEvents(limit, offset int) ([]models.CompanyEvent, error) {
	var events []models.CompanyEvent
	query := "SELECT * FROM company_events ORDER BY event_date DESC LIMIT $1 OFFSET $2"
	err := r.client.ExecuteQuery(query, &events, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	if events == nil {
		events = []models.CompanyEvent{}
	}
	return events, nil
}
