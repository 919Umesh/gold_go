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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/profile [get]
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

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's full name or phone
// @Tags auth
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Profile updates"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/profile/update [put]
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

// UpdateKYC godoc
// @Summary Update user KYC and Role (Admin only)
// @Description Update KYC status and Role of a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body UpdateKYCAdmin true "KYC and Role updates"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/users/{user_id}/kyc [put]
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

// UploadProfileImage godoc
// @Summary Upload profile image
// @Description Upload profile image for current user (max 5MB)
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/profile/image [post]
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
