package appwrite

import (
	"fmt"

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
	endpoint := "https://fra.cloud.appwrite.io/v1"
	projectID := "698599a4000fcd54a56a"
	apiKey := "standard_54423236d4faade3f231047beafa66688bd930cc2c63de086bb76cf3939f564d993019d5da21b7e29a3b5f49e39d1010e9065777540b31616b304a4f1357580e1902891d2b83191ae092e23c80498fe432d821fe7e4ad38256890df1f716dec2340e2712dcb1eee93e979f59dc920ba58193b7cf3f29eff9ba46e8128f42dd96"
	databaseID := "6985bccb00216e0a7171"

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
		// Default or fail? Let's require it for now to avoid ambiguity
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
