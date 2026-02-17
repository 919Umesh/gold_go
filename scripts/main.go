package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
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

	// 1. DELETE EXISTING COLLECTIONS
	slog.Info("🗑️  Deleting existing collections...")
	for _, c := range collections {
		_, err := client.Databases.DeleteCollection(dbID, c.ID)
		if err == nil {
			slog.Info("Deleted collection", "id", c.ID)
		} else {
			// Ignore error if collection doesn't exist
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 2. CREATE COLLECTIONS
	// Permissions: Create, Read, Update for any user (NO Delete)
	permissions := []string{
		"create(\"users\")",
		"read(\"users\")",
		"update(\"users\")",
	}

	for _, c := range collections {
		createCollectionWithPermissions(client, dbID, c.ID, c.Name, permissions)
	}

	// 3. CREATE ATTRIBUTES
	createAttributes(client, dbID)

	// 4. CREATE STORAGE BUCKET
	createStorageBucket(client)

	slog.Info("✅ Appwrite schema setup completed!")

	// 5. SEED DATA
	// Setup initial data: 25 companies with events and transactions
	if err := setupCompaniesData(client); err != nil {
		log.Fatalf("Failed to setup companies data: %v", err)
	}
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
	// NEW: Profile field (optional) - URL or File ID
	createStringAttr(client, dbID, "users", "profile_image_id", 100, false)

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
	createStringAttr(client, dbID, "market_events", "image_url", 500, false)

	// ==================== VIRTUAL WALLETS ====================
	createStringAttr(client, dbID, "virtual_wallets", "user_id", 50, true)
	createFloatAttr(client, dbID, "virtual_wallets", "balance", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_invested", true)
	createFloatAttr(client, dbID, "virtual_wallets", "total_profit_loss", true)
	createFloatAttr(client, dbID, "virtual_wallets", "fiat_balance", false)
	createBooleanAttr(client, dbID, "virtual_wallets", "locked", false)
	createIntegerAttr(client, dbID, "virtual_wallets", "version", false)

	// ==================== WALLET TRANSACTIONS ====================
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

	// ==================== STOCK TRANSACTIONS ====================
	createStringAttr(client, dbID, "stock_transactions", "user_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "company_id", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "type", 20, true)
	createIntegerAttr(client, dbID, "stock_transactions", "quantity", true)
	createFloatAttr(client, dbID, "stock_transactions", "price_per_share", true)
	createFloatAttr(client, dbID, "stock_transactions", "total_amount", true)
	createStringAttr(client, dbID, "stock_transactions", "status", 50, true)
	createStringAttr(client, dbID, "stock_transactions", "reference_id", 100, false)
	createStringAttr(client, dbID, "stock_transactions", "created_at", 50, false)

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

func createStorageBucket(client *appwrite.Client) {
	bucketID := "user-profiles"
	name := "User Profiles"
	allowedFileExtensions := []string{"jpg", "jpeg", "png", "gif"}

	_, err := client.Storage.GetBucket(bucketID)
	if err == nil {
		slog.Info("Storage bucket already exists", "id", bucketID)
		return
	}

	_, err = client.Storage.CreateBucket(
		bucketID,
		name,
		client.Storage.WithCreateBucketPermissions([]string{"read(\"any\")", "write(\"users\")"}), // Public read, users write
		client.Storage.WithCreateBucketFileSecurity(true),
		client.Storage.WithCreateBucketEnabled(true),
		client.Storage.WithCreateBucketMaximumFileSize(5*1024*1024), // 5MB
		client.Storage.WithCreateBucketAllowedFileExtensions(allowedFileExtensions),
	)

	if err != nil {
		log.Printf("Failed to create storage bucket %s: %v", name, err)
	} else {
		slog.Info("✅ Created storage bucket", "id", bucketID)
	}
}

// Helpers for creation
func createStringAttr(client *appwrite.Client, dbID, colID, key string, size int, required bool) {
	_, err := client.Databases.CreateStringAttribute(dbID, colID, key, size, required)
	if err == nil {
		slog.Info("Created String attribute", "collection", colID, "key", key)
	}
	time.Sleep(50 * time.Millisecond)
}
func createIntegerAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateIntegerAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Integer attribute", "collection", colID, "key", key)
	}
	time.Sleep(50 * time.Millisecond)
}
func createFloatAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateFloatAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Float attribute", "collection", colID, "key", key)
	}
	time.Sleep(50 * time.Millisecond)
}
func createBooleanAttr(client *appwrite.Client, dbID, colID, key string, required bool) {
	_, err := client.Databases.CreateBooleanAttribute(dbID, colID, key, required)
	if err == nil {
		slog.Info("Created Boolean attribute", "collection", colID, "key", key)
	}
	time.Sleep(50 * time.Millisecond)
}
func createIndex(client *appwrite.Client, dbID, colID, key, typeStr string, attributes []string) {
	_, err := client.Databases.CreateIndex(dbID, colID, key, typeStr, attributes)
	if err == nil {
		slog.Info("Created Index", "collection", colID, "key", key)
	}
	time.Sleep(50 * time.Millisecond)
}

// ===============================================
// DATA SETUP
// ===============================================

// Company data for 25 companies
var companiesData = []struct {
	symbol      string
	name        string
	sector      string
	marketCap   float64
	description string
	foundedYear int
	employees   int
}{
	// Banking
	{"NABIL", "Nabil Bank Limited", "Banking", 180000000000, "Leading private sector bank in Nepal", 1984, 3500},
	{"GBIME", "Global IME Bank Limited", "Banking", 160000000000, "Largest commercial bank by network", 2007, 4200},
	{"NICA", "Nepal Investment Mega Bank", "Banking", 150000000000, "Merged large commercial bank", 1986, 4000},
	{"NMB", "NMB Bank Limited", "Banking", 140000000000, "Corporate-focused commercial bank", 1996, 3200},
	{"HBL", "Himalayan Bank Limited", "Banking", 130000000000, "Joint venture bank with strong presence", 1993, 2800},
	// IT
	{"NTC", "Nepal Telecom", "Information Technology", 200000000000, "Largest telecom & IT service provider", 2004, 5000},
	{"NITC", "Nepal IT Corporation", "Information Technology", 90000000000, "Government-backed IT service company", 1998, 1500},
	{"F1SOFT", "F1Soft International", "Information Technology", 70000000000, "Digital payment and fintech company", 2004, 1200},
	{"ESewa", "eSewa Digital Services", "Information Technology", 60000000000, "Fintech wallet operator", 2009, 900},
	{"IMS", "IMS Software Solutions", "Information Technology", 30000000000, "IT consulting and development", 2005, 600},
	// Hydropower
	{"HYDRO1", "Upper Tamakoshi Hydropower", "Hydropower", 120000000000, "Large hydropower producer", 2007, 900},
	{"BPC", "Butwal Power Company", "Hydropower", 55000000000, "Integrated power producer", 1966, 600},
	{"CHCL", "Chilime Hydropower", "Hydropower", 60000000000, "Hydropower generation", 1995, 450},
	{"API", "Api Power Company", "Hydropower", 40000000000, "Hydropower & energy developer", 2003, 350},
	{"RSHP", "Rosuwa Shyamkhola Hydro Power", "Hydropower", 35000000000, "Small-medium hydropower", 2010, 280},
	// Insurance
	{"NLIC", "Nepal Life Insurance Company", "Insurance", 130000000000, "Largest life insurance provider", 2001, 1200},
	{"SICL", "Shikhar Insurance Company", "Insurance", 45000000000, "Leading non-life insurance provider", 2004, 800},
	{"NMBHL", "NMB Health Insurance", "Insurance", 35000000000, "Health insurance specialist", 2010, 700},
	// Pharma
	{"APOLLONP", "Apollo Nepal Hospitals", "Pharma", 50000000000, "Hospital & healthcare chain", 2015, 1500},
	{"MHPL", "Medical Health Products Limited", "Pharma", 40000000000, "Pharmaceutical manufacturer", 2008, 900},
	// Manufacturing
	{"HDL", "Himalayan Distillery Limited", "Manufacturing", 90000000000, "Distillery & beverages producer", 1985, 700},
	{"UNL", "Unilever Nepal Limited", "Manufacturing", 80000000000, "FMCG with global brands", 1992, 300},
	{"BNL", "Bottlers Nepal Limited", "Manufacturing", 60000000000, "Bottled beverages producer", 1979, 450},
	{"SHIVM", "Shivam Cements Limited", "Manufacturing", 50000000000, "Cement manufacturer", 2003, 1100},
	// Real Estate
	{"DLFNP", "Nepal Housing Development Co.", "Real Estate", 40000000000, "Real estate developer", 2008, 500},
}

var eventTypes = []struct {
	eventType   models.MarketEventType
	title       string
	description string
	impact      float64
}{
	{models.MarketEventEarnings, "Q3 Financial Results", "Strong quarterly earnings announcement with profit growth", 3.5},
	{models.MarketEventDividend, "Dividend Announcement", "50% dividend payout to shareholders", 2.0},
	{models.MarketEventNews, "New Product Launch", "Company launches innovative product in market", 2.5},
	{models.MarketEventMerger, "Merger Announcement", "Strategic merger with competitor announced", 4.5},
	{models.MarketEventIPO, "Subsidiary IPO", "Subsidiary company to be listed on stock exchange", 3.0},
	{models.MarketEventNews, "Market Expansion", "Company expands operations to new region", 2.8},
	{models.MarketEventNews, "Credit Rating Upgrade", "Credit rating agency upgrades company rating", 2.2},
	{models.MarketEventNews, "Contract Win", "Company wins major government contract", 3.2},
	{models.MarketEventDividend, "Special Dividend", "Special one-time dividend of ₹100 per share", 5.0},
	{models.MarketEventNews, "Tech Partnership", "Strategic partnership with global technology company", 2.3},
}

var transactionTypes = []models.StockTransactionType{models.StockTransactionBuy, models.StockTransactionSell}
var transactionStatuses = []models.StockTransactionStatus{models.StockTransactionCompleted, models.StockTransactionCompleted, models.StockTransactionCompleted, models.StockTransactionPending}

func setupCompaniesData(client *appwrite.Client) error {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting Data Setup")
	fmt.Println(strings.Repeat("=", 80))

	// Step 1: Create 25 Companies
	fmt.Println("\n📦 Step 1: Creating 25 companies...")
	createdCompanies := make([]models.Company, 0)

	for _, comp := range companiesData {
		// Use Decimal for creating model struct (verification mainly)
		mCap := decimal.NewFromFloat(comp.marketCap)

		company := models.Company{
			Symbol:      comp.symbol,
			Name:        comp.name,
			Sector:      comp.sector,
			MarketCap:   mCap,
			Description: comp.description,
			FoundedYear: comp.foundedYear,
			Employees:   comp.employees,
			IsActive:    true,
		}

		// Prepare map for Appwrite (decimals converted to float64)
		data := map[string]interface{}{
			"symbol":       company.Symbol,
			"name":         company.Name,
			"sector":       company.Sector,
			"market_cap":   mCap.InexactFloat64(),
			"description":  company.Description,
			"founded_year": company.FoundedYear,
			"employees":    company.Employees,
			"is_active":    company.IsActive,
		}

		doc, err := client.Databases.CreateDocument(client.Config.DatabaseID, "companies", id.Unique(), data)
		if err != nil {
			fmt.Printf("❌ Failed to create company %s: %v\n", comp.symbol, err)
			continue
		}

		// Store ID back in company struct
		company.ID = doc.Id
		createdCompanies = append(createdCompanies, company)
		fmt.Printf("✅ Created: %s (%s)\n", company.Name, company.Sector)
	}

	if len(createdCompanies) == 0 {
		return fmt.Errorf("failed to create companies")
	}

	// Step 2: Create Stock Prices (30 days for ML)
	fmt.Println("\n📊 Step 2: Creating 30 days of stock prices...")
	totalPricesCreated := 0

	for _, company := range createdCompanies {
		// Start with random base price
		basePriceStr := decimal.NewFromFloat(100.0 + rand.Float64()*500.0)

		for i := 29; i >= 0; i-- {
			timestamp := time.Now().AddDate(0, 0, -i)

			// Fluctuation
			changePercent := (rand.Float64() - 0.5) * 0.05 // +/- 2.5%
			changeFactor := decimal.NewFromFloat(1.0 + changePercent)

			open := basePriceStr
			closePrice := open.Mul(changeFactor)

			// Limit high/low
			// high := max(open, closePrice) * (1 + rand*0.01)
			high := decimal.Max(open, closePrice).Mul(decimal.NewFromFloat(1.0 + rand.Float64()*0.01))
			// low := min(open, closePrice) * (1 - rand*0.01)
			low := decimal.Min(open, closePrice).Mul(decimal.NewFromFloat(1.0 - rand.Float64()*0.01))

			volume := int64(rand.Intn(5000000)) + 1000000

			priceData := map[string]interface{}{
				"company_id":  company.ID,
				"open_price":  open.InexactFloat64(),
				"high_price":  high.InexactFloat64(),
				"low_price":   low.InexactFloat64(),
				"close_price": closePrice.InexactFloat64(),
				"volume":      volume,
				"timestamp":   timestamp.Format(time.RFC3339),
				"timeframe":   "1D",
			}

			_, err := client.Databases.CreateDocument(client.Config.DatabaseID, "stock_prices", id.Unique(), priceData)
			if err == nil {
				totalPricesCreated++
			}

			basePriceStr = closePrice
		}
		fmt.Printf("✅ Prices for %s\n", company.Symbol)
	}

	// Step 3: Market Events (10 per company = 250 total)
	fmt.Println("\n📰 Step 3: Creating market events...")
	totalEvents := 0
	for _, company := range createdCompanies {
		for i := 0; i < 10; i++ {
			evType := eventTypes[i%len(eventTypes)]
			impact := decimal.NewFromFloat(evType.impact + (rand.Float64() - 0.5))
			eventDate := time.Now().AddDate(0, 0, -(i+1)*10)

			eventData := map[string]interface{}{
				"company_id":        company.ID,
				"event_type":        string(evType.eventType),
				"title":             evType.title,
				"description":       evType.description,
				"impact_percentage": impact.InexactFloat64(),
				"event_date":        eventDate.Format(time.RFC3339),
				"image_url":         fmt.Sprintf("https://placehold.co/600x400?text=%s", strings.ReplaceAll(evType.title, " ", "+")),
			}
			_, err := client.Databases.CreateDocument(client.Config.DatabaseID, "market_events", id.Unique(), eventData)
			if err == nil {
				totalEvents++
			}
		}
	}
	fmt.Printf("✅ Created %d events\n", totalEvents)

	// Step 4: Transactions (30 per company = 750 total)
	fmt.Println("\n💰 Step 4: Creating stock transactions (target 30/company)...")
	totalTx := 0
	testUsers := []string{"user_001", "user_002", "user_003"}

	for _, company := range createdCompanies {
		for i := 0; i < 30; i++ {
			userID := testUsers[rand.Intn(len(testUsers))]
			txType := transactionTypes[rand.Intn(len(transactionTypes))]
			status := transactionStatuses[rand.Intn(len(transactionStatuses))]

			qty := rand.Intn(100) + 1
			price := decimal.NewFromFloat(100.0 + rand.Float64()*500.0).Round(2)
			total := price.Mul(decimal.NewFromInt(int64(qty)))

			txData := map[string]interface{}{
				"user_id":         userID,
				"company_id":      company.ID,
				"type":            string(txType),
				"quantity":        qty,
				"price_per_share": price.InexactFloat64(),
				"total_amount":    total.InexactFloat64(),
				"status":          string(status),
				"reference_id":    fmt.Sprintf("TXN_%s_%d", company.Symbol, time.Now().UnixNano()),
				"created_at":      time.Now().Add(time.Duration(-i) * time.Hour).Format(time.RFC3339),
			}
			_, err := client.Databases.CreateDocument(client.Config.DatabaseID, "stock_transactions", id.Unique(), txData)
			if err == nil {
				totalTx++
			}
		}
		if totalTx%1000 == 0 {
			fmt.Printf("... %d transactions\n", totalTx)
		}
	}
	fmt.Printf("✅ Created %d total transactions\n", totalTx)

	return nil
}
