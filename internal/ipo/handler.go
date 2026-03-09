package ipo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ──────────────────── Admin: Create Company ────────────────────

type CreateCompanyRequest struct {
	Symbol      string `json:"symbol" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Sector      string `json:"sector"`
	TotalSupply int64  `json:"total_supply" binding:"required,gt=0"`
}

// CreateCompany godoc
// @Summary Create a company (Admin only)
// @Description Register a new company in the system
// @Tags admin
// @Accept json
// @Produce json
// @Param request body CreateCompanyRequest true "Company details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/companies [post]
func (h *Handler) CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Sector == "" {
		req.Sector = "General"
	}

	company, err := h.service.CreateCompany(req.Symbol, req.Name, req.Sector, req.TotalSupply)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "company created",
		"company": company,
	})
}

// ──────────────────── Admin: Launch IPO ────────────────────

type LaunchIPORequest struct {
	CompanyID       string `json:"company_id" binding:"required"`
	PricePerShare   string `json:"price_per_share" binding:"required"`
	TotalShares     int64  `json:"total_shares" binding:"required,gt=0"`
	MaxPerApplicant int64  `json:"max_per_applicant" binding:"required,gt=0"`
	OpenAt          string `json:"open_at" binding:"required"`
	CloseAt         string `json:"close_at" binding:"required"`
}

// LaunchIPO godoc
// @Summary Launch an IPO (Admin only)
// @Description Start an IPO for a specific company
// @Tags admin
// @Accept json
// @Produce json
// @Param request body LaunchIPORequest true "IPO details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/ipos [post]
func (h *Handler) LaunchIPO(c *gin.Context) {
	var req LaunchIPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := decimal.NewFromString(req.PricePerShare)
	if err != nil || !price.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price_per_share"})
		return
	}

	openAt, err := time.Parse(time.RFC3339, req.OpenAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid open_at format (use RFC3339)"})
		return
	}

	closeAt, err := time.Parse(time.RFC3339, req.CloseAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid close_at format (use RFC3339)"})
		return
	}

	ipo, err := h.service.LaunchIPO(req.CompanyID, price, req.TotalShares, req.MaxPerApplicant, openAt, closeAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "IPO launched",
		"ipo":     ipo,
	})
}

// ──────────────────── Admin: Allocate IPO ────────────────────

// AllocateIPO godoc
// @Summary Allocate IPO shares (Admin only)
// @Description Trigger the lottery-based allocation for a closed IPO
// @Tags admin
// @Produce json
// @Param id path string true "IPO ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/ipos/{id}/allocate [post]
func (h *Handler) AllocateIPO(c *gin.Context) {
	ipoID := c.Param("id")
	if ipoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IPO ID required"})
		return
	}

	result, err := h.service.AllocateIPO(ipoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "IPO allocation complete",
		"result":  result,
	})
}

// ──────────────────── User: Apply for IPO ────────────────────

type ApplyIPORequest struct {
	SharesRequested int64 `json:"shares_requested" binding:"required,gt=0"`
}

// ApplyForIPO godoc
// @Summary Apply for IPO
// @Description Apply for shares in an open IPO
// @Tags ipo
// @Accept json
// @Produce json
// @Param id path string true "IPO ID"
// @Param request body ApplyIPORequest true "Application details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /ipos/{id}/apply [post]
func (h *Handler) ApplyForIPO(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipoID := c.Param("id")
	if ipoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IPO ID required"})
		return
	}

	var req ApplyIPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.service.ApplyForIPO(userID, ipoID, req.SharesRequested)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "IPO application submitted",
		"application": app,
	})
}

// ──────────────────── List / Detail ────────────────────

// ListIPOs godoc
// @Summary List all IPOs
// @Description Get a list of all IPOs (active, closed, allocated)
// @Tags ipo
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ipos [get]
func (h *Handler) ListIPOs(c *gin.Context) {
	ipos, err := h.service.ListIPOs(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ipos": ipos})
}

// GetIPO godoc
// @Summary Get IPO details
// @Description Get details of a specific IPO by ID
// @Tags ipo
// @Produce json
// @Param id path string true "IPO ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /ipos/{id} [get]
func (h *Handler) GetIPO(c *gin.Context) {
	ipoID := c.Param("id")
	if ipoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IPO ID required"})
		return
	}

	ipo, err := h.service.GetIPO(ipoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ipo": ipo})
}
