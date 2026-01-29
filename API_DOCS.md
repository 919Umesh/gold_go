# Gold Go - Complete API Documentation

**Base URL**: `http://localhost:8080/api/v1` (local) or `https://your-app.onrender.com/api/v1` (production)

**Version**: 1.0  
**Last Updated**: January 29, 2026

---

## Table of Contents

1. [Authentication](#authentication)
2. [Gold Trading](#gold-trading)
3. [Wallet Management](#wallet-management)
4. [Stock Market](#stock-market)
5. [Stock Trading](#stock-trading)
6. [Admin Operations](#admin-operations)
7. [Error Codes](#error-codes)

---

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Register User

**Endpoint**: `POST /auth/register`  
**Authentication**: None  
**Rate Limit**: 60 requests per hour

**Request Body**:
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "9841234567",
  "password": "SecurePass123",
  "role": "user"
}
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9841234567",
    "role": "user",
    "kyc_status": "pending"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Login

**Endpoint**: `POST /auth/login`  
**Authentication**: None  
**Rate Limit**: 60 requests per 5 minutes

**Request Body**:
```json
{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Get Profile

**Endpoint**: `GET /auth/profile`  
**Authentication**: Required  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9841234567",
    "role": "user",
    "kyc_status": "pending",
    "created_at": "2026-01-29T10:00:00Z"
  }
}
```

---

### Update Profile

**Endpoint**: `PUT /auth/profile/update`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "full_name": "John Updated",
  "phone": "9841234568"
}
```

**Response** (200 OK):
```json
{
  "message": "Profile updated successfully",
  "user": {...}
}
```

---

## Gold Trading

### Get Current Gold Price

**Endpoint**: `GET /gold/price`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "price_per_gram": 6523.45,
  "updated_at": "2026-01-29T12:00:00Z",
  "source": "provider"
}
```

---

### Get Gold Price History

**Endpoint**: `GET /gold/history?days=30`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: 600 requests per minute

**Query Parameters**:
- `days` (optional): Number of days (default: 7, max: 365)

**Response** (200 OK):
```json
{
  "prices": [
    {
      "id": 1,
      "price_per_gram": 6500.00,
      "source": "provider",
      "updated_at": "2026-01-28T12:00:00Z"
    }
  ]
}
```

---

## Wallet Management

All wallet endpoints require authentication.

### Get Wallet Balance

**Endpoint**: `GET /wallet`  
**Authentication**: Required  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "wallet": {
    "id": 1,
    "user_id": 1,
    "fiat_balance": 50000.00,
    "gold_grams": 10.5,
    "locked": false
  }
}
```

---

### Top Up Wallet

**Endpoint**: `POST /wallet/topup`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "amount": 10000.00,
  "reference_id": "payment-ref-12345"
}
```

**Response** (200 OK):
```json
{
  "message": "Wallet topped up successfully",
  "wallet": {
    "fiat_balance": 60000.00,
    "gold_grams": 10.5
  },
  "transaction": {
    "id": 123,
    "type": "topup",
    "amount": 10000.00,
    "status": "success"
  }
}
```

---

### Buy Gold

**Endpoint**: `POST /wallet/buy`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "grams": 5.0,
  "price_per_gram": 6500.00,
  "reference_id": "buy-ref-456"
}
```

**Response** (200 OK):
```json
{
  "message": "Gold purchased successfully",
  "wallet": {
    "fiat_balance": 27500.00,
    "gold_grams": 15.5
  },
  "transaction": {
    "id": 124,
    "type": "buy",
    "amount": 32500.00,
    "gold_grams": 5.0,
    "price_per_gram": 6500.00
  }
}
```

---

### Sell Gold

**Endpoint**: `POST /wallet/sell`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "grams": 2.0,
  "price_per_gram": 6550.00,
  "reference_id": "sell-ref-789"
}
```

**Response** (200 OK):
```json
{
  "message": "Gold sold successfully",
  "wallet": {
    "fiat_balance": 40600.00,
    "gold_grams": 13.5
  },
  "transaction": {
    "id": 125,
    "type": "sell",
    "amount": 13100.00,
    "gold_grams": 2.0,
    "price_per_gram": 6550.00
  }
}
```

---

### Get Transaction History

**Endpoint**: `GET /transaction`  
**Authentication**: Required  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "transactions": [
    {
      "id": 125,
      "user_id": 1,
      "type": "sell",
      "amount": 13100.00,
      "gold_grams": 2.0,
      "price_per_gram": 6550.00,
      "status": "success",
      "created_at": "2026-01-29T12:30:00Z"
    }
  ]
}
```

---

## Stock Market

All stock market data endpoints are public (no authentication required).

### List All Companies

**Endpoint**: `GET /stocks?limit=50&offset=0`  
**Authentication**: None  
**Cache**: 2 minutes  
**Rate Limit**: Standard

**Query Parameters**:
- `limit` (optional): Results per page (default: 50, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response** (200 OK):
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

---

### Search Companies

**Endpoint**: `GET /stocks/search?q=nabil`  
**Authentication**: None  
**Rate Limit**: Standard

**Query Parameters**:
- `q` (required): Search query (searches name, symbol, sector)

**Response** (200 OK):
```json
{
  "companies": [...]
}
```

---

### Get Companies by Sector

**Endpoint**: `GET /stocks/sector/Banking?limit=50&offset=0`  
**Authentication**: None  
**Cache**: 2 minutes  
**Rate Limit**: Standard

**Available Sectors**: Banking, Hydropower, Insurance, Manufacturing, Hotels, Finance

**Response** (200 OK):
```json
{
  "sector": "Banking",
  "companies": [...],
  "limit": 50,
  "offset": 0
}
```

---

### Get Company Details

**Endpoint**: `GET /stocks/:symbol`  
**Authentication**: None  
**Cache**: 2 minutes  
**Rate Limit**: Standard

**Example**: `GET /stocks/NABIL`

**Response** (200 OK):
```json
{
  "company": {
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
}
```

---

### Get Current Stock Price

**Endpoint**: `GET /stocks/:symbol/price`  
**Authentication**: None  
**Cache**: 30 seconds  
**Rate Limit**: Standard

**Example**: `GET /stocks/NABIL/price`

**Response** (200 OK):
```json
{
  "price": {
    "id": 12345,
    "company_id": 1,
    "open_price": 1200.50,
    "high_price": 1215.00,
    "low_price": 1195.00,
    "close_price": 1210.00,
    "volume": 125000,
    "timestamp": "2026-01-29T12:00:00Z",
    "timeframe": "1d"
  }
}
```

---

### Get Price History

**Endpoint**: `GET /stocks/:symbol/history?timeframe=1d&days=30`  
**Authentication**: None  
**Cache**: 5 minutes  
**Rate Limit**: Standard

**Query Parameters**:
- `timeframe` (optional): 1m, 5m, 15m, 1h, 1d, 1w, 1M (default: 1d)
- `days` (optional): Number of days (default: 30)

**Example**: `GET /stocks/NABIL/history?timeframe=1d&days=7`

**Response** (200 OK):
```json
{
  "symbol": "NABIL",
  "timeframe": "1d",
  "prices": [
    {
      "open_price": 1200.00,
      "high_price": 1220.00,
      "low_price": 1190.00,
      "close_price": 1210.00,
      "volume": 150000,
      "timestamp": "2026-01-29T00:00:00Z"
    }
  ]
}
```

---

### Get Market Overview

**Endpoint**: `GET /stocks/market-overview`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "total_companies": 21,
  "top_gainers": [
    {
      "id": 5,
      "symbol": "NICA",
      "name": "NIC Asia Bank",
      "current_price": 1050.00,
      "previous_price": 1000.00,
      "change": 50.00,
      "change_percent": 5.00
    }
  ],
  "top_losers": [...],
  "most_active": [...]
}
```

---

### Get Top Gainers

**Endpoint**: `GET /stocks/top-gainers?limit=10`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: Standard

**Query Parameters**:
- `limit` (optional): Number of results (default: 10)

**Response** (200 OK):
```json
{
  "gainers": [
    {
      "symbol": "NICA",
      "name": "NIC Asia Bank",
      "current_price": 1050.00,
      "previous_price": 1000.00,
      "change": 50.00,
      "change_percent": 5.00
    }
  ]
}
```

---

### Get Top Losers

**Endpoint**: `GET /stocks/top-losers?limit=10`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: Standard

**Response**: Same format as Top Gainers

---

### Get Most Active Stocks

**Endpoint**: `GET /stocks/most-active?limit=10`  
**Authentication**: None  
**Cache**: 1 minute  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "active": [
    {
      "id": 1,
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking"
    }
  ]
}
```

---

### Get Upcoming Events

**Endpoint**: `GET /stocks/:symbol/events`  
**Authentication**: None  
**Cache**: 5 minutes  
**Rate Limit**: Standard

**Example**: `GET /stocks/NABIL/events`

**Response** (200 OK):
```json
{
  "events": [
    {
      "id": 1,
      "company_id": 1,
      "event_type": "earnings",
      "title": "Q1 Earnings Report Released",
      "description": "Quarterly financial results",
      "impact_percentage": 2.5,
      "event_date": "2026-02-15T00:00:00Z"
    }
  ]
}
```

---

## Stock Trading

All trading endpoints require authentication.

### Get Virtual Wallet

**Endpoint**: `GET /trading/wallet`  
**Authentication**: Required  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "wallet": {
    "id": 1,
    "user_id": 1,
    "balance": 950000.00,
    "total_invested": 50000.00,
    "total_profit_loss": 2500.00,
    "created_at": "2026-01-29T10:00:00Z",
    "updated_at": "2026-01-29T12:00:00Z"
  }
}
```

**Note**: New users automatically receive NPR 1,000,000 (10 lakh) virtual balance.

---

### Get Portfolio

**Endpoint**: `GET /trading/portfolio`  
**Authentication**: Required  
**Rate Limit**: Standard

**Response** (200 OK):
```json
{
  "total_value": 52500.00,
  "total_invested": 50000.00,
  "total_profit_loss": 2500.00,
  "profit_loss_percent": 5.00,
  "holdings": [
    {
      "id": 1,
      "user_id": 1,
      "company_id": 1,
      "company_symbol": "NABIL",
      "company_name": "Nabil Bank Limited",
      "quantity": 50,
      "avg_buy_price": 1000.00,
      "total_invested": 50000.00,
      "current_price": 1050.00,
      "current_value": 52500.00,
      "profit_loss": 2500.00,
      "profit_loss_percent": 5.00
    }
  ]
}
```

---

### Buy Stock

**Endpoint**: `POST /trading/buy`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "symbol": "NABIL",
  "quantity": 10
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Stock purchased successfully",
  "quantity": 10,
  "price_per_share": 1210.00,
  "total_amount": 12100.00,
  "new_balance": 937900.00
}
```

**Error Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "Insufficient balance. Required: NPR 12100.00, Available: NPR 10000.00"
}
```

---

### Sell Stock

**Endpoint**: `POST /trading/sell`  
**Authentication**: Required  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "symbol": "NABIL",
  "quantity": 5
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Stock sold successfully",
  "quantity": 5,
  "price_per_share": 1210.00,
  "total_amount": 6050.00,
  "new_balance": 943950.00
}
```

**Error Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "Insufficient shares. You own 3 shares"
}
```

---

### Get Transaction History

**Endpoint**: `GET /trading/transactions?limit=50&offset=0`  
**Authentication**: Required  
**Rate Limit**: Standard

**Query Parameters**:
- `limit` (optional): Results per page (default: 50, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response** (200 OK):
```json
{
  "transactions": [
    {
      "id": 1,
      "user_id": 1,
      "company_id": 1,
      "type": "buy",
      "quantity": 10,
      "price_per_share": 1210.00,
      "total_amount": 12100.00,
      "status": "completed",
      "created_at": "2026-01-29T12:00:00Z"
    }
  ],
  "limit": 50,
  "offset": 0
}
```

---

## Admin Operations

All admin endpoints require authentication with admin role.

### Update User KYC Status

**Endpoint**: `PUT /admin/users/:user_id/kyc`  
**Authentication**: Required (Admin only)  
**Rate Limit**: Standard

**Request Body**:
```json
{
  "kyc_status": "approved"
}
```

**Valid KYC Statuses**: pending, approved, rejected

**Response** (200 OK):
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": 1,
    "kyc_status": "approved"
  }
}
```

---

### Seed Stock Market Data

**Endpoint**: `POST /admin/seed-stocks`  
**Authentication**: Required (Admin only)  
**Rate Limit**: None

**Request Body**: None

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Database seeded successfully",
  "companies_created": 21,
  "prices_created": 5355,
  "events_created": 20
}
```

**Error Response** (400 Bad Request):
```json
{
  "message": "Database already seeded. Delete existing data first if you want to re-seed."
}
```

**Note**: This endpoint can only be called once. It populates the database with 21 Nepali companies and 1 year of historical price data.

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters or validation error |
| 401 | Unauthorized - Missing or invalid authentication token |
| 403 | Forbidden - Insufficient permissions (e.g., not admin) |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists (e.g., duplicate email) |
| 500 | Internal Server Error |

---

## Market Simulation

### How It Works

1. **Market Hours**: The market simulator runs during Nepal Stock Exchange (NEPSE) hours:
   - **Days**: Sunday - Thursday
   - **Time**: 11:00 AM - 3:00 PM Nepal Time
   - **Closed**: Friday & Saturday

2. **Price Updates**: Stock prices update every 5 minutes during market hours

3. **Algorithm**: Uses Geometric Brownian Motion for realistic price movements

4. **Initial Balance**: New users receive NPR 1,000,000 (10 lakh) virtual currency for stock trading

---

## Seeded Companies

The system includes 21 demo Nepali companies across 6 sectors:

**Banking (5)**:
- NABIL - Nabil Bank Limited
- SCB - Standard Chartered Bank Nepal
- HBL - Himalayan Bank Limited
- EBL - Everest Bank Limited
- NICA - NIC Asia Bank

**Hydropower (4)**:
- CHCL - Chilime Hydropower Company
- NHPC - Nepal Hydro Power Company
- UMHL - Upper Marsyangdi Hydropower
- RADHI - Radhi Bidyut Company

**Insurance (4)**:
- NLIC - Nepal Life Insurance Company
- NLICL - National Life Insurance Company
- SICL - Shikhar Insurance Company
- PRIN - Prime Insurance Company

**Manufacturing (3)**:
- UNL - Unilever Nepal Limited
- NRIC - Nepal Reinsurance Company
- BNT - Bottlers Nepal (Terai)

**Hotels (3)**:
- OHL - Oriental Hotels Limited
- TRHPR - Taragaon Regency Hotel
- SHL - Soaltee Hotel Limited

**Finance (2)**:
- GUFL - Goodwill Finance Limited
- CFCL - Central Finance Company

Each company has 1 year of historical price data for charting and analysis.

---

## Example Workflows

### Complete Trading Workflow

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "9841234567",
    "password": "Test123!",
    "role": "user"
  }'

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }'

# 3. Browse stocks
curl http://localhost:8080/api/v1/stocks/market-overview

# 4. Check virtual wallet (should have NPR 10 lakh)
curl http://localhost:8080/api/v1/trading/wallet \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 5. Buy stocks
curl -X POST http://localhost:8080/api/v1/trading/buy \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NABIL",
    "quantity": 10
  }'

# 6. Check portfolio
curl http://localhost:8080/api/v1/trading/portfolio \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 7. Sell stocks
curl -X POST http://localhost:8080/api/v1/trading/sell \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NABIL",
    "quantity": 5
  }'
```

---

## Rate Limiting

All endpoints are rate-limited to prevent abuse. Specific limits are noted in each endpoint description. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.

---

## Support

For issues or questions:
- Check the deployment guide: `DEPLOYMENT_GUIDE.md`
- Review stock deployment: `STOCK_DEPLOYMENT.md`
- Contact: support@example.com
