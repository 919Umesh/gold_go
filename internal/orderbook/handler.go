package orderbook

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ──────────────────── Place Buy Order ────────────────────

type PlaceBuyOrderRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,gt=0"`
	Price     string `json:"price" binding:"required"`
	OrderType string `json:"order_type"` // limit (default), market
}

// PlaceBuyOrder godoc
func (h *Handler) PlaceBuyOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req PlaceBuyOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil || !price.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	orderType := req.OrderType
	if orderType == "" {
		orderType = "limit"
	}

	order, matches, err := h.service.PlaceBuyOrder(userID, req.CompanyID, req.Quantity, price, orderType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "buy order placed",
		"order":   order,
		"matches": len(matches),
	})
}

// ──────────────────── Place Sell Order ────────────────────

type PlaceSellOrderRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,gt=0"`
	Price     string `json:"price" binding:"required"`
}

// PlaceSellOrder godoc
func (h *Handler) PlaceSellOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req PlaceSellOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil || !price.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	order, matches, err := h.service.PlaceSellOrder(userID, req.CompanyID, req.Quantity, price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "sell order placed",
		"order":   order,
		"matches": len(matches),
	})
}

// ──────────────────── Cancel Order ────────────────────

// CancelOrder godoc
func (h *Handler) CancelOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID required"})
		return
	}

	if err := h.service.CancelOrder(userID, orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

// ──────────────────── Order Book ────────────────────

// GetOrderBook godoc
func (h *Handler) GetOrderBook(c *gin.Context) {
	companyID := c.Param("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id required"})
		return
	}

	book, err := h.service.GetOrderBook(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_book": book})
}

// ──────────────────── User Orders ────────────────────

// GetUserOrders godoc
func (h *Handler) GetUserOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.service.GetUserOrders(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// ──────────────────── Portfolio ────────────────────

// GetPortfolio godoc
func (h *Handler) GetPortfolio(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	portfolio, err := h.service.GetPortfolio(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

// ──────────────────── User Trades ────────────────────

// GetUserTrades godoc
func (h *Handler) GetUserTrades(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trades, err := h.service.GetUserTrades(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trades": trades})
}

// ──────────────────── Company Trades ────────────────────

// GetCompanyTrades godoc
func (h *Handler) GetCompanyTrades(c *gin.Context) {
	companyID := c.Param("id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company ID required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	transactions, err := h.service.GetCompanyTransactions(companyID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id":   companyID,
		"transactions": transactions,
		"count":        len(transactions),
	})
}
