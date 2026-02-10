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

func (h *StockHandler) GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")

	company, err := h.service.GetCompany(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Company not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}

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

func (h *StockHandler) GetCurrentPrice(c *gin.Context) {
	symbol := c.Param("symbol")

	price, err := h.service.GetCurrentPrice(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Price not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

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

func (h *StockHandler) GetMarketOverview(c *gin.Context) {
	overview, err := h.service.GetMarketOverview()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch market overview")
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *StockHandler) GetTopGainers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	gainers, err := h.service.GetTopGainers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top gainers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"gainers": gainers})
}

func (h *StockHandler) GetTopLosers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	losers, err := h.service.GetTopLosers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top losers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"losers": losers})
}

func (h *StockHandler) GetMostActive(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	active, err := h.service.GetMostActive(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch most active")
		return
	}

	c.JSON(http.StatusOK, gin.H{"active": active})
}

func (h *StockHandler) GetUpcomingEvents(c *gin.Context) {
	symbol := c.Param("symbol")

	events, err := h.service.GetUpcomingEvents(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Events not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
