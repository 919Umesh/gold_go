package models

import (
	"time"
)

type User struct {
	ID           string    `json:"$id,omitempty"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"password_hash"`
	KYCStatus    string    `json:"kyc_status"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"$createdAt,omitempty"`
	UpdatedAt    time.Time `json:"$updatedAt,omitempty"`

	Wallet *Wallet `json:"-"`
}
