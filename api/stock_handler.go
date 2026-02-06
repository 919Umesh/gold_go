package api

import (
	"net/http"
	"strconv"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service stock.Service
}

func NewStockHandler(service stock.Service) *StockHandler {
	return &StockHandler{service: service}
}

// GetCompany godoc
// @Summary Get company details
// @Tags stocks
// @Param symbol path string true "Company Symbol"
// @Success 200 {object} models.Company
// @Router /api/v1/stocks/{symbol} [get]
func (h *StockHandler) GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")

	company, err := h.service.GetCompany(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Company not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}

// ListCompanies godoc
// @Summary List all companies
// @Tags stocks
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.Company
// @Router /api/v1/stocks [get]
func (h *StockHandler) ListCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	companies, err := h.service.ListCompanies(limit, offset)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"limit":     limit,
		"offset":    offset,
	})
}

// SearchCompanies godoc
// @Summary Search companies
// @Tags stocks
// @Param q query string true "Search query"
// @Success 200 {array} models.Company
// @Router /api/v1/stocks/search [get]
func (h *StockHandler) SearchCompanies(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Search query required")
		return
	}

	companies, err := h.service.SearchCompanies(query)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Search failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}

// GetCompaniesBySector godoc
// @Summary Get companies by sector
// @Tags stocks
// @Param sector path string true "Sector name"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.Company
// @Router /api/v1/stocks/sector/{sector} [get]
func (h *StockHandler) GetCompaniesBySector(c *gin.Context) {
	sector := c.Param("sector")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	companies, err := h.service.GetCompaniesBySector(sector, limit, offset)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sector":    sector,
		"companies": companies,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetCurrentPrice godoc
// @Summary Get current stock price
// @Tags stocks
// @Param symbol path string true "Company Symbol"
// @Success 200 {object} models.StockPrice
// @Router /api/v1/stocks/{symbol}/price [get]
func (h *StockHandler) GetCurrentPrice(c *gin.Context) {
	symbol := c.Param("symbol")

	price, err := h.service.GetCurrentPrice(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Price not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

// GetPriceHistory godoc
// @Summary Get price history
// @Tags stocks
// @Param symbol path string true "Company Symbol"
// @Param timeframe query string false "Timeframe (1m, 5m, 15m, 1h, 1d, 1w, 1M)" default("1d")
// @Param days query int false "Number of days" default(30)
// @Success 200 {array} models.StockPrice
// @Router /api/v1/stocks/{symbol}/history [get]
func (h *StockHandler) GetPriceHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	timeframe := c.DefaultQuery("timeframe", "1d")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	prices, err := h.service.GetPriceHistory(symbol, timeframe, days)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "History not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol":    symbol,
		"timeframe": timeframe,
		"prices":    prices,
	})
}

// GetMarketOverview godoc
// @Summary Get market overview
// @Tags stocks
// @Success 200 {object} stock.MarketOverview
// @Router /api/v1/stocks/market-overview [get]
func (h *StockHandler) GetMarketOverview(c *gin.Context) {
	overview, err := h.service.GetMarketOverview()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch market overview")
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetTopGainers godoc
// @Summary Get top gaining stocks
// @Tags stocks
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} stock.CompanyWithChange
// @Router /api/v1/stocks/top-gainers [get]
func (h *StockHandler) GetTopGainers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	gainers, err := h.service.GetTopGainers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top gainers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"gainers": gainers})
}

// GetTopLosers godoc
// @Summary Get top losing stocks
// @Tags stocks
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} stock.CompanyWithChange
// @Router /api/v1/stocks/top-losers [get]
func (h *StockHandler) GetTopLosers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	losers, err := h.service.GetTopLosers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top losers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"losers": losers})
}

// GetMostActive godoc
// @Summary Get most active stocks
// @Tags stocks
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.Company
// @Router /api/v1/stocks/most-active [get]
func (h *StockHandler) GetMostActive(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	active, err := h.service.GetMostActive(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch most active")
		return
	}

	c.JSON(http.StatusOK, gin.H{"active": active})
}

// GetUpcomingEvents godoc
// @Summary Get upcoming events for a company
// @Tags stocks
// @Param symbol path string true "Company Symbol"
// @Success 200 {array} models.MarketEvent
// @Router /api/v1/stocks/{symbol}/events [get]
func (h *StockHandler) GetUpcomingEvents(c *gin.Context) {
	symbol := c.Param("symbol")

	events, err := h.service.GetUpcomingEvents(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Events not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
