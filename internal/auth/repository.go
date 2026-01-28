package auth

import (
	"time"

	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *models.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users 
		SET full_name = ?, email = ?, phone = ?, password_hash = ?, kyc_status = ?, role = ?, updated_at = ?
		WHERE id = ?
	`
	return r.db.Exec(query, user.FullName, user.Email, user.Phone, user.PasswordHash, user.KYCStatus, user.Role, user.UpdatedAt, user.ID).Error
}

func (r *repository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (full_name, email, phone, password_hash, kyc_status, role, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING id
	`
	// returning ID since GORM model usually expects it to be populated
	err := r.db.Raw(query, user.FullName, user.Email, user.Phone, user.PasswordHash, user.KYCStatus, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID).Error
	return err
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = ? LIMIT 1`
	err := r.db.Raw(query, email).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *repository) FindByID(id uint) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = ? LIMIT 1`
	err := r.db.Raw(query, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *repository) ExistsByEmail(email string) (bool, error) {
	var count int64
	query := `SELECT count(*) FROM users WHERE email = ?`
	err := r.db.Raw(query, email).Scan(&count).Error
	return count > 0, err
}
