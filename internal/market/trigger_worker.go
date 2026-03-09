package market

import (
	"log/slog"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

type TriggerWorker struct {
	client     *supabase.Client
	placeOrder func(userID, companyID string, qty int64, price decimal.Decimal) error
}

func NewTriggerWorker(client *supabase.Client) *TriggerWorker {
	return &TriggerWorker{client: client}
}

func (tw *TriggerWorker) SetOrderPlacer(fn func(userID, companyID string, qty int64, price decimal.Decimal) error) {
	tw.placeOrder = fn
}

func (tw *TriggerWorker) CheckTriggers(companyID string, currentPrice decimal.Decimal) {
	triggers, err := tw.getActiveTriggers(companyID)
	if err != nil {
		slog.Error("TriggerWorker: failed to get triggers", "company_id", companyID, "error", err)
		return
	}

	for _, trigger := range triggers {
		shouldFire := false

		switch trigger.Direction {
		case models.TriggerDirectionAbove:
			if currentPrice.GreaterThanOrEqual(trigger.TriggerPrice) {
				shouldFire = true
			}
		case models.TriggerDirectionBelow:
			if currentPrice.LessThanOrEqual(trigger.TriggerPrice) {
				shouldFire = true
			}
		}

		if shouldFire {
			slog.Info("TriggerWorker: firing trigger",
				"trigger_id", trigger.ID,
				"user_id", trigger.UserID,
				"price", currentPrice.String(),
				"trigger_price", trigger.TriggerPrice.String(),
			)

			tw.updateTriggerStatus(trigger.ID, models.TriggerStatusTriggered)

			if tw.placeOrder != nil {
				if err := tw.placeOrder(trigger.UserID, trigger.CompanyID, trigger.SharesQty, currentPrice); err != nil {
					slog.Error("TriggerWorker: failed to place auto-sell", "error", err)
				}
			}
		}
	}
}

func (tw *TriggerWorker) getActiveTriggers(companyID string) ([]models.PriceTrigger, error) {
	var triggers []models.PriceTrigger
	query := "SELECT * FROM price_triggers WHERE company_id = $1 AND status = $2"
	err := tw.client.ExecuteQuery(query, &triggers, companyID, models.TriggerStatusActive)
	if err != nil {
		return nil, err
	}
	if triggers == nil {
		triggers = []models.PriceTrigger{}
	}
	return triggers, nil
}

func (tw *TriggerWorker) updateTriggerStatus(triggerID, status string) {
	query := `UPDATE price_triggers SET status = $1 WHERE id = $2 RETURNING *`
	var t models.PriceTrigger
	if err := tw.client.ExecuteUpdate(query, &t, status, triggerID); err != nil {
		slog.Error("TriggerWorker: failed to update trigger status", "error", err)
	}
}

func (tw *TriggerWorker) CreateTrigger(trigger *models.PriceTrigger) error {
	query := `INSERT INTO price_triggers (user_id, company_id, trigger_price, shares_qty, direction, status)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	return tw.client.ExecuteInsert(query, trigger,
		trigger.UserID, trigger.CompanyID, trigger.TriggerPrice.String(),
		trigger.SharesQty, trigger.Direction, trigger.Status)
}

func (tw *TriggerWorker) CancelTrigger(triggerID, userID string) error {
	var trigger models.PriceTrigger
	query := "SELECT * FROM price_triggers WHERE id = $1"
	err := tw.client.ExecuteQueryRow(query, &trigger, triggerID)
	if err != nil {
		return err
	}
	if trigger.UserID != userID {
		return models.ErrNotOwner
	}

	tw.updateTriggerStatus(triggerID, models.TriggerStatusCancelled)
	return nil
}

func (tw *TriggerWorker) GetUserTriggers(userID string) ([]models.PriceTrigger, error) {
	var triggers []models.PriceTrigger
	query := "SELECT * FROM price_triggers WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2"
	err := tw.client.ExecuteQuery(query, &triggers, userID, 50)
	if err != nil {
		return nil, err
	}
	if triggers == nil {
		triggers = []models.PriceTrigger{}
	}
	return triggers, nil
}
