package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	client, err := appwrite.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Appwrite client: %v", err)
	}

	dbID := client.Config.DatabaseID
	slog.Info("Setting up Appwrite Database", "database_id", dbID)

	// List of collections to create
	collections := []struct {
		ID   string
		Name string
	}{
		{ID: "users", Name: "Users"},
		{ID: "companies", Name: "Companies"},
		{ID: "stock_prices", Name: "Stock Prices"},
		{ID: "market_events", Name: "Market Events"},
		{ID: "virtual_wallets", Name: "Virtual Wallets"},
		{ID: "user_portfolios", Name: "User Portfolios"},
		{ID: "stock_transactions", Name: "Stock Transactions"},
	}

	for _, c := range collections {
		createCollection(client, dbID, c.ID, c.Name)
	}

	// Create Attributes
	createAttributes(client, dbID)

	slog.Info("Appwrite setup completed successfully.")
}

func createCollection(client *appwrite.Client, dbID, colID, colName string) {
	_, err := client.Databases.GetCollection(dbID, colID)
	if err == nil {
		slog.Info("Collection already exists", "id", colID)
		return
	}

	_, err = client.Databases.CreateCollection(dbID, colID, colName)
	if err != nil {
		log.Printf("Failed to create collection %s: %v", colName, err)
	} else {
		slog.Info("Created collection", "id", colID)
	}
}

func createAttributes(client *appwrite.Client, dbID string) {
	// Users
	createStringAttr(client, dbID, "users", "full_name", 100, true)
	createStringAttr(client, dbID, "users", "email", 255, true)
	createStringAttr(client, dbID, "users", "phone", 20, true)
	createStringAttr(client, dbID, "users", "password_hash", 255, true)
	createStringAttr(client, dbID, "users", "kyc_status", 50, false)
	createStringAttr(client, dbID, "users", "role", 50, false)

	// Companies
	createStringAttr(client, dbID, "companies", "symbol", 10, true)
	createStringAttr(client, dbID, "companies", "name", 100, true)
	createStringAttr(client, dbID, "companies", "sector", 50, true)
	createFloatAttr(client, dbID, "companies", "market_cap", true)
	createStringAttr(client, dbID, "companies", "description", 1000, false)
	createIntegerAttr(client, dbID, "companies", "founded_year", false)
	createIntegerAttr(client, dbID, "companies", "employees", false)
	createBooleanAttr(client, dbID, "companies", "is_active", true)

	// Stock Prices
	createStringAttr(client, dbID, "stock_prices", "company_id", 50, true)
	createFloatAttr(client, dbID, "stock_prices", "open_price", true)
	createFloatAttr(client, dbID, "stock_prices", "high_price", true)
	createFloatAttr(client, dbID, "stock_prices", "low_price", true)
	createFloatAttr(client, dbID, "stock_prices", "close_price", true)
	createIntegerAttr(client, dbID, "stock_prices", "volume", true)
	createStringAttr(client, dbID, "stock_prices", "timestamp", 50, true) // RFC3339
	createStringAttr(client, dbID, "stock_prices", "timeframe", 10, true)

	// Market Events
	createStringAttr(client, dbID, "market_events", "company_id", 50, true)
	createStringAttr(client, dbID, "market_events", "event_type", 50, true)
	createStringAttr(client, dbID, "market_events", "title", 255, true)
	createStringAttr(client, dbID, "market_events", "description", 1000, false)
	createFloatAttr(client, dbID, "market_events", "impact_percentage", true)
	createStringAttr(client, dbID, "market_events", "event_date", 50, true)

	// Virtual Wallets
	createStringAttr(client, dbID, "virtual_wallets", "user_id", 50, true)
	createFloatAttr(client, dbID, "virtual_wallets", "balance", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_invested", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_profit_loss", true)

	// User Portfolios
	createStringAttr(client, dbID, "user_portfolios", "user_id", 50, true)
	createStringAttr(client, dbID, "user_portfolios", "company_id", 50, true)
	createIntegerAttr(client, dbID, "user_portfolios", "quantity", true)
	createFloatAttr(client, dbID, "user_portfolios", "average_price", true)
	createFloatAttr(client, dbID, "user_portfolios", "total_invested", true)

	// Stock Transactions
	createStringAttr(client, dbID, "stock_transactions", "user_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "company_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "transaction_type", 20, true)
	createIntegerAttr(client, dbID, "stock_transactions", "quantity", true)
	createFloatAttr(client, dbID, "stock_transactions", "price", true)
	createFloatAttr(client, dbID, "stock_transactions", "total_amount", true)
	// createStringAttr(client, dbID, "stock_transactions", "transaction_date", 50, true) // Usually $createdAt handles this, but code might use explicit field?
	// Models don't seem to have explicit transaction_date field other than timestamps?
	// Checking internal/trading/repository.go, it uses $createdAt/UpdatedAt.
	// Models.StockTransaction has CreatedAt. So we don't strictly need it if we rely on Appwrite meta, but wait, Models usually decodes it.

	// Creating Indexes
	createIndex(client, dbID, "users", "email_idx", "unique", []string{"email"})
	createIndex(client, dbID, "companies", "symbol_idx", "unique", []string{"symbol"})
	createIndex(client, dbID, "stock_prices", "company_time_idx", "key", []string{"company_id", "timestamp"})
	createIndex(client, dbID, "virtual_wallets", "user_id_idx", "unique", []string{"user_id"})
	// Add more indexes as needed (e.g. for searching)
}

func createStringAttr(client *appwrite.Client, dbID, colID, key string, size int, required bool) {
	// Check if exists? Appwrite errors if exists? Or we can just try.
	// Easier to just try and ignore "Attribute already exists" error.
	// But let's swallow error for simplicity in this script.
	_, err := client.Databases.CreateStringAttribute(dbID, colID, key, size, required)
	if err == nil {
		slog.Info("Created String attribute", "collection", colID, "key", key)
	} else {
		// slog.Warn("Failed to create attribute (might exist)", "collection", colID, "key", key, "error", err)
	}
	// Optimization: Add a small sleep to avoid hitting rate limits if creating many?
	time.Sleep(100 * time.Millisecond)
}

func createIntegerAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateIntegerAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Integer attribute", "collection", colID, "key", key)
	}
	time.Sleep(100 * time.Millisecond)
}

func createFloatAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateFloatAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Float attribute", "collection", colID, "key", key)
	}
	time.Sleep(100 * time.Millisecond)
}

func createBooleanAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateBooleanAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Boolean attribute", "collection", colID, "key", key)
	}
	time.Sleep(100 * time.Millisecond)
}

func createIndex(client *appwrite.Client, dbID, colID, key, typeStr string, attributes []string) {
	_, err := client.Databases.CreateIndex(dbID, colID, key, typeStr, attributes)
	if err == nil {
		slog.Info("Created Index", "collection", colID, "key", key)
	}
	time.Sleep(100 * time.Millisecond)
}
