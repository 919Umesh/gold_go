package api

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type AdminHandler struct {
	stockRepo stock.Repository
	marketSim *market.Simulator
}

func NewAdminHandler(stockRepo stock.Repository, marketSim *market.Simulator) *AdminHandler {
	return &AdminHandler{
		stockRepo: stockRepo,
		marketSim: marketSim,
	}
}

type NepaliCompany struct {
	Symbol      string
	Name        string
	Sector      string
	MarketCap   float64
	Description string
	FoundedYear int
	Employees   int
	TotalShares int64
}

var nepaliCompanies = []NepaliCompany{
	{Symbol: "NABIL", Name: "Nabil Bank Limited", Sector: "Banking", MarketCap: 180000000000, Description: "Leading private sector bank in Nepal", FoundedYear: 1984, Employees: 3500, TotalShares: 1800000000},
	{Symbol: "GBIME", Name: "Global IME Bank Limited", Sector: "Banking", MarketCap: 160000000000, Description: "Largest commercial bank by network", FoundedYear: 2007, Employees: 4200, TotalShares: 1600000000},
	{Symbol: "NICA", Name: "Nepal Investment Mega Bank", Sector: "Banking", MarketCap: 150000000000, Description: "Merged large commercial bank", FoundedYear: 1986, Employees: 4000, TotalShares: 1500000000},
	{Symbol: "NMB", Name: "NMB Bank Limited", Sector: "Banking", MarketCap: 140000000000, Description: "Corporate-focused commercial bank", FoundedYear: 1996, Employees: 3200, TotalShares: 1400000000},
	{Symbol: "HBL", Name: "Himalayan Bank Limited", Sector: "Banking", MarketCap: 130000000000, Description: "Joint venture bank with strong presence", FoundedYear: 1993, Employees: 2800, TotalShares: 1300000000},

	{Symbol: "NTC", Name: "Nepal Telecom", Sector: "Information Technology", MarketCap: 200000000000, Description: "Largest telecom & IT service provider", FoundedYear: 2004, Employees: 5000, TotalShares: 2000000000},
	{Symbol: "NITC", Name: "Nepal IT Corporation", Sector: "Information Technology", MarketCap: 90000000000, Description: "Government-backed IT service company", FoundedYear: 1998, Employees: 1500, TotalShares: 900000000},
	{Symbol: "F1SOFT", Name: "F1Soft International", Sector: "Information Technology", MarketCap: 70000000000, Description: "Digital payment and fintech company", FoundedYear: 2004, Employees: 1200, TotalShares: 700000000},
	{Symbol: "ESewa", Name: "eSewa Digital Services", Sector: "Information Technology", MarketCap: 60000000000, Description: "Fintech wallet operator", FoundedYear: 2009, Employees: 900, TotalShares: 600000000},
	{Symbol: "IMS", Name: "IMS Software Solutions", Sector: "Information Technology", MarketCap: 30000000000, Description: "IT consulting and development", FoundedYear: 2005, Employees: 600, TotalShares: 300000000},

	{Symbol: "HYDRO1", Name: "Upper Tamakoshi Hydropower", Sector: "Hydropower", MarketCap: 120000000000, Description: "Large hydropower producer", FoundedYear: 2007, Employees: 900, TotalShares: 1200000000},
	{Symbol: "BPC", Name: "Butwal Power Company", Sector: "Hydropower", MarketCap: 55000000000, Description: "Integrated power producer", FoundedYear: 1966, Employees: 600, TotalShares: 550000000},
	{Symbol: "CHCL", Name: "Chilime Hydropower", Sector: "Hydropower", MarketCap: 60000000000, Description: "Hydropower generation", FoundedYear: 1995, Employees: 450, TotalShares: 600000000},
	{Symbol: "API", Name: "Api Power Company", Sector: "Hydropower", MarketCap: 40000000000, Description: "Hydropower & energy developer", FoundedYear: 2003, Employees: 350, TotalShares: 400000000},
	{Symbol: "RSHP", Name: "Rosuwa Shyamkhola Hydro Power", Sector: "Hydropower", MarketCap: 35000000000, Description: "Small-medium hydropower", FoundedYear: 2010, Employees: 280, TotalShares: 350000000},

	{Symbol: "NLIC", Name: "Nepal Life Insurance Company", Sector: "Insurance", MarketCap: 130000000000, Description: "Largest life insurance provider", FoundedYear: 2001, Employees: 1200, TotalShares: 1300000000},
	{Symbol: "SICL", Name: "Shikhar Insurance Company", Sector: "Insurance", MarketCap: 45000000000, Description: "Leading non-life insurance provider", FoundedYear: 2004, Employees: 800, TotalShares: 450000000},
	{Symbol: "NMBHL", Name: "NMB Health Insurance", Sector: "Insurance", MarketCap: 35000000000, Description: "Health insurance specialist", FoundedYear: 2010, Employees: 700, TotalShares: 350000000},

	{Symbol: "APOLLONP", Name: "Apollo Nepal Hospitals", Sector: "Pharma", MarketCap: 50000000000, Description: "Hospital & healthcare chain", FoundedYear: 2015, Employees: 1500, TotalShares: 500000000},
	{Symbol: "MHPL", Name: "Medical Health Products Limited", Sector: "Pharma", MarketCap: 40000000000, Description: "Pharmaceutical manufacturer", FoundedYear: 2008, Employees: 900, TotalShares: 400000000},

	{Symbol: "HDL", Name: "Himalayan Distillery Limited", Sector: "Manufacturing", MarketCap: 90000000000, Description: "Distillery & beverages producer", FoundedYear: 1985, Employees: 700, TotalShares: 900000000},
	{Symbol: "UNL", Name: "Unilever Nepal Limited", Sector: "Manufacturing", MarketCap: 80000000000, Description: "FMCG with global brands", FoundedYear: 1992, Employees: 300, TotalShares: 800000000},
	{Symbol: "BNL", Name: "Bottlers Nepal Limited", Sector: "Manufacturing", MarketCap: 60000000000, Description: "Bottled beverages producer", FoundedYear: 1979, Employees: 450, TotalShares: 600000000},
	{Symbol: "SHIVM", Name: "Shivam Cements Limited", Sector: "Manufacturing", MarketCap: 50000000000, Description: "Cement manufacturer", FoundedYear: 2003, Employees: 1100, TotalShares: 500000000},

	{Symbol: "DLFNP", Name: "Nepal Housing Development Co.", Sector: "Real Estate", MarketCap: 40000000000, Description: "Real estate developer", FoundedYear: 2008, Employees: 500, TotalShares: 400000000},
}

var adminEventTemplates = []struct {
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
	{models.MarketEventDividend, "Interim Dividend", "Interim dividend of ₹50 per share", 1.8},
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
	{models.MarketEventDividend, "Special Dividend", "Special one-time dividend of ₹100 per share", 5.0},
	{models.MarketEventNews, "Export Agreement", "Major export agreement signed with foreign entity", 3.5},
	{models.MarketEventIPO, "FPO Announcement", "Further public offering announced", 2.8},
	{models.MarketEventNews, "Infrastructure Investment", "Major infrastructure investment announced", 2.0},
	{models.MarketEventNews, "Leadership Change", "New CEO appointed with strong track record", 1.8},
	{models.MarketEventEarnings, "Q4 Financial Results", "Year-end earnings show record profits", 4.2},
}

func (h *AdminHandler) SeedStockData(c *gin.Context) {
	companies, err := h.stockRepo.ListCompanies(1, 0)
	if err == nil && len(companies) > 0 {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Database already seeded. Delete existing data via Supabase dashboard if re-seeding needed.")
		return
	}

	slog.Info("Starting database seeding via API...")

	companiesCreated := 0
	pricesCreated := 0
	eventsCreated := 0

	companyMap := make(map[string]string)

	for _, nc := range nepaliCompanies {
		existing, err := h.stockRepo.GetCompanyBySymbol(nc.Symbol)
		if err == nil && existing != nil {
			companyMap[nc.Symbol] = existing.ID
			continue
		}

		company := &models.Company{
			Symbol:          nc.Symbol,
			Name:            nc.Name,
			Sector:          nc.Sector,
			MarketCap:       decimal.NewFromFloat(nc.MarketCap),
			Description:     nc.Description,
			FoundedYear:     nc.FoundedYear,
			Employees:       nc.Employees,
			TotalShares:     nc.TotalShares,
			AvailableShares: nc.TotalShares,
			IsActive:        true,
		}

		if err := h.stockRepo.CreateCompany(company); err != nil {
			slog.Error("Failed to create company", "symbol", nc.Symbol, "error", err)
			continue
		}

		companyMap[nc.Symbol] = company.ID
		companiesCreated++
	}

	initialPrice := decimal.NewFromInt(100)
	for _, nc := range nepaliCompanies {
		companyID, exists := companyMap[nc.Symbol]
		if !exists {
			continue
		}

		stockPrice := &models.StockPrice{
			CompanyID:  companyID,
			OpenPrice:  initialPrice,
			HighPrice:  initialPrice,
			LowPrice:   initialPrice,
			ClosePrice: initialPrice,
			Volume:     0,
			Timestamp:  time.Now(),
			Timeframe:  "1D",
		}

		if err := h.stockRepo.CreateStockPrice(stockPrice); err != nil {
			slog.Error("Failed to create initial price", "symbol", nc.Symbol, "error", err)
			continue
		}
		pricesCreated++
	}

	for _, nc := range nepaliCompanies {
		companyID, exists := companyMap[nc.Symbol]
		if !exists {
			continue
		}

		for i := 0; i < 25; i++ {
			evTemplate := adminEventTemplates[i%len(adminEventTemplates)]
			impact := evTemplate.impact + (rand.Float64()-0.5)*1.0
			eventDate := time.Now().Add(time.Duration(rand.Intn(365)+1) * 24 * time.Hour)

			event := &models.MarketEvent{
				CompanyID:        companyID,
				EventType:        evTemplate.eventType,
				Title:            evTemplate.title,
				Description:      fmt.Sprintf("%s - %s", evTemplate.description, nc.Name),
				ImpactPercentage: impact,
				EventDate:        eventDate,
			}

			if err := h.stockRepo.CreateMarketEvent(event); err == nil {
				eventsCreated++
			}
		}
	}

	slog.Info("Database seeding completed", "companies", companiesCreated, "prices", pricesCreated, "events", eventsCreated)

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"message":           "Database seeded successfully",
		"companies_created": companiesCreated,
		"prices_created":    pricesCreated,
		"events_created":    eventsCreated,
		"notes":             "All companies start at ₹100. No test users or transactions created. Users register via API, wallet starts at ₹0.",
	})
}

func (h *AdminHandler) TriggerMarketUpdate(c *gin.Context) {
	slog.Info("Admin triggered a manual market price update")
	h.marketSim.UpdateMarket()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Market prices updated successfully.",
	})
}
