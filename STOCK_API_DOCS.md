# Stock Market API Documentation

## Overview
This document describes the API endpoints for the virtual stock market trading system designed for Nepali beginners to learn stock trading.

**Base URL**: `http://localhost:8080/api/v1`

---

## Public Endpoints (No Authentication Required)

### Stock Market Data

#### 1. List All Companies
Get a paginated list of all companies.

```http
GET /stocks?limit=50&offset=0
```

**Query Parameters:**
- `limit` (optional): Number of results (default: 50, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```json
{
  "companies": [
    {
      "id": 1,
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 50000000000,
      "description": "Leading commercial bank in Nepal",
      "founded_year": 1984,
      "employees": 1200,
      "is_active": true
    }
  ],
  "limit": 50,
  "offset": 0
}
```

#### 2. Search Companies
Search for companies by name, symbol, or sector.

```http
GET /stocks/search?q=nabil
```

#### 3. Get Company Details
Get detailed information about a specific company.

```http
GET /stocks/NABIL
```

#### 4. Get Current Stock Price
Get the latest price for a stock.

```http
GET /stocks/NABIL/price
```

#### 5. Get Price History
Get historical price data for charting.

```http
GET /stocks/NABIL/history?timeframe=1d&days=30
```

**Query Parameters:**
- `timeframe` (optional): 1m, 5m, 15m, 1h, 1d, 1w, 1M (default: 1d)
- `days` (optional): Number of days (default: 30)

#### 6. Get Market Overview
Get a summary of market performance.

```http
GET /stocks/market-overview
```

#### 7. Get Top Gainers/Losers/Most Active
```http
GET /stocks/top-gainers?limit=10
GET /stocks/top-losers?limit=10
GET /stocks/most-active?limit=10
```

---

## Protected Endpoints (Authentication Required)

**Authentication:** Include JWT token in Authorization header
```
Authorization: Bearer <your_jwt_token>
```

### Trading Operations

#### 1. Get Virtual Wallet
```http
GET /trading/wallet
```

#### 2. Get Portfolio
```http
GET /trading/portfolio
```

#### 3. Buy Stock
```http
POST /trading/buy
Content-Type: application/json

{
  "symbol": "NABIL",
  "quantity": 10
}
```

#### 4. Sell Stock
```http
POST /trading/sell
Content-Type: application/json

{
  "symbol": "NABIL",
  "quantity": 5
}
```

#### 5. Get Transaction History
```http
GET /trading/transactions?limit=50&offset=0
```

---

## Market Simulation

### How It Works

1. **Market Hours**: Sunday - Thursday, 11:00 AM - 3:00 PM Nepal Time
2. **Price Updates**: Every 5 minutes during market hours
3. **Algorithm**: Geometric Brownian Motion for realistic price movements
4. **Initial Balance**: NPR 1,000,000 (10 lakh) virtual currency

---

## Seeded Companies

The system includes 21 demo Nepali companies across sectors:
- **Banking**: NABIL, SCB, HBL, EBL, NICA
- **Hydropower**: CHCL, NHPC, UMHL, RADHI
- **Insurance**: NLIC, NLICL, SICL, PRIN
- **Manufacturing**: UNL, NRIC, BNT
- **Hotels**: OHL, TRHPR, SHL
- **Finance**: GUFL, CFCL

Each company has 1 year of historical price data for charting and analysis.
