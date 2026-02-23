package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	client, err := supabase.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	slog.Info("Setting up Supabase Database (Seeding Data)")
	slog.Info("NOTE: Run the SQL migration in supabase/migrations/001_complete_schema.sql first via Supabase SQL Editor")

	if err := setupCompaniesData(client); err != nil {
		log.Fatalf("Failed to setup companies data: %v", err)
	}
}

var companiesData = []struct {
	symbol      string
	name        string
	sector      string
	marketCap   float64
	description string
	foundedYear int
	employees   int
	totalShares int64
}{

	{"NABIL", "Nabil Bank Limited", "Banking", 180000000000, "Leading private sector bank in Nepal", 1984, 3500, 1800000000},
	{"GBIME", "Global IME Bank Limited", "Banking", 160000000000, "Largest commercial bank by network", 2007, 4200, 1600000000},
	{"NICA", "Nepal Investment Mega Bank", "Banking", 150000000000, "Merged large commercial bank", 1986, 4000, 1500000000},
	{"NMB", "NMB Bank Limited", "Banking", 140000000000, "Corporate-focused commercial bank", 1996, 3200, 1400000000},
	{"HBL", "Himalayan Bank Limited", "Banking", 130000000000, "Joint venture bank with strong presence", 1993, 2800, 1300000000},

	{"NTC", "Nepal Telecom", "Information Technology", 200000000000, "Largest telecom & IT service provider", 2004, 5000, 2000000000},
	{"NITC", "Nepal IT Corporation", "Information Technology", 90000000000, "Government-backed IT service company", 1998, 1500, 900000000},
	{"F1SOFT", "F1Soft International", "Information Technology", 70000000000, "Digital payment and fintech company", 2004, 1200, 700000000},
	{"ESewa", "eSewa Digital Services", "Information Technology", 60000000000, "Fintech wallet operator", 2009, 900, 600000000},
	{"IMS", "IMS Software Solutions", "Information Technology", 30000000000, "IT consulting and development", 2005, 600, 300000000},

	{"HYDRO1", "Upper Tamakoshi Hydropower", "Hydropower", 120000000000, "Large hydropower producer", 2007, 900, 1200000000},
	{"BPC", "Butwal Power Company", "Hydropower", 55000000000, "Integrated power producer", 1966, 600, 550000000},
	{"CHCL", "Chilime Hydropower", "Hydropower", 60000000000, "Hydropower generation", 1995, 450, 600000000},
	{"API", "Api Power Company", "Hydropower", 40000000000, "Hydropower & energy developer", 2003, 350, 400000000},
	{"RSHP", "Rosuwa Shyamkhola Hydro Power", "Hydropower", 35000000000, "Small-medium hydropower", 2010, 280, 350000000},

	{"NLIC", "Nepal Life Insurance Company", "Insurance", 130000000000, "Largest life insurance provider", 2001, 1200, 1300000000},
	{"SICL", "Shikhar Insurance Company", "Insurance", 45000000000, "Leading non-life insurance provider", 2004, 800, 450000000},
	{"NMBHL", "NMB Health Insurance", "Insurance", 35000000000, "Health insurance specialist", 2010, 700, 350000000},

	{"APOLLONP", "Apollo Nepal Hospitals", "Pharma", 50000000000, "Hospital & healthcare chain", 2015, 1500, 500000000},
	{"MHPL", "Medical Health Products Limited", "Pharma", 40000000000, "Pharmaceutical manufacturer", 2008, 900, 400000000},

	{"HDL", "Himalayan Distillery Limited", "Manufacturing", 90000000000, "Distillery & beverages producer", 1985, 700, 900000000},
	{"UNL", "Unilever Nepal Limited", "Manufacturing", 80000000000, "FMCG with global brands", 1992, 300, 800000000},
	{"BNL", "Bottlers Nepal Limited", "Manufacturing", 60000000000, "Bottled beverages producer", 1979, 450, 600000000},
	{"SHIVM", "Shivam Cements Limited", "Manufacturing", 50000000000, "Cement manufacturer", 2003, 1100, 500000000},

	{"DLFNP", "Nepal Housing Development Co.", "Real Estate", 40000000000, "Real estate developer", 2008, 500, 400000000},
}

var eventTypes = []struct {
	eventType   models.MarketEventType
	title       string
	description string
	impact      float64
}{
	{models.MarketEventEarnings, "Q1 Financial Results", "Strong quarterly earnings with revenue growth", 3.5},
	{models.MarketEventDividend, "Annual Dividend Announcement", "Annual dividend payout to shareholders", 2.0},
	{models.MarketEventNews, "New Product Launch", "Company launches innovative product in market", 2.5},
	{models.MarketEventMerger, "Merger Discussion", "Strategic merger discussions with competitor", 4.5},
	{models.MarketEventIPO, "Subsidiary IPO Filing", "Subsidiary company files for IPO", 3.0},
	{models.MarketEventNews, "Market Expansion", "Company expands operations to new region", 2.8},
	{models.MarketEventNews, "Credit Rating Upgrade", "Credit rating agency upgrades company rating", 2.2},
	{models.MarketEventNews, "Government Contract Win", "Company wins major government contract", 3.2},
	{models.MarketEventDividend, "Interim Dividend", "Interim dividend of 50 per share", 1.8},
	{models.MarketEventNews, "Tech Partnership", "Strategic partnership with global technology company", 2.3},
	{models.MarketEventEarnings, "Q2 Financial Results", "Second quarter results show steady growth", 3.0},
	{models.MarketEventNews, "Board Meeting Outcome", "Key decisions from board meeting", 1.5},
	{models.MarketEventDividend, "Bonus Share Announcement", "Company announces bonus shares", 4.0},
	{models.MarketEventNews, "CSR Initiative Launch", "Company launches major CSR program", 1.0},
	{models.MarketEventEarnings, "Q3 Financial Results", "Third quarter earnings exceed expectations", 3.8},
	{models.MarketEventNews, "New Branch Opening", "Company opens new branches across country", 1.2},
	{models.MarketEventMerger, "Acquisition Announcement", "Company announces acquisition of smaller firm", 5.0},
	{models.MarketEventNews, "Regulatory Compliance", "Company achieves new regulatory compliance", 1.5},
	{models.MarketEventEarnings, "Annual Report Release", "Annual financial report shows positive trend", 2.5},
	{models.MarketEventDividend, "Special Dividend", "Special one-time dividend of 100 per share", 5.0},
	{models.MarketEventNews, "Export Agreement", "Major export agreement signed with foreign entity", 3.5},
	{models.MarketEventIPO, "FPO Announcement", "Further public offering announced", 2.8},
	{models.MarketEventNews, "Infrastructure Investment", "Major infrastructure investment announced", 2.0},
	{models.MarketEventNews, "Leadership Change", "New CEO appointed with strong track record", 1.8},
	{models.MarketEventEarnings, "Q4 Financial Results", "Year-end earnings show record profits", 4.2},
}

func setupCompaniesData(client *supabase.Client) error {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Starting Data Setup (Supabase)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\nStep 1: Creating/fetching companies...")

	createdCompanies := make([]models.Company, 0)

	for _, comp := range companiesData {
		mCap := decimal.NewFromFloat(comp.marketCap)

		var existingCompany models.Company
		err := client.ExecuteQueryRow("SELECT * FROM companies WHERE symbol = $1", &existingCompany, comp.symbol)

		if err == nil {

			createdCompanies = append(createdCompanies, existingCompany)
			fmt.Printf("Found existing: %s (%s) - %d total shares\n", existingCompany.Name, existingCompany.Sector, existingCompany.TotalShares)
			continue
		}

		var company models.Company
		err = client.ExecuteInsert(
			"INSERT INTO companies (symbol, name, sector, market_cap, description, founded_year, employees, total_shares, available_shares, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *",
			&company,
			comp.symbol, comp.name, comp.sector, mCap.InexactFloat64(),
			comp.description, comp.foundedYear, comp.employees, comp.totalShares, comp.totalShares, true)
		if err != nil {
			fmt.Printf("Failed to create company %s: %v\n", comp.symbol, err)
			continue
		}

		createdCompanies = append(createdCompanies, company)
		fmt.Printf("Created: %s (%s) - %d total shares\n", company.Name, company.Sector, comp.totalShares)
	}

	if len(createdCompanies) == 0 {
		return fmt.Errorf("failed to create companies")
	}

	fmt.Println("\nStep 2: Creating initial stock price (100) for each company...")
	totalPricesCreated := 0

	for _, company := range createdCompanies {
		initialPrice := decimal.NewFromInt(100)

		err := client.ExecuteInsert(
			"INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, timestamp, timeframe) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *",
			nil,
			company.ID, initialPrice.InexactFloat64(), initialPrice.InexactFloat64(),
			initialPrice.InexactFloat64(), initialPrice.InexactFloat64(),
			0, time.Now().Format(time.RFC3339), "1D")
		if err == nil {
			totalPricesCreated++
			fmt.Printf("Initial price set for %s: 100\n", company.Symbol)
		} else {
			fmt.Printf("Failed to set initial price for %s: %v\n", company.Symbol, err)
		}
	}

	fmt.Println("\nStep 3: Creating 25 market events per company...")
	totalEvents := 0
	for _, company := range createdCompanies {
		for i := 0; i < 25; i++ {
			evType := eventTypes[i%len(eventTypes)]
			impact := evType.impact + (rand.Float64()-0.5)*1.0
			eventDate := time.Now().Add(time.Duration(rand.Intn(365)+1) * 24 * time.Hour)

			err := client.ExecuteInsert(
				"INSERT INTO market_events (company_id, event_type, title, description, impact_percentage, event_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
				nil,
				company.ID, string(evType.eventType), evType.title,
				evType.description, impact, eventDate.Format(time.RFC3339))
			if err == nil {
				totalEvents++
			}
		}
		fmt.Printf("25 events for %s\n", company.Symbol)
	}
	fmt.Printf("Created %d total events\n", totalEvents)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("Setup Complete!\n")
	fmt.Printf("   Companies:      %d\n", len(createdCompanies))
	fmt.Printf("   Initial Prices: %d (all 100)\n", totalPricesCreated)
	fmt.Printf("   Market Events:  %d (25 per company)\n", totalEvents)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nNotes:")
	fmt.Println("   - No test users, transactions, or historical prices created")
	fmt.Println("   - Users register via API, wallet starts with 0 balance")
	fmt.Println("   - Users must top up wallet before buying stocks")
	fmt.Println("   - Stock prices start at 100, prices change via buy/sell trades (order-driven)")
	fmt.Println("   - Each company has total_shares and available_shares")
	fmt.Println("   - Users can only buy shares that are available in the market")
	fmt.Println("   - When buying, available_shares decreases; when selling, it increases")
	fmt.Println("   - Trading volume is tracked on each transaction")

	return nil
}
