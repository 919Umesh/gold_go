package auth

import (
	"errors"
	"fmt"

	"github.com/umesh/gold_investment/models"
	"github.com/umesh/gold_investment/pkg/utils"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Register(fullName, email, phone, password string) (*models.User, error)
	Login(email, password string) (*models.User, string, error)
	GetProfile(userID uint) (*models.User, error)
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

func (s *service) Register(fullName, email, phone, password string) (*models.User, error) {
	exists, err := s.repo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	return user, nil
}

func (s *service) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("token generation failed: %w", err)
	}

	return user, token, nil
}

func (s *service) GetProfile(userID uint) (*models.User, error) {
	return s.repo.FindByID(userID)
}
