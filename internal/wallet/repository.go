package wallet

import (
	"time"

	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	WithLock(userID uint, fn func(*models.Wallet) error) error

	CreateTransaction(transaction *models.Transaction) error
	UpdateTransaction(transaction *models.Transaction) error

	GetUserTransaction(userID uint) ([]models.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `SELECT * FROM wallets WHERE user_id = ? LIMIT 1`
	err := r.db.Raw(query, userID).Scan(&wallet).Error
	if err != nil {
		return nil, err
	}
	if wallet.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &wallet, nil
}

func (r *repository) Create(wallet *models.Wallet) error {
	query := `INSERT INTO wallets (user_id, fiat_balance, gold_grams, locked, version) VALUES (?, ?, ?, ?, ?)`
	return r.db.Exec(query, wallet.UserID, wallet.FiatBalance, wallet.GoldGrams, wallet.Locked, wallet.Version).Error
}

func (r *repository) Update(wallet *models.Wallet) error {
	query := `
		UPDATE wallets 
		SET fiat_balance = ?, gold_grams = ?, locked = ?, version = ? 
		WHERE id = ?
	`
	return r.db.Exec(query, wallet.FiatBalance, wallet.GoldGrams, wallet.Locked, wallet.Version, wallet.ID).Error
}

func (r *repository) CreateTransaction(transaction *models.Transaction) error {
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	query := `
		INSERT INTO transactions 
		(user_id, type, amount, gold_grams, price_per_gram, status, reference_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	// returning id to populate transaction.ID
	return r.db.Raw(query, transaction.UserID, transaction.Type, transaction.Amount, transaction.GoldGrams, transaction.PricePerGram, transaction.Status, transaction.ReferenceID, transaction.CreatedAt, transaction.UpdatedAt).Scan(&transaction.ID).Error
}

func (r *repository) UpdateTransaction(transaction *models.Transaction) error {
	transaction.UpdatedAt = time.Now()
	query := `
		UPDATE transactions 
		SET user_id=?, type=?, amount=?, gold_grams=?, price_per_gram=?, status=?, reference_id=?, updated_at=? 
		WHERE id = ?
	`
	return r.db.Exec(query, transaction.UserID, transaction.Type, transaction.Amount, transaction.GoldGrams, transaction.PricePerGram, transaction.Status, transaction.ReferenceID, transaction.UpdatedAt, transaction.ID).Error
}

func (r *repository) GetUserTransaction(userID uint) ([]models.Transaction, error) {
	var transaction []models.Transaction
	query := ` 
			SELECT * 
			FROM transactions 
			WHERE user_id = ? 
			ORDER BY created_at DESC 
		`
	err := r.db.Raw(query, userID).Scan(&transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *repository) WithLock(userID uint, fn func(*models.Wallet) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		query := `SELECT * FROM wallets WHERE user_id = ? FOR UPDATE`
		if err := tx.Raw(query, userID).Scan(&wallet).Error; err != nil {
			return err
		}

		// If using raw SQL, we must verify a record was actually found
		if wallet.ID == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := fn(&wallet); err != nil {
			return err
		}

		updateQuery := `
			UPDATE wallets 
			SET fiat_balance = ?, gold_grams = ?, locked = ?, version = ? 
			WHERE id = ?
		`
		return tx.Exec(updateQuery, wallet.FiatBalance, wallet.GoldGrams, wallet.Locked, wallet.Version, wallet.ID).Error
	})
}
