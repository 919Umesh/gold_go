#!/bin/bash
# Stock Market Simulator - Complete Setup and Run Guide

echo "==============================================="
echo "Stock Market Simulator - Setup Script"
echo "==============================================="
echo ""

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

cd /Users/inovaaratechnologies/gold_go

echo -e "${BLUE}Step 1: Setting up Appwrite database...${NC}"
echo "Running: go run scripts/setup_appwrite.go"
echo ""

# Run the setup script
go run scripts/setup_appwrite.go

echo ""
echo -e "${GREEN}✅ Database setup complete!${NC}"
echo ""
echo "This setup has created:"
echo "  • 25 Nepali companies across 5 sectors"
echo "  • 750 stock prices (30 days per company)"
echo "  • 250 market events (10 per company)"
echo "  • 750+ transactions for market data"
echo "  • 3 test users with ₹10,00,000 balance each"
echo ""

echo -e "${BLUE}Step 2: Starting API server...${NC}"
echo "Running: go run cmd/main.go"
echo ""
echo "Server will be available at: http://localhost:8080"
echo ""

# Start the server
go run cmd/main.go
