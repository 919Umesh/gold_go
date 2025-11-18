package wallet

import (
	"log"

	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	WithLock(userID uint, fn func(*models.Wallet) error) error

	//Add the transaction data
	CreateTransaction(transaction *models.Transaction) error
	UpdateTransaction(transaction *models.Transaction) error
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

func (r *repository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *repository) UpdateTransaction(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *repository) GetUserTransaction(userID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	query := ` 
				SELECT * 
				FROM transactions 
				WHERE user_id = ? 
				ORDER BY created_at ASC 
				`
	err := r.db.Raw(query, userID).Scan(&transaction).Error

	if err != nil {
		return nil, err
	}
	return &transaction, nil
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
		log.Print("Before saving the data")
		log.Print(&wallet)
		return tx.Save(&wallet).Error
	})
}
