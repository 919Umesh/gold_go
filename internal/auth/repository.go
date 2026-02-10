package auth

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/models"

	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

const (
	CollectionUsers = "users"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *models.User) error
}

type repository struct {
	client *appwrite.Client
}

func NewRepository(client *appwrite.Client) Repository {
	return &repository{client: client}
}

func (r *repository) Create(user *models.User) error {
	data := map[string]interface{}{
		"full_name":     user.FullName,
		"email":         user.Email,
		"phone":         user.Phone,
		"password_hash": user.PasswordHash,
		"kyc_status":    user.KYCStatus,
		"role":          user.Role,
	}

	docID := id.Unique()

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionUsers,
		docID,
		data,
	)
	if err != nil {
		return err
	}

	// Use Decode to populate CreatedAt/UpdatedAt
	return appwrite.Decode(resp, user)
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionUsers,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("email", email),
			query.Limit(1),
		}),
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user models.User
	if err := appwrite.DecodeListItem(resp, 0, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

func (r *repository) FindByID(id string) (*models.User, error) {
	doc, err := r.client.Databases.GetDocument(
		r.client.Config.DatabaseID,
		CollectionUsers,
		id,
	)

	if err != nil {
		return nil, err
	}

	var user models.User
	if err := appwrite.Decode(doc, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

func (r *repository) ExistsByEmail(email string) (bool, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionUsers,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("email", email),
			query.Limit(1),
		}),
	)
	if err != nil {
		return false, err
	}
	return len(resp.Documents) > 0, nil
}

func (r *repository) Update(user *models.User) error {
	data := map[string]interface{}{
		"full_name":  user.FullName,
		"phone":      user.Phone,
		"kyc_status": user.KYCStatus,
		"role":       user.Role,
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionUsers,
		user.ID,
		r.client.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, user)
}
