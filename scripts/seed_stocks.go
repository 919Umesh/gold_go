package main

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/models"
	"github.com/joho/godotenv"
)

// NepaliCompany represents a demo company for seeding
type NepaliCompany struct {
	Symbol      string
	Name        string
	Sector      string
	MarketCap   float64
	Description string
	FoundedYear int
	Employees   int
	BasePrice   float64
}

var nepaliCompanies = []NepaliCompany{
	// Banking Sector
	{Symbol: "NABIL", Name: "Nabil Bank Limited", Sector: "Banking", MarketCap: 50000000000, Description: "Leading commercial bank in Nepal", FoundedYear: 1984, Employees: 1200, BasePrice: 1200},
	{Symbol: "SCB", Name: "Standard Chartered Bank Nepal", Sector: "Banking", MarketCap: 45000000000, Description: "International banking services", FoundedYear: 1987, Employees: 1000, BasePrice: 850},
	{Symbol: "HBL", Name: "Himalayan Bank Limited", Sector: "Banking", MarketCap: 40000000000, Description: "Joint venture bank with Habib Bank", FoundedYear: 1993, Employees: 950, BasePrice: 650},
	{Symbol: "EBL", Name: "Everest Bank Limited", Sector: "Banking", MarketCap: 38000000000, Description: "Punjab National Bank joint venture", FoundedYear: 1994, Employees: 900, BasePrice: 720},
	{Symbol: "NICA", Name: "NIC Asia Bank", Sector: "Banking", MarketCap: 42000000000, Description: "Merged commercial bank", FoundedYear: 1998, Employees: 1100, BasePrice: 980},

	// Hydropower Sector
	{Symbol: "CHCL", Name: "Chilime Hydropower Company", Sector: "Hydropower", MarketCap: 15000000000, Description: "42.5 MW hydropower plant", FoundedYear: 1996, Employees: 200, BasePrice: 580},
	{Symbol: "NHPC", Name: "Nepal Hydro Power Company", Sector: "Hydropower", MarketCap: 12000000000, Description: "Renewable energy producer", FoundedYear: 2010, Employees: 180, BasePrice: 420},
	{Symbol: "UMHL", Name: "Upper Marsyangdi Hydropower", Sector: "Hydropower", MarketCap: 18000000000, Description: "70 MW hydropower project", FoundedYear: 2008, Employees: 250, BasePrice: 350},
	{Symbol: "RADHI", Name: "Radhi Bidyut Company", Sector: "Hydropower", MarketCap: 10000000000, Description: "Small hydropower developer", FoundedYear: 2012, Employees: 150, BasePrice: 280},

	// Insurance Sector
	{Symbol: "NLIC", Name: "Nepal Life Insurance Company", Sector: "Insurance", MarketCap: 8000000000, Description: "Life insurance provider", FoundedYear: 1988, Employees: 500, BasePrice: 1850},
	{Symbol: "NLICL", Name: "National Life Insurance Company", Sector: "Insurance", MarketCap: 7500000000, Description: "Life and health insurance", FoundedYear: 1998, Employees: 450, BasePrice: 1200},
	{Symbol: "SICL", Name: "Shikhar Insurance Company", Sector: "Insurance", MarketCap: 6000000000, Description: "General insurance services", FoundedYear: 2002, Employees: 300, BasePrice: 950},
	{Symbol: "PRIN", Name: "Prime Insurance Company", Sector: "Insurance", MarketCap: 5500000000, Description: "Non-life insurance", FoundedYear: 2005, Employees: 280, BasePrice: 820},

	// Manufacturing Sector
	{Symbol: "UNL", Name: "Unilever Nepal Limited", Sector: "Manufacturing", MarketCap: 25000000000, Description: "FMCG products manufacturer", FoundedYear: 1992, Employees: 800, BasePrice: 18500},
	{Symbol: "NRIC", Name: "Nepal Reinsurance Company", Sector: "Manufacturing", MarketCap: 4000000000, Description: "Reinsurance services", FoundedYear: 2013, Employees: 120, BasePrice: 1100},
	{Symbol: "BNT", Name: "Bottlers Nepal (Terai)", Sector: "Manufacturing", MarketCap: 12000000000, Description: "Beverage bottling company", FoundedYear: 1985, Employees: 600, BasePrice: 14200},

	// Hotels & Tourism
	{Symbol: "OHL", Name: "Oriental Hotels Limited", Sector: "Hotels", MarketCap: 3500000000, Description: "Luxury hotel chain", FoundedYear: 1970, Employees: 400, BasePrice: 580},
	{Symbol: "TRHPR", Name: "Taragaon Regency Hotel", Sector: "Hotels", MarketCap: 2800000000, Description: "Premium hospitality services", FoundedYear: 1995, Employees: 250, BasePrice: 420},
	{Symbol: "SHL", Name: "Soaltee Hotel Limited", Sector: "Hotels", MarketCap: 4200000000, Description: "Five-star hotel operator", FoundedYear: 1965, Employees: 500, BasePrice: 350},

	// Finance Companies
	{Symbol: "GUFL", Name: "Goodwill Finance Limited", Sector: "Finance", MarketCap: 3000000000, Description: "Financial services provider", FoundedYear: 2007, Employees: 200, BasePrice: 680},
	{Symbol: "CFCL", Name: "Central Finance Company", Sector: "Finance", MarketCap: 2500000000, Description: "Lending and investment", FoundedYear: 2005, Employees: 180, BasePrice: 520},
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	// Initialize config and database
	cfg := config.InitConfig()
	db := config.ConnectDatabase(cfg)

	slog.Info("Starting database seeding...")

	// Seed companies
	slog.Info("Seeding companies...")
	companyIDs := make(map[string]uint)

	for _, nc := range nepaliCompanies {
		company := &models.Company{
			Symbol:      nc.Symbol,
			Name:        nc.Name,
			Sector:      nc.Sector,
			MarketCap:   nc.MarketCap,
			Description: nc.Description,
			FoundedYear: nc.FoundedYear,
			Employees:   nc.Employees,
			IsActive:    true,
		}

		// Check if company already exists
		var existing models.Company
		err := db.Where("symbol = ?", nc.Symbol).First(&existing).Error
		if err == nil {
			slog.Info("Company already exists", "symbol", nc.Symbol)
			companyIDs[nc.Symbol] = existing.ID
			continue
		}

		if err := db.Create(company).Error; err != nil {
			log.Printf("Failed to create company %s: %v", nc.Symbol, err)
			continue
		}

		companyIDs[nc.Symbol] = company.ID
		slog.Info("Created company", "symbol", nc.Symbol, "id", company.ID)
	}

	// Generate historical price data (1 year)
	slog.Info("Generating historical price data...")

	for _, nc := range nepaliCompanies {
		companyID, exists := companyIDs[nc.Symbol]
		if !exists {
			continue
		}

		// Check if historical data already exists
		var count int64
		db.Model(&models.StockPrice{}).Where("company_id = ?", companyID).Count(&count)
		if count > 0 {
			slog.Info("Historical data already exists", "symbol", nc.Symbol)
			continue
		}

		// Generate 1 year of daily data
		startDate := time.Now().AddDate(-1, 0, 0)
		currentPrice := nc.BasePrice

		for days := 0; days < 365; days++ {
			timestamp := startDate.AddDate(0, 0, days)

			// Skip weekends (Friday and Saturday in Nepal)
			if timestamp.Weekday() == time.Friday || timestamp.Weekday() == time.Saturday {
				continue
			}

			// Simulate price movement using random walk
			change := (rand.Float64() - 0.5) * 0.05 // -2.5% to +2.5%
			currentPrice = currentPrice * (1 + change)

			// Generate OHLC data
			open := currentPrice * (1 + (rand.Float64()-0.5)*0.01)
			high := currentPrice * (1 + rand.Float64()*0.02)
			low := currentPrice * (1 - rand.Float64()*0.02)
			close := currentPrice

			// Ensure high is highest and low is lowest
			high = math.Max(high, math.Max(open, close))
			low = math.Min(low, math.Min(open, close))

			volume := int64(10000 + rand.Intn(990000))

			stockPrice := &models.StockPrice{
				CompanyID:  companyID,
				OpenPrice:  math.Round(open*100) / 100,
				HighPrice:  math.Round(high*100) / 100,
				LowPrice:   math.Round(low*100) / 100,
				ClosePrice: math.Round(close*100) / 100,
				Volume:     volume,
				Timestamp:  timestamp,
				Timeframe:  "1d",
			}

			if err := db.Create(stockPrice).Error; err != nil {
				log.Printf("Failed to create price for %s: %v", nc.Symbol, err)
			}
		}

		slog.Info("Generated historical data", "symbol", nc.Symbol, "days", 365)
	}

	// Generate some market events
	slog.Info("Generating market events...")

	eventTypes := []models.MarketEventType{
		models.MarketEventEarnings,
		models.MarketEventNews,
		models.MarketEventDividend,
	}

	for i := 0; i < 20; i++ {
		nc := nepaliCompanies[rand.Intn(len(nepaliCompanies))]
		companyID := companyIDs[nc.Symbol]

		event := &models.MarketEvent{
			CompanyID:        companyID,
			EventType:        eventTypes[rand.Intn(len(eventTypes))],
			Title:            fmt.Sprintf("Important announcement for %s", nc.Name),
			Description:      "This is a simulated market event for educational purposes",
			ImpactPercentage: (rand.Float64() - 0.5) * 10,
			EventDate:        time.Now().Add(time.Duration(rand.Intn(60)) * 24 * time.Hour),
		}

		if err := db.Create(event).Error; err != nil {
			log.Printf("Failed to create event: %v", err)
		}
	}

	slog.Info("Database seeding completed successfully!")
	slog.Info("Summary", "companies", len(nepaliCompanies), "historical_days", 365)
}
