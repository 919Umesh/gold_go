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
	slog.Info("Setting up Appwrite Database (Fresh Setup)", "database_id", dbID)

	// List of collections to create with permissions
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
		{ID: "transactions", Name: "Wallet Transactions"},
		{ID: "stock_transactions", Name: "Stock Transactions"},
	}

	// Permissions: Create, Read, Update for any user (NO Delete)
	permissions := []string{
		"create(\"users\")",
		"read(\"users\")",
		"update(\"users\")",
	}

	for _, c := range collections {
		createCollectionWithPermissions(client, dbID, c.ID, c.Name, permissions)
	}

	// Create Attributes
	createAttributes(client, dbID)

	slog.Info("✅ Appwrite setup completed successfully!")
	slog.Info("Collections created with permissions: Create, Read, Update (NO Delete)")
	slog.Info("New fields added: profile (users), image_url (market_events)")
}

func createCollectionWithPermissions(client *appwrite.Client, dbID, colID, colName string, permissions []string) {
	_, err := client.Databases.GetCollection(dbID, colID)
	if err == nil {
		slog.Info("Collection already exists", "id", colID)
		return
	}

	_, err = client.Databases.CreateCollection(
		dbID,
		colID,
		colName,
		appwrite.WithCreateCollectionPermissions(permissions),
		appwrite.WithCreateCollectionDocumentSecurity(true),
	)
	if err != nil {
		log.Printf("Failed to create collection %s: %v", colName, err)
	} else {
		slog.Info("✅ Created collection with permissions", "id", colID)
	}
	time.Sleep(200 * time.Millisecond)
}

func createAttributes(client *appwrite.Client, dbID string) {
	// ==================== USERS ====================
	createStringAttr(client, dbID, "users", "full_name", 100, true)
	createStringAttr(client, dbID, "users", "email", 255, true)
	createStringAttr(client, dbID, "users", "phone", 20, true)
	createStringAttr(client, dbID, "users", "password_hash", 255, true)
	createStringAttr(client, dbID, "users", "kyc_status", 50, false)
	createStringAttr(client, dbID, "users", "role", 50, false)
	// NEW: Profile field (optional) - URL to profile image
	createStringAttr(client, dbID, "users", "profile", 500, false)

	// ==================== COMPANIES ====================
	createStringAttr(client, dbID, "companies", "symbol", 10, true)
	createStringAttr(client, dbID, "companies", "name", 100, true)
	createStringAttr(client, dbID, "companies", "sector", 50, true)
	createFloatAttr(client, dbID, "companies", "market_cap", true)
	createStringAttr(client, dbID, "companies", "description", 1000, false)
	createIntegerAttr(client, dbID, "companies", "founded_year", false)
	createIntegerAttr(client, dbID, "companies", "employees", false)
	createBooleanAttr(client, dbID, "companies", "is_active", true)

	// ==================== STOCK PRICES ====================
	createStringAttr(client, dbID, "stock_prices", "company_id", 50, true)
	createFloatAttr(client, dbID, "stock_prices", "open_price", true)
	createFloatAttr(client, dbID, "stock_prices", "high_price", true)
	createFloatAttr(client, dbID, "stock_prices", "low_price", true)
	createFloatAttr(client, dbID, "stock_prices", "close_price", true)
	createIntegerAttr(client, dbID, "stock_prices", "volume", true)
	createStringAttr(client, dbID, "stock_prices", "timestamp", 50, true)
	createStringAttr(client, dbID, "stock_prices", "timeframe", 10, true)

	// ==================== MARKET EVENTS ====================
	createStringAttr(client, dbID, "market_events", "company_id", 50, true)
	createStringAttr(client, dbID, "market_events", "event_type", 50, true)
	createStringAttr(client, dbID, "market_events", "title", 255, true)
	createStringAttr(client, dbID, "market_events", "description", 1000, false)
	createFloatAttr(client, dbID, "market_events", "impact_percentage", true)
	createStringAttr(client, dbID, "market_events", "event_date", 50, true)
	// NEW: Image URL for market events (optional)
	createStringAttr(client, dbID, "market_events", "image_url", 500, false)

	// ==================== VIRTUAL WALLETS ====================
	createStringAttr(client, dbID, "virtual_wallets", "user_id", 50, true)
	createFloatAttr(client, dbID, "virtual_wallets", "balance", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_invested", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_profit_loss", true)
	createFloatAttr(client, dbID, "virtual_wallets", "fiat_balance", false)
	createBooleanAttr(client, dbID, "virtual_wallets", "locked", false)
	createIntegerAttr(client, dbID, "virtual_wallets", "version", false)

	// ==================== WALLET TRANSACTIONS (NO company_id!) ====================
	// This is for fiat top-ups only - NO company_id required
	createStringAttr(client, dbID, "transactions", "user_id", 50, true)
	createStringAttr(client, dbID, "transactions", "type", 50, true)
	createFloatAttr(client, dbID, "transactions", "amount", true)
	createStringAttr(client, dbID, "transactions", "status", 50, true)
	createStringAttr(client, dbID, "transactions", "reference_id", 100, true)

	// ==================== USER PORTFOLIOS ====================
	createStringAttr(client, dbID, "user_portfolios", "user_id", 50, true)
	createStringAttr(client, dbID, "user_portfolios", "company_id", 50, true)
	createIntegerAttr(client, dbID, "user_portfolios", "quantity", true)
	createFloatAttr(client, dbID, "user_portfolios", "average_price", true)
	createFloatAttr(client, dbID, "user_portfolios", "total_invested", true)

	// ==================== STOCK TRANSACTIONS (HAS company_id!) ====================
	// This is for buy/sell stock trades - company_id IS required
	createStringAttr(client, dbID, "stock_transactions", "user_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "company_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "type", 20, true)
	createIntegerAttr(client, dbID, "stock_transactions", "quantity", true)
	createFloatAttr(client, dbID, "stock_transactions", "price_per_share", true)
	createFloatAttr(client, dbID, "stock_transactions", "total_amount", true)
	createStringAttr(client, dbID, "stock_transactions", "status", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "reference_id", 100, false)

	// ==================== INDEXES ====================
	slog.Info("Creating indexes...")
	createIndex(client, dbID, "users", "email_idx", "unique", []string{"email"})
	createIndex(client, dbID, "companies", "symbol_idx", "unique", []string{"symbol"})
	createIndex(client, dbID, "stock_prices", "company_time_idx", "key", []string{"company_id", "timestamp"})
	createIndex(client, dbID, "virtual_wallets", "user_id_idx", "unique", []string{"user_id"})
	createIndex(client, dbID, "transactions", "user_id_idx", "key", []string{"user_id"})
	createIndex(client, dbID, "stock_transactions", "user_id_idx", "key", []string{"user_id"})
	createIndex(client, dbID, "user_portfolios", "user_company_idx", "key", []string{"user_id", "company_id"})
	createIndex(client, dbID, "market_events", "company_id_idx", "key", []string{"company_id"})
}

func createStringAttr(client *appwrite.Client, dbID, colID, key string, size int, required bool) {
	_, err := client.Databases.CreateStringAttribute(dbID, colID, key, size, required)
	if err == nil {
		slog.Info("Created String attribute", "collection", colID, "key", key)
	}
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
