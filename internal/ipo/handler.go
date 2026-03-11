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
func (h *Handler) ListIPOs(c *gin.Context) {
	ipos, err := h.service.ListIPOs(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ipos": ipos})
}

// GetIPO godoc
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

// GetIPOApplications godoc
func (h *Handler) GetIPOApplications(c *gin.Context) {
	ipoID := c.Param("id")
	if ipoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IPO ID required"})
		return
	}

	apps, err := h.service.GetIPOApplications(ipoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ipo_id":       ipoID,
		"applications": apps,
		"count":        len(apps),
	})
}
