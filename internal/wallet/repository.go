package wallet

import (
	"github.com/umesh/gold_investment/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	WithLock(userID uint, fn func(*models.Wallet) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *repository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *repository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *repository) WithLock(userID uint, fn func(*models.Wallet) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		if err := fn(&wallet); err != nil {
			return err
		}

		return tx.Save(&wallet).Error
	})
}
