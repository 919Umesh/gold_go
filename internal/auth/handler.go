package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/919Umesh/stock_market_sim/pkg/apperr"
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
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Register(req.FullName, req.Email, req.Phone, req.Password, req.Role)
	if err != nil {
		if err == ErrUserExists {
			apperr.RespondWithMessage(c, http.StatusConflict, "user already exists")
			return
		}
		slog.Error("Registration error", "error", err)
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "registration failed")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    ToUserResponse(user),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message: "login successful",
		Token:   token,
		User:    ToUserResponse(user),
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	uid, ok := userID.(string)
	if !ok {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "invalid user id type")
		return
	}

	user, err := h.service.GetProfile(uid)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": ToUserResponse(user)})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	uid, ok := userID.(string)
	if !ok {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "invalid user id type")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
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
		apperr.RespondWithMessage(c, http.StatusBadRequest, "no fields to update")
		return
	}

	user, err := h.service.UpdateProfile(uid, updates)

	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "profile update failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    ToUserResponse(user),
	})
}

func (h *Handler) UpdateKYC(c *gin.Context) {
	userIDStr := c.Param("user_id")

	var request UpdateKYCAdmin
	if err := c.ShouldBindJSON(&request); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.UpdateUserKYCStatus(userIDStr, request.KYCStatus, request.Role)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "KYC update failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC status updated successfully",
		"user":    ToUserResponse(user),
	})
}

func (h *Handler) UploadProfileImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	uid, ok := userID.(string)
	if !ok {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "invalid user id type")
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "image too large (max 5MB)")
		return
	}

	user, err := h.service.UploadProfileImage(uid, file, header.Filename)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "failed to upload image: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile image uploaded successfully",
		"user":    ToUserResponse(user),
	})
}
