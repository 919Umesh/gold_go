package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
	Role     string `json:"role" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Fullname string `json:"full_name,omitempty" binding:"omitempty,min=2,max=100"`
	Phone    string `json:"phone,omitempty" binding:"omitempty,min=10,max=15"`
}

type UpdateKYCAdmin struct {
	KYCStatus string `json:"kyc_status" binding:"required,oneof=pending verified rejected under_review"`
	Role      string `json:"role" binding:"required,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(req.FullName, req.Email, req.Phone, req.Password, req.Role)
	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	user, err := h.service.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	log.Print("--------------UserID---------------")
	log.Print(userID)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if req.Fullname != "" {
		updates["full_name"] = req.Fullname
	}

	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	user, err := h.service.UpdateProfile(userID.(uint), updates)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "profile update failed"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    user,
	})
}

func (h *Handler) UpdateKYC(c *gin.Context) {
	userIDStr := c.Param("user_id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	var request UpdateKYCAdmin
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUserKYCStatus(uint(userID), request.KYCStatus, request.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KYC update failed"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC status updated successfully",
		"user":    user,
	})
}
