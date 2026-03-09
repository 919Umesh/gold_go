# Stock Market Simulator API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All authenticated endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## Health Check

### GET /health
Returns service status.

**Response:**
```json
{ "status": "ok", "service": "stock-market-simulator" }
```

---

## Auth Endpoints

### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "+9779800000000",
  "password": "password123"
}
```

### POST /auth/login
Login and receive a JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { "id": "uuid", "full_name": "John Doe", "email": "john@example.com", "profile_image_url": "https://..." }
}
```

### GET /auth/profile 🔒
Get authenticated user's profile (includes `profile_image_url`).

### PUT /auth/profile 🔒
Update profile information.

### PUT /auth/kyc 🔒
Update KYC verification status.

### POST /auth/profile/image 🔒
Upload profile image. Deletes old image from Supabase storage before uploading new one.

**Request:** `multipart/form-data` with `image` field.

---

## Market Data Endpoints (Public)

### GET /market/companies
List all companies.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 50 | Max results |
| offset | int | 0 | Pagination offset |

### GET /market/companies/new
Get recently listed companies.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/companies/old
Get oldest listed companies.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/companies/:id
Get company detail by ID.

### GET /market/live
Get live trading data for all companies.

**Response:**
```json
{
  "data": [
    {
      "symbol": "NABIL",
      "company_id": "uuid",
      "company_name": "Nabil Bank Limited",
      "sector": "Commercial Banks",
      "ltp": 1250.00,
      "change_percent": 2.35,
      "open": 1220.00,
      "high": 1260.00,
      "low": 1215.00,
      "volume": 15000,
      "previous_close": 1221.00,
      "difference": 29.00,
      "turnover": 18750000.00,
      "last_updated": "2025-03-09T12:00:00Z"
    }
  ],
  "count": 15
}
```

### GET /market/index
Get market index summary (NEPSE-like index).

**Response:**
```json
{
  "data": {
    "index_value": 2145.67,
    "change": 12.34,
    "change_percent": 0.58,
    "total_turnover": 125000000.00,
    "total_volume": 250000,
    "total_market_cap": 2145670000000.00,
    "advances": 8,
    "declines": 5,
    "unchanged": 2,
    "total_companies": 15,
    "previous_close": 2133.33,
    "timestamp": "2025-03-09T12:00:00Z"
  }
}
```

### GET /market/candlestick
Get OHLCV candlestick data for chart rendering.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| symbol | string | required | Stock symbol (e.g., NABIL) |
| timeframe | string | 1D | 1m, 5m, 15m, 1h, 1D |
| days | int | 90 | Number of days of history |

**Response:**
```json
{
  "data": [
    {
      "timestamp": "2025-01-01T00:00:00Z",
      "open": 1220.00,
      "high": 1260.00,
      "low": 1200.00,
      "close": 1250.00,
      "volume": 15000,
      "turnover": 18750000.00,
      "change_percent": 2.46
    }
  ],
  "symbol": "NABIL",
  "timeframe": "1D",
  "count": 90
}
```

### GET /market/top-gainers
Get top gaining stocks by percentage change.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/top-losers
Get top losing stocks by percentage change.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/most-active
Get most actively traded stocks by volume.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/top-turnover
Get stocks with highest turnover.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Max results |

### GET /market/sectors
Get sector performance summary.

**Response:**
```json
{
  "data": [
    {
      "sector": "Commercial Banks",
      "avg_change": 1.85,
      "total_turnover": 50000000.00,
      "total_volume": 80000,
      "total_market_cap": 1500000000.00,
      "company_count": 5
    }
  ],
  "count": 7
}
```

### GET /market/sectors/:sector/companies
Get all companies in a specific sector.

### GET /market/stream
Server-Sent Events (SSE) stream for real-time price updates.

**Headers:** `Content-Type: text/event-stream`

---

## Price Triggers 🔒

### POST /market/triggers
Create a price trigger for auto-sell.

**Request Body:**
```json
{
  "company_id": "uuid",
  "trigger_price": 1300.00,
  "shares_qty": 100,
  "direction": "above"
}
```

Direction: `above` (sell when price goes above), `below` (sell when price goes below).

### GET /market/triggers
Get authenticated user's price triggers.

### PUT /market/triggers/:id/cancel
Cancel a price trigger.

---

## Company Events (Public)

### GET /events/
Get all company events.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 50 | Max results |
| offset | int | 0 | Pagination offset |

### GET /events/company/:company_id
Get events for a specific company.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 20 | Max results |

### GET /events/upcoming
Get upcoming events across all companies.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 20 | Max results |

### GET /events/type/:event_type
Get events by type.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 20 | Max results |

**Event Types:** `agm`, `dividend`, `bonus_share`, `right_share`, `quarterly_report`, `board_meeting`, `financial_results`, `stock_split`, `merger_acquisition`, `ipo_announcement`

**Event Statuses:** `upcoming`, `completed`, `cancelled`, `postponed`

---

## Wallet Endpoints 🔒

### GET /wallet/
Get all user wallets (main + trading).

### GET /wallet/main
Get main wallet balance.

### GET /wallet/trading
Get trading wallet balance.

### POST /wallet/topup
Top up main wallet.

**Request Body:**
```json
{ "amount": 100000 }
```

### POST /wallet/transfer
Transfer between main and trading wallets.

**Request Body:**
```json
{ "amount": 50000, "direction": "main_to_trading" }
```

### GET /wallet/transfer-history
Get wallet transfer history.

---

## Order / Trading Endpoints 🔒

### POST /orders/buy
Place a buy order.

**Request Body:**
```json
{
  "company_id": "uuid",
  "price": 1250.00,
  "quantity": 100,
  "order_type": "limit"
}
```

### POST /orders/sell
Place a sell order.

### PUT /orders/:id/cancel
Cancel a pending order.

### GET /orders/book/:company_id
Get order book for a company.

### GET /orders/my-orders
Get authenticated user's orders.

### GET /orders/portfolio
Get authenticated user's portfolio.

### GET /orders/trades
Get authenticated user's trade history.

---

## IPO Endpoints

### GET /ipo/
List all IPOs.

### GET /ipo/:id
Get IPO detail.

### POST /ipo/apply 🔒
Apply for an IPO.

### POST /ipo/create 🔒 (Admin)
Create a new company via IPO.

### POST /ipo/launch 🔒 (Admin)
Launch an IPO for a company.

### POST /ipo/allocate 🔒 (Admin)
Allocate IPO shares to applicants.

---

## Interactive Swagger UI
Available at: `http://localhost:8080/swagger/index.html`

---

## Seed Data
Run the seed script to populate 15 Nepali companies with events and candlestick data:
```bash
go run scripts/main.go
```

This creates:
- **15 companies** across 7 sectors (Commercial Banks, Telecom, Manufacturing, Life Insurance, Non-Life Insurance, Hydropower, Infrastructure Development Bank)
- **150 company events** (10 per company: AGM, dividend, bonus share, right share, quarterly report, board meeting, financial results, stock split, merger/acquisition, IPO announcement)
- **1,350 candlestick records** (90 days of daily OHLCV data per company)
