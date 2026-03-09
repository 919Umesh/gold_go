package models

import (
	"time"
)

type User struct {
	ID              string    `json:"id,omitempty"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	PasswordHash    string    `json:"password_hash,omitempty"`
	KYCStatus       string    `json:"kyc_status"`
	Role            string    `json:"role"`
	IsAdmin         bool      `json:"is_admin"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}
