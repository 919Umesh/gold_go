package appwrite

import (
	"fmt"
	"os"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
)

type Client struct {
	Client    client.Client
	Databases *databases.Databases
	Config    *Config
}

type Config struct {
	Endpoint   string
	ProjectID  string
	APIKey     string
	DatabaseID string
}

func NewClient() (*Client, error) {
	endpoint := getEnv("APPWRITE_ENDPOINT", "https://fra.cloud.appwrite.io/v1")
	projectID := getEnv("APPWRITE_PROJECT_ID", "698599a4000fcd54a56a")
	apiKey := getEnv("APPWRITE_API_KEY", "")
	databaseID := getEnv("APPWRITE_DATABASE_ID", "6985bccb00216e0a7171")

	if endpoint == "" {
		endpoint = "https://fra.cloud.appwrite.io/v1"
	}
	if projectID == "" {
		return nil, fmt.Errorf("APPWRITE_PROJECT_ID is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("APPWRITE_API_KEY is required")
	}
	if databaseID == "" {
		return nil, fmt.Errorf("APPWRITE_DATABASE_ID is required")
	}

	c := client.New()
	c.Endpoint = endpoint
	c.AddHeader("X-Appwrite-Project", projectID)
	c.AddHeader("X-Appwrite-Key", apiKey)

	db := databases.New(c)

	return &Client{
		Client:    c,
		Databases: db,
		Config: &Config{
			Endpoint:   endpoint,
			ProjectID:  projectID,
			APIKey:     apiKey,
			DatabaseID: databaseID,
		},
	}, nil
}


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
