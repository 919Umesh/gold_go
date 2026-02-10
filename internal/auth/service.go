package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/919Umesh/stock_market_sim/pkg/utils"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Thread-safe in-memory cache for user lookups
	cacheMu   sync.RWMutex
	userCache = make(map[string]*models.User)
)

type Service interface {
	Register(fullName, email, phone, password, role string) (*models.User, error)
	Login(email, password string) (*models.User, string, error)
	GetProfile(userID string) (*models.User, error)
	UpdateProfile(userID string, updates map[string]interface{}) (*models.User, error)
	UpdateUserKYCStatus(userID string, kycStatus, role string) (*models.User, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *service) Register(fullName, email, phone, password, role string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
		Role:         role,
		KYCStatus:    "pending",
	}

	if err := s.repo.Create(user); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	// Cache the user after successful registration
	cacheMu.Lock()
	userCache[email] = user
	cacheMu.Unlock()

	slog.Info("user registered successfully", "email", email, "user_id", user.ID)
	return user, nil
}

func (s *service) Login(email, password string) (*models.User, string, error) {
	// Check cache first (thread-safe)
	cacheMu.RLock()
	cachedUser, exists := userCache[email]
	cacheMu.RUnlock()

	if exists {
		if err := utils.ComparePassword(cachedUser.PasswordHash, password); err == nil {
			token, err := utils.GenerateToken(cachedUser.ID, s.jwtSecret)
			if err != nil {
				return nil, "", fmt.Errorf("token generation failed: %w", err)
			}
			return cachedUser, token, nil
		}
		return nil, "", ErrInvalidCredentials
	}

	// Not in cache, query repository
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Cache the user for future lookups
	cacheMu.Lock()
	userCache[email] = user
	cacheMu.Unlock()

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("token generation failed: %w", err)
	}

	return user, token, nil
}

func (s *service) UpdateProfile(userID string, updates map[string]interface{}) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if fullname, ok := updates["full_name"].(string); ok {
		user.FullName = fullname
	}

	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("profile update failed: %w", err)
	}

	// Invalidate cache for this user
	cacheMu.Lock()
	delete(userCache, user.Email)
	cacheMu.Unlock()

	return user, nil
}

func (s *service) UpdateUserKYCStatus(userID string, kycStatus, role string) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.KYCStatus = kycStatus
	user.Role = role

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("KYC status update failed: %w", err)
	}

	return user, nil
}

func (s *service) GetProfile(userID string) (*models.User, error) {
	return s.repo.FindByID(userID)
}
