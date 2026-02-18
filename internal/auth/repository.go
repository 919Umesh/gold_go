package auth

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/google/uuid"
)

const (
	TableUsers         = "users"
	BucketUserProfiles = "user-profiles"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *models.User) error
	UploadProfileImage(file multipart.File, filename string) (string, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}


func (r *repository) Create(user *models.User) error {
	query := `INSERT INTO users (full_name, email, phone, password_hash, kyc_status, role)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	return r.client.ExecuteInsert(query, user,
		user.FullName, user.Email, user.Phone, user.PasswordHash, user.KYCStatus, user.Role)
}


func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	err := r.client.ExecuteQueryRow(query, &user, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *repository) FindByID(id string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.client.ExecuteQueryRow(query, &user, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}


func (r *repository) ExistsByEmail(email string) (bool, error) {
	_, err := r.FindByEmail(email)
	if err != nil {
		return false, nil 
	}
	return true, nil
}

func (r *repository) Update(user *models.User) error {
	query := `UPDATE users SET full_name = $1, phone = $2, kyc_status = $3, role = $4,
			  profile_image_id = $5 WHERE id = $6 RETURNING *`
	return r.client.ExecuteUpdate(query, user,
		user.FullName, user.Phone, user.KYCStatus, user.Role, user.ProfileImageID, user.ID)
}


func (r *repository) UploadProfileImage(f multipart.File, filename string) (string, error) {

	path := fmt.Sprintf("profiles/%s_%s", uuid.New().String(), filename)

	uploadURL := fmt.Sprintf("%s/object/%s/%s", r.client.StorageURL(), BucketUserProfiles, path)

	req, err := http.NewRequest("POST", uploadURL, io.Reader(f))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}


	r.client.SetHeaders(req)
	req.Header.Set("Content-Type", "image/*")


	resp, err := r.client.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload error (status %d): %s", resp.StatusCode, string(respBody))
	}
	publicURL := fmt.Sprintf("%s/object/public/%s/%s", r.client.StorageURL(), BucketUserProfiles, path)
	return publicURL, nil
}
