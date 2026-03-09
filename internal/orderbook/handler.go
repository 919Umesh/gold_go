package orderbook

import (
	"net/http"

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
// @Summary Place a buy order
// @Description Place a new buy limit or market order
// @Tags trading
// @Accept json
// @Produce json
// @Param request body PlaceBuyOrderRequest true "Buy order details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/buy [post]
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
// @Summary Place a sell order
// @Description Place a new sell limit order
// @Tags trading
// @Accept json
// @Produce json
// @Param request body PlaceSellOrderRequest true "Sell order details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/sell [post]
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
// @Summary Cancel an order
// @Description Cancel an open buy or sell order
// @Tags trading
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/{id} [delete]
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
// @Summary Get company order book
// @Description Get current buy and sell orders for a specific company
// @Tags market
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {object} OrderBookView
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orderbook/{company_id} [get]
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
// @Summary Get user's orders
// @Description Get a list of open and historical orders for the current user
// @Tags trading
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/my [get]
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
// @Summary Get user portfolio
// @Description Get a list of all stock holdings for the current user
// @Tags trading
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /portfolio [get]
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
// @Summary Get user's trades
// @Description Get a list of all executed trades for the current user
// @Tags trading
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /trades [get]
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
