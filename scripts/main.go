package main

import (
	"fmt"
	"log"
	"log/slog"
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
	totalSupply int64
}{
	{"NABIL", "Nabil Bank Limited", "Banking", 10000},
	{"GBIME", "Global IME Bank Limited", "Banking", 10000},
	{"NICA", "Nepal Investment Mega Bank", "Banking", 10000},
	{"NMB", "NMB Bank Limited", "Banking", 10000},
	{"HBL", "Himalayan Bank Limited", "Banking", 10000},

	{"NTC", "Nepal Telecom", "Information Technology", 10000},
	{"F1SOFT", "F1Soft International", "Information Technology", 10000},

	{"HYDRO1", "Upper Tamakoshi Hydropower", "Hydropower", 10000},
	{"BPC", "Butwal Power Company", "Hydropower", 10000},
	{"CHCL", "Chilime Hydropower", "Hydropower", 10000},

	{"NLIC", "Nepal Life Insurance Company", "Insurance", 10000},
	{"SICL", "Shikhar Insurance Company", "Insurance", 10000},

	{"HDL", "Himalayan Distillery Limited", "Manufacturing", 10000},
	{"UNL", "Unilever Nepal Limited", "Manufacturing", 10000},
	{"SHIVM", "Shivam Cements Limited", "Manufacturing", 10000},
}

func setupCompaniesData(client *supabase.Client) error {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Starting Data Setup (Supabase) — New Schema V2")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\nStep 1: Creating/fetching companies...")

	createdCompanies := make([]models.Company, 0)

	for _, comp := range companiesData {
		var existingCompany models.Company
		err := client.ExecuteQueryRow("SELECT * FROM companies WHERE symbol = $1", &existingCompany, comp.symbol)

		if err == nil {
			createdCompanies = append(createdCompanies, existingCompany)
			fmt.Printf("Found existing: %s (%s) - %d total supply\n", existingCompany.Name, existingCompany.Sector, existingCompany.TotalSupply)
			continue
		}

		var company models.Company
		err = client.ExecuteInsert(
			"INSERT INTO companies (symbol, name, sector, total_supply, current_price, market_cap, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *",
			&company,
			comp.symbol, comp.name, comp.sector, comp.totalSupply,
			"0", "0", true)
		if err != nil {
			fmt.Printf("Failed to create company %s: %v\n", comp.symbol, err)
			continue
		}

		createdCompanies = append(createdCompanies, company)
		fmt.Printf("Created: %s (%s) - %d total supply\n", company.Name, company.Sector, comp.totalSupply)
	}

	if len(createdCompanies) == 0 {
		return fmt.Errorf("failed to create companies")
	}

	fmt.Println("\nStep 2: Creating initial stock price record for each company...")
	totalPricesCreated := 0

	for _, company := range createdCompanies {
		initialPrice := decimal.NewFromInt(100)

		err := client.ExecuteInsert(
			"INSERT INTO stock_prices (company_id, open_price, high_price, low_price, close_price, volume, timestamp, timeframe) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *",
			nil,
			company.ID, initialPrice.String(), initialPrice.String(),
			initialPrice.String(), initialPrice.String(),
			0, time.Now().Format(time.RFC3339), "1D")
		if err == nil {
			totalPricesCreated++
			fmt.Printf("Initial price set for %s: 100\n", company.Symbol)
		} else {
			fmt.Printf("Failed to set initial price for %s: %v\n", company.Symbol, err)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("Setup Complete!\n")
	fmt.Printf("   Companies:      %d\n", len(createdCompanies))
	fmt.Printf("   Initial Prices: %d (all 100)\n", totalPricesCreated)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\nNotes:")
	fmt.Println("   - Users register via API, wallets (main + trading) start with 0 balance")
	fmt.Println("   - Users top up main wallet, then transfer to trading wallet")
	fmt.Println("   - Admin launches IPOs, users apply, admin allocates via lottery")
	fmt.Println("   - Only IPO-allocated users can initially sell shares")
	fmt.Println("   - Prices update after every matched trade with high sensitivity")

	return nil
}
