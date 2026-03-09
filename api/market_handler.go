package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// MarketHandler handles market data / SSE / price trigger endpoints
type MarketHandler struct {
	priceEngine   *market.PriceEngine
	eventHub      *market.EventHub
	triggerWorker *market.TriggerWorker
}

func NewMarketHandler(pe *market.PriceEngine, eh *market.EventHub, tw *market.TriggerWorker) *MarketHandler {
	return &MarketHandler{
		priceEngine:   pe,
		eventHub:      eh,
		triggerWorker: tw,
	}
}

// ──────────────────── Market Data ────────────────────

// GetLiveTradingData godoc
// @Summary Get live trading data
// @Description Get live price and volume data for all companies
// @Tags market
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/live [get]
func (h *MarketHandler) GetLiveTradingData(c *gin.Context) {
	data, err := h.priceEngine.GetLiveTradingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetMarketIndex godoc
// @Summary Get market index
// @Description Get overall market indicators (advances, declines, market cap)
// @Tags market
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/index [get]
func (h *MarketHandler) GetMarketIndex(c *gin.Context) {
	idx, err := h.priceEngine.GetMarketIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"index": idx})
}

// GetCandlestickData godoc
// @Summary Get candlestick data
// @Description Get OHLCV data for a specific symbol
// @Tags market
// @Produce json
// @Param symbol path string true "Company Symbol"
// @Param timeframe query string false "Timeframe (default 1D)"
// @Param days query int false "Number of days (default 30)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /market/companies/{symbol}/candles [get]
func (h *MarketHandler) GetCandlestickData(c *gin.Context) {
	symbol := c.Param("symbol")
	timeframe := c.DefaultQuery("timeframe", "1D")
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)
	if days <= 0 {
		days = 30
	}

	candles, err := h.priceEngine.GetCandlestickData(symbol, timeframe, days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"candles": candles})
}

// ──────────────────── SSE Stream ────────────────────

// StreamPrices godoc
// @Summary Stream prices (SSE)
// @Description Real-time price updates via Server-Sent Events
// @Tags market
// @Produce text/event-stream
// @Success 200 {string} string "SSE stream"
// @Router /market/stream [get]
func (h *MarketHandler) StreamPrices(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ch := h.eventHub.Subscribe()
	defer h.eventHub.Unsubscribe(ch)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// ──────────────────── Price Triggers ────────────────────

type CreateTriggerRequest struct {
	CompanyID    string `json:"company_id" binding:"required"`
	TriggerPrice string `json:"trigger_price" binding:"required"`
	SharesQty    int64  `json:"shares_qty" binding:"required,gt=0"`
	Direction    string `json:"direction" binding:"required"` // above, below
}

// CreateTrigger godoc
// @Summary Create price trigger
// @Description Create an auto-sell trigger based on price target
// @Tags trading
// @Accept json
// @Produce json
// @Param request body CreateTriggerRequest true "Trigger details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /triggers [post]
func (h *MarketHandler) CreateTrigger(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := decimal.NewFromString(req.TriggerPrice)
	if err != nil || !price.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trigger_price"})
		return
	}

	if req.Direction != models.TriggerDirectionAbove && req.Direction != models.TriggerDirectionBelow {
		c.JSON(http.StatusBadRequest, gin.H{"error": "direction must be 'above' or 'below'"})
		return
	}

	trigger := &models.PriceTrigger{
		UserID:       userID,
		CompanyID:    req.CompanyID,
		TriggerPrice: price,
		SharesQty:    req.SharesQty,
		Direction:    req.Direction,
		Status:       models.TriggerStatusActive,
	}

	if err := h.triggerWorker.CreateTrigger(trigger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "price trigger created",
		"trigger": trigger,
	})
}

// CancelTrigger godoc
// @Summary Cancel price trigger
// @Description Cancel an active price trigger
// @Tags trading
// @Produce json
// @Param id path string true "Trigger ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /triggers/{id} [delete]
func (h *MarketHandler) CancelTrigger(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	triggerID := c.Param("id")
	if err := h.triggerWorker.CancelTrigger(triggerID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trigger cancelled"})
}

// GetUserTriggers godoc
// @Summary Get user's triggers
// @Description Get a list of active price triggers for the current user
// @Tags trading
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /triggers [get]
func (h *MarketHandler) GetUserTriggers(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	triggers, err := h.triggerWorker.GetUserTriggers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"triggers": triggers})
}

// ──────────────────── Companies List ────────────────────

// ListCompanies godoc
// @Summary List companies
// @Description Get a list of all trading companies with summary data
// @Tags market
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/companies [get]
func (h *MarketHandler) ListCompanies(c *gin.Context) {
	data, err := h.priceEngine.GetLiveTradingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type companyInfo struct {
		Symbol   string          `json:"symbol"`
		Name     string          `json:"company_name"`
		Sector   string          `json:"sector"`
		Price    decimal.Decimal `json:"current_price"`
		Change   decimal.Decimal `json:"change_percent"`
		Volume   int64           `json:"volume"`
		Turnover decimal.Decimal `json:"turnover"`
	}

	companies := make([]companyInfo, 0, len(data))
	for _, d := range data {
		companies = append(companies, companyInfo{
			Symbol:   d.Symbol,
			Name:     d.CompanyName,
			Sector:   d.Sector,
			Price:    d.LTP,
			Change:   d.ChangePercent,
			Volume:   d.Volume,
			Turnover: d.Turnover,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"count":     len(companies),
	})
}

// GetCompanyDetail godoc
// @Summary Get company details
// @Description Get detailed info for a specific company by symbol
// @Tags market
// @Produce json
// @Param symbol path string true "Company Symbol"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /market/companies/{symbol} [get]
func (h *MarketHandler) GetCompanyDetail(c *gin.Context) {
	symbol := c.Param("symbol")
	_ = symbol
	// All live data
	data, err := h.priceEngine.GetLiveTradingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, d := range data {
		if d.Symbol == symbol {
			c.JSON(http.StatusOK, gin.H{"company": d})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("company %s not found", symbol)})
}
