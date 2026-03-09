package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

type SeedCompany struct {
	Symbol            string
	Name              string
	Sector            string
	TotalSupply       int64
	CurrentPrice      float64
	SharesOutstanding int64
	EPS               float64
	PERatio           float64
	BookValue         float64
	PBV               float64
	Week52High        float64
	Week52Low         float64
	Avg120Day         float64
	Yield1Year        float64
	ListedDate        string
	Description       string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	client, err := supabase.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}

	companies := []SeedCompany{
		{Symbol: "NABIL", Name: "Nabil Bank Limited", Sector: "Commercial Banks", TotalSupply: 10000000, CurrentPrice: 1250.00, SharesOutstanding: 9500000, EPS: 45.23, PERatio: 27.63, BookValue: 312.50, PBV: 4.0, Week52High: 1400.00, Week52Low: 980.00, Avg120Day: 1180.00, Yield1Year: 8.5, ListedDate: "2001-03-15", Description: "One of the largest commercial banks in Nepal with extensive banking services."},
		{Symbol: "GBIME", Name: "Global IME Bank Limited", Sector: "Commercial Banks", TotalSupply: 12000000, CurrentPrice: 310.00, SharesOutstanding: 11500000, EPS: 18.75, PERatio: 16.53, BookValue: 195.60, PBV: 1.59, Week52High: 380.00, Week52Low: 250.00, Avg120Day: 305.00, Yield1Year: 5.2, ListedDate: "2007-06-20", Description: "A leading commercial bank formed by merger of multiple banks."},
		{Symbol: "NICA", Name: "NIC Asia Bank Limited", Sector: "Commercial Banks", TotalSupply: 14000000, CurrentPrice: 620.00, SharesOutstanding: 13000000, EPS: 28.40, PERatio: 21.83, BookValue: 230.00, PBV: 2.70, Week52High: 720.00, Week52Low: 500.00, Avg120Day: 600.00, Yield1Year: 6.8, ListedDate: "2005-08-12", Description: "One of the largest banks in Nepal with wide branch network."},
		{Symbol: "NMB", Name: "NMB Bank Limited", Sector: "Commercial Banks", TotalSupply: 8000000, CurrentPrice: 440.00, SharesOutstanding: 7500000, EPS: 22.10, PERatio: 19.91, BookValue: 210.00, PBV: 2.10, Week52High: 520.00, Week52Low: 350.00, Avg120Day: 420.00, Yield1Year: 4.5, ListedDate: "2008-01-10", Description: "A major commercial bank providing comprehensive banking solutions."},
		{Symbol: "HBL", Name: "Himalayan Bank Limited", Sector: "Commercial Banks", TotalSupply: 9000000, CurrentPrice: 1480.00, SharesOutstanding: 8500000, EPS: 52.30, PERatio: 28.30, BookValue: 350.00, PBV: 4.23, Week52High: 1600.00, Week52Low: 1100.00, Avg120Day: 1400.00, Yield1Year: 9.2, ListedDate: "1998-07-05", Description: "Established commercial bank with a strong track record."},
		{Symbol: "NTC", Name: "Nepal Telecom", Sector: "Telecom", TotalSupply: 15000000, CurrentPrice: 880.00, SharesOutstanding: 14000000, EPS: 40.50, PERatio: 21.73, BookValue: 420.00, PBV: 2.10, Week52High: 1020.00, Week52Low: 720.00, Avg120Day: 860.00, Yield1Year: 7.1, ListedDate: "2003-11-18", Description: "The largest telecom provider in Nepal."},
		{Symbol: "SHIVM", Name: "Shivam Cements Limited", Sector: "Manufacturing", TotalSupply: 5000000, CurrentPrice: 560.00, SharesOutstanding: 4800000, EPS: 30.10, PERatio: 18.60, BookValue: 280.00, PBV: 2.00, Week52High: 650.00, Week52Low: 420.00, Avg120Day: 530.00, Yield1Year: 3.8, ListedDate: "2015-04-22", Description: "A leading cement manufacturer in Nepal."},
		{Symbol: "NLIC", Name: "Nepal Life Insurance Company", Sector: "Life Insurance", TotalSupply: 7000000, CurrentPrice: 1100.00, SharesOutstanding: 6500000, EPS: 48.20, PERatio: 22.82, BookValue: 310.00, PBV: 3.55, Week52High: 1280.00, Week52Low: 880.00, Avg120Day: 1050.00, Yield1Year: 6.5, ListedDate: "2006-02-14", Description: "The largest life insurance company in Nepal."},
		{Symbol: "SICL", Name: "Shikhar Insurance Company Limited", Sector: "Non-Life Insurance", TotalSupply: 3000000, CurrentPrice: 780.00, SharesOutstanding: 2800000, EPS: 35.60, PERatio: 21.91, BookValue: 250.00, PBV: 3.12, Week52High: 900.00, Week52Low: 620.00, Avg120Day: 750.00, Yield1Year: 5.0, ListedDate: "2010-09-08", Description: "A leading non-life insurance company."},
		{Symbol: "HDL", Name: "Himalayan Distillery Limited", Sector: "Manufacturing", TotalSupply: 2000000, CurrentPrice: 1650.00, SharesOutstanding: 1900000, EPS: 62.40, PERatio: 26.44, BookValue: 480.00, PBV: 3.44, Week52High: 1850.00, Week52Low: 1300.00, Avg120Day: 1580.00, Yield1Year: 10.5, ListedDate: "2004-05-30", Description: "One of the largest distillery companies in Nepal."},
		{Symbol: "UNL", Name: "Unilever Nepal Limited", Sector: "Manufacturing", TotalSupply: 1500000, CurrentPrice: 14200.00, SharesOutstanding: 1400000, EPS: 310.50, PERatio: 45.73, BookValue: 1200.00, PBV: 11.83, Week52High: 16000.00, Week52Low: 11500.00, Avg120Day: 13800.00, Yield1Year: 12.0, ListedDate: "1995-01-20", Description: "The Nepal subsidiary of Unilever, a major FMCG company."},
		{Symbol: "CHCL", Name: "Chilime Hydropower Company", Sector: "Hydropower", TotalSupply: 6000000, CurrentPrice: 520.00, SharesOutstanding: 5500000, EPS: 24.30, PERatio: 21.40, BookValue: 200.00, PBV: 2.60, Week52High: 620.00, Week52Low: 400.00, Avg120Day: 500.00, Yield1Year: 4.2, ListedDate: "2011-07-15", Description: "A major hydropower company operating in Nepal."},
		{Symbol: "NHPC", Name: "Nepal Hydro Power Company", Sector: "Hydropower", TotalSupply: 4000000, CurrentPrice: 380.00, SharesOutstanding: 3800000, EPS: 15.80, PERatio: 24.05, BookValue: 160.00, PBV: 2.38, Week52High: 460.00, Week52Low: 300.00, Avg120Day: 370.00, Yield1Year: 3.5, ListedDate: "2012-11-25", Description: "A hydropower company contributing to Nepal energy sector."},
		{Symbol: "BPC", Name: "Butwal Power Company", Sector: "Hydropower", TotalSupply: 8000000, CurrentPrice: 420.00, SharesOutstanding: 7500000, EPS: 19.50, PERatio: 21.54, BookValue: 185.00, PBV: 2.27, Week52High: 500.00, Week52Low: 340.00, Avg120Day: 410.00, Yield1Year: 4.8, ListedDate: "2003-03-12", Description: "One of the pioneer hydropower companies of Nepal."},
		{Symbol: "NIFRA", Name: "Nepal Infrastructure Bank", Sector: "Infrastructure Development Bank", TotalSupply: 10000000, CurrentPrice: 220.00, SharesOutstanding: 9500000, EPS: 8.60, PERatio: 25.58, BookValue: 120.00, PBV: 1.83, Week52High: 280.00, Week52Low: 170.00, Avg120Day: 215.00, Yield1Year: 2.5, ListedDate: "2019-12-01", Description: "Nepal's first infrastructure development bank."},
	}

	eventTypes := []string{"agm", "dividend", "bonus_share", "right_share", "quarterly_report", "board_meeting", "financial_results", "stock_split", "merger_acquisition", "ipo_announcement"}
	eventTitles := map[string]string{
		"agm":                "Annual General Meeting",
		"dividend":           "Cash Dividend Distribution",
		"bonus_share":        "Bonus Share Issuance",
		"right_share":        "Right Share Offering",
		"quarterly_report":   "Quarterly Financial Report",
		"board_meeting":      "Board of Directors Meeting",
		"financial_results":  "Annual Financial Results",
		"stock_split":        "Stock Split Announcement",
		"merger_acquisition": "Merger/Acquisition Update",
		"ipo_announcement":   "IPO Announcement",
	}
	eventDescs := map[string]string{
		"agm":                "Annual General Meeting of shareholders to discuss financial performance and future plans.",
		"dividend":           "Distribution of cash dividend to shareholders for the fiscal year.",
		"bonus_share":        "Issuance of bonus shares to existing shareholders based on holdings.",
		"right_share":        "Right share offering to existing shareholders at a discounted price.",
		"quarterly_report":   "Publication of quarterly financial report with revenue and profit details.",
		"board_meeting":      "Board of Directors meeting to discuss strategic decisions.",
		"financial_results":  "Annual financial results publication with detailed P&L and balance sheet.",
		"stock_split":        "Stock split to increase liquidity and make shares more accessible.",
		"merger_acquisition": "Update on planned merger or acquisition activities.",
		"ipo_announcement":   "Announcement regarding upcoming IPO or FPO activities.",
	}

	fmt.Println("=== Stock Market Simulator Seeder ===")
	fmt.Println()

	for i, comp := range companies {
		marketCap := comp.CurrentPrice * float64(comp.TotalSupply)

		// Insert company using ExecuteInsert
		var createdCompany models.Company
		insertQuery := `INSERT INTO companies (symbol, name, sector, total_supply, current_price, market_cap, shares_outstanding, eps, pe_ratio, book_value, pbv, week_52_high, week_52_low, avg_120_day, yield_1_year, listed_date, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING *`
		err := client.ExecuteInsert(insertQuery, &createdCompany,
			comp.Symbol, comp.Name, comp.Sector, comp.TotalSupply,
			comp.CurrentPrice, marketCap,
			comp.SharesOutstanding, comp.EPS, comp.PERatio,
			comp.BookValue, comp.PBV, comp.Week52High, comp.Week52Low,
			comp.Avg120Day, comp.Yield1Year, comp.ListedDate, comp.Description,
		)
		if err != nil {
			log.Printf("Failed to create company %s: %v", comp.Symbol, err)
			continue
		}

		companyID := createdCompany.ID
		if companyID == "" {
			log.Printf("Failed to get ID for company %s", comp.Symbol)
			continue
		}

		fmt.Printf("[%d/15] Created company: %s (%s) ID: %s\n", i+1, comp.Symbol, comp.Name, companyID)

		// Seed 10 events per company
		baseDate := time.Now()
		for j, eventType := range eventTypes {
			daysOffset := rand.Intn(365) - 180
			eventDate := baseDate.AddDate(0, 0, daysOffset)
			status := models.EventStatusUpcoming
			if eventDate.Before(time.Now()) {
				if rand.Intn(2) == 0 {
					status = models.EventStatusCompleted
				} else {
					status = models.EventStatusCancelled
				}
			}
			fiscalYear := fmt.Sprintf("FY %d/%02d", eventDate.Year()-1, eventDate.Year()%100)
			title := fmt.Sprintf("%s - %s %s", comp.Symbol, eventTitles[eventType], fiscalYear)

			var createdEvent models.CompanyEvent
			eventQuery := `INSERT INTO company_events (company_id, event_type, title, description, event_date, fiscal_year, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
			if err := client.ExecuteInsert(eventQuery, &createdEvent,
				companyID, eventType, title, eventDescs[eventType],
				eventDate.Format("2006-01-02T15:04:05Z"),
				fiscalYear, status,
			); err != nil {
				log.Printf("  Failed to create event %d for %s: %v", j+1, comp.Symbol, err)
			}
		}
		fmt.Printf("  Created 10 events for %s\n", comp.Symbol)

		// Seed 90 days of OHLCV candlestick data
		price := comp.CurrentPrice
		for day := 89; day >= 0; day-- {
			date := time.Now().AddDate(0, 0, -day)
			dayStr := date.Format("2006-01-02T00:00:00Z")

			changePct := (rand.Float64() - 0.48) * 4.0
			openPrice := price
			closePrice := openPrice * (1 + changePct/100.0)

			if closePrice < comp.Week52Low*0.9 {
				closePrice = comp.Week52Low * (0.9 + rand.Float64()*0.1)
			}
			if closePrice > comp.Week52High*1.1 {
				closePrice = comp.Week52High * (1.0 + rand.Float64()*0.1)
			}

			high := openPrice * (1 + rand.Float64()*2.0/100.0)
			low := openPrice * (1 - rand.Float64()*2.0/100.0)
			if closePrice > high {
				high = closePrice * (1 + rand.Float64()*0.5/100.0)
			}
			if closePrice < low {
				low = closePrice * (1 - rand.Float64()*0.5/100.0)
			}

			volume := int64(rand.Intn(50000) + 1000)
			turnover := closePrice * float64(volume)

			priceDiff := closePrice - openPrice
			pctChange := 0.0
			if openPrice > 0 {
				pctChange = (priceDiff / openPrice) * 100.0
			}

			var createdPrice models.StockPrice
			priceQuery := `INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, turnover, change_percent, timestamp, timeframe) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`
			if err := client.ExecuteInsert(priceQuery, &createdPrice,
				companyID, openPrice, high, low, closePrice, volume, turnover, pctChange, dayStr, "1D",
			); err != nil {
				log.Printf("  Failed to create price for %s day %d: %v", comp.Symbol, day, err)
			}

			price = closePrice
		}
		fmt.Printf("  Created 90 days of OHLCV data for %s\n", comp.Symbol)

		// Update company current_price to last seeded price
		finalPrice := decimal.NewFromFloat(price)
		finalMarketCap := finalPrice.Mul(decimal.NewFromInt(comp.TotalSupply))
		updateQuery := fmt.Sprintf(
			"UPDATE companies SET current_price = %s, market_cap = %s WHERE id = '%s'",
			finalPrice.StringFixed(2), finalMarketCap.StringFixed(2), companyID,
		)
		if err := client.ExecuteUpdate(updateQuery, nil); err != nil {
			log.Printf("  Failed to update final price for %s: %v", comp.Symbol, err)
		}
	}

	fmt.Println()
	fmt.Println("=== Seeding Complete ===")
	fmt.Printf("Companies: 15\n")
	fmt.Printf("Events: 150 (10 per company)\n")
	fmt.Printf("Candlestick records: 1350 (90 days x 15 companies)\n")

	os.Exit(0)
}
