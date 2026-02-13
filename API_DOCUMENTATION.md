# Stock Market Simulator - Complete API Documentation

## 📋 Table of Contents

1. [Base Information](#base-information)
2. [Authentication](#authentication)
3. [Public Endpoints](#public-endpoints)
4. [Protected Endpoints](#protected-endpoints)
5. [Admin Endpoints](#admin-endpoints)
6. [Data Models](#data-models)
7. [Error Handling](#error-handling)
8. [Examples](#examples)

---

## Base Information

**API Version:** v1

**Content-Type:** `application/json`

**Authentication:** JWT Bearer Token (except for public endpoints)

---

## Authentication

### Login

Authenticate a user and receive a JWT token.

**Endpoint:** `POST /auth/login`

**Access:** Public

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Success Response (200 OK):**
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_123",
    "full_name": "John Doe",
    "email": "user@example.com",
    "phone": "+977981234567",
    "kyc_status": "verified",
    "role": "user"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "invalid credentials"
}
```

---

### Register

Create a new user account.

**Endpoint:** `POST /auth/register`

**Access:** Public

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "user@example.com",
  "phone": "+977981234567",
  "password": "securepassword",
  "role": "user"
}
```

**Validation:**
- `full_name`: 2-100 characters
- `email`: Valid email format
- `phone`: 10-15 characters
- `password`: Minimum 6 characters
- `role`: "user" or "admin"

**Success Response (201 Created):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "user_123",
    "full_name": "John Doe",
    "email": "user@example.com",
    "phone": "+977981234567",
    "kyc_status": "pending",
    "role": "user"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "error": "user already exists"
}
```

---

## Public Endpoints

### Stock Market - Browse & Search

#### List All Companies

**Endpoint:** `GET /stocks`

**Access:** Public

**Query Parameters:**
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)
- `sort` (optional): Sort field (name, symbol, price)

**Success Response (200 OK):**
```json
{
  "companies": [
    {
      "id": "company_1",
      "name": "Nepal Telecom",
      "symbol": "NTC",
      "current_price": 1250.50,
      "market_cap": "500B",
      "sector": "Telecommunications"
    }
  ],
  "total": 45,
  "limit": 50,
  "offset": 0
}
```

---

#### Search Companies

**Endpoint:** `GET /stocks/search`

**Access:** Public

**Query Parameters:**
- `q` (required): Search query (name or symbol)

**Success Response (200 OK):**
```json
{
  "companies": [
    {
      "id": "company_1",
      "name": "Nepal Telecom",
      "symbol": "NTC",
      "current_price": 1250.50,
      "market_cap": "500B",
      "sector": "Telecommunications"
    }
  ]
}
```

---

#### Get Company Details

**Endpoint:** `GET /stocks/:symbol`

**Access:** Public

**URL Parameters:**
- `symbol`: Company symbol (e.g., "NTC")

**Success Response (200 OK):**
```json
{
  "company": {
    "id": "company_1",
    "name": "Nepal Telecom",
    "symbol": "NTC",
    "current_price": 1250.50,
    "market_cap": "500B",
    "sector": "Telecommunications",
    "description": "State-owned telecommunications company",
    "pe_ratio": 15.5,
    "dividend_yield": 3.2
  }
}
```

---

#### Get Current Stock Price

**Endpoint:** `GET /stocks/:symbol/price`

**Access:** Public

**URL Parameters:**
- `symbol`: Company symbol (e.g., "NTC")

**Success Response (200 OK):**
```json
{
  "symbol": "NTC",
  "price": 1250.50,
  "currency": "NPR",
  "timestamp": "2026-02-12T10:30:00Z",
  "change_percent": 2.5,
  "change_amount": 30.50
}
```

---

#### Get Price History

**Endpoint:** `GET /stocks/:symbol/history`

**Access:** Public

**Query Parameters:**
- `days` (optional): Number of days to fetch (default: 30)

**Success Response (200 OK):**
```json
{
  "symbol": "NTC",
  "prices": [
    {
      "date": "2026-02-12",
      "open": 1220.00,
      "high": 1260.00,
      "low": 1210.00,
      "close": 1250.50,
      "volume": 1000000
    }
  ]
}
```

---

#### Get Market Overview

**Endpoint:** `GET /stocks/market-overview`

**Access:** Public

**Success Response (200 OK):**
```json
{
  "market_status": "open",
  "total_companies": 45,
  "gainers_count": 18,
  "losers_count": 12,
  "unchanged_count": 15,
  "market_index": 2150.75,
  "index_change_percent": 1.2
}
```

---

#### Get Top Gainers

**Endpoint:** `GET /stocks/top-gainers`

**Access:** Public

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Success Response (200 OK):**
```json
{
  "gainers": [
    {
      "symbol": "NTC",
      "name": "Nepal Telecom",
      "price": 1250.50,
      "change_percent": 5.2,
      "change_amount": 61.75
    }
  ]
}
```

---

#### Get Top Losers

**Endpoint:** `GET /stocks/top-losers`

**Access:** Public

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Success Response (200 OK):**
```json
{
  "losers": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank",
      "price": 850.00,
      "change_percent": -3.1,
      "change_amount": -27.50
    }
  ]
}
```

---

#### Get Most Active Stocks

**Endpoint:** `GET /stocks/most-active`

**Access:** Public

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Success Response (200 OK):**
```json
{
  "most_active": [
    {
      "symbol": "NTC",
      "volume": 5000000,
      "price": 1250.50,
      "trades": 25000
    }
  ]
}
```

---

#### Get Upcoming Market Events

**Endpoint:** `GET /stocks/:symbol/events`

**Access:** Public

**URL Parameters:**
- `symbol`: Company symbol (e.g., "NTC")

**Success Response (200 OK):**
```json
{
  "events": [
    {
      "id": "event_1",
      "company_symbol": "NTC",
      "event_type": "dividend",
      "description": "Dividend payout - 50 NPR per share",
      "scheduled_date": "2026-03-15T00:00:00Z",
      "impact": "positive"
    }
  ]
}
```

---

## Protected Endpoints

**Requires:** `Authorization: Bearer <token>` header

### Authentication - User Profile

#### Get User Profile

**Endpoint:** `GET /auth/profile`

**Access:** Protected (JWT required)

**Success Response (200 OK):**
```json
{
  "user": {
    "id": "user_123",
    "full_name": "John Doe",
    "email": "user@example.com",
    "phone": "+977981234567",
    "kyc_status": "verified",
    "role": "user",
    "created_at": "2026-01-15T08:30:00Z"
  }
}
```

---

#### Update User Profile

**Endpoint:** `PUT /auth/profile/update`

**Access:** Protected (JWT required)

**Request Body:**
```json
{
  "full_name": "John Updated",
  "phone": "+977981234567"
}
```

**Validation:**
- `full_name` (optional): 2-100 characters
- `phone` (optional): 10-15 characters

**Success Response (200 OK):**
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "user_123",
    "full_name": "John Updated",
    "email": "user@example.com",
    "phone": "+977981234567",
    "kyc_status": "verified"
  }
}
```

---

### Virtual Wallet

#### Get Wallet Balance

**Endpoint:** `GET /wallet`

**Access:** Protected (JWT required)

**Success Response (200 OK):**
```json
{
  "wallet": {
    "id": "wallet_123",
    "user_id": "user_123",
    "balance": 1000000.00,
    "total_invested": 150000.00,
    "total_profit_loss": 25000.00,
    "currency": "NPR"
  }
}
```

---

#### Top Up Wallet

**Endpoint:** `POST /wallet/topup`

**Access:** Protected (JWT required)

**Request Body:**
```json
{
  "amount": 50000.00,
  "reference": "TopUp from bank"
}
```

**Validation:**
- `amount`: Minimum 100, Maximum 10,000,000 NPR

**Success Response (200 OK):**
```json
{
  "message": "wallet topped up successfully",
  "new_balance": 1050000.00,
  "transaction_id": "txn_456"
}
```

---

#### Get Transaction History

**Endpoint:** `GET /transaction`

**Access:** Protected (JWT required)

**Query Parameters:**
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Success Response (200 OK):**
```json
{
  "transactions": [
    {
      "id": "txn_123",
      "user_id": "user_123",
      "type": "topup",
      "amount": 50000.00,
      "timestamp": "2026-02-12T10:30:00Z",
      "status": "completed"
    }
  ],
  "limit": 50,
  "offset": 0,
  "total": 120
}
```

---

### Stock Trading

#### Get Trading Wallet

**Endpoint:** `GET /trading/wallet`

**Access:** Protected (JWT required)

**Success Response (200 OK):**
```json
{
  "wallet": {
    "id": "wallet_123",
    "user_id": "user_123",
    "balance": 850000.00,
    "total_invested": 150000.00,
    "total_profit_loss": 25000.00,
    "currency": "NPR",
    "created_at": "2026-01-15T08:30:00Z"
  }
}
```

---

#### Get User Portfolio

**Endpoint:** `GET /trading/portfolio`

**Access:** Protected (JWT required)

**Success Response (200 OK):**
```json
{
  "portfolio": {
    "total_value": 175000.00,
    "total_invested": 150000.00,
    "total_profit_loss": 25000.00,
    "profit_loss_percent": 16.67,
    "holdings": [
      {
        "id": "portfolio_123",
        "company_symbol": "NTC",
        "company_name": "Nepal Telecom",
        "quantity": 100,
        "average_price": 1200.00,
        "current_price": 1250.50,
        "current_value": 125050.00,
        "profit_loss": 5050.00,
        "profit_loss_percent": 4.21
      }
    ]
  }
}
```

---

#### Buy Stock

**Endpoint:** `POST /trading/buy`

**Access:** Protected (JWT required)

**Request Body:**
```json
{
  "symbol": "NTC",
  "quantity": 50
}
```

**Validation:**
- `symbol`: Stock symbol (required)
- `quantity`: Minimum 1 share

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock purchase successful",
  "transaction_id": "stock_txn_123",
  "company_symbol": "NTC",
  "quantity": 50,
  "price_per_share": 1250.50,
  "total_amount": 62525.00,
  "new_balance": 787475.00
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "Insufficient balance in wallet"
}
```

---

#### Sell Stock

**Endpoint:** `POST /trading/sell`

**Access:** Protected (JWT required)

**Request Body:**
```json
{
  "symbol": "NTC",
  "quantity": 30
}
```

**Validation:**
- `symbol`: Stock symbol (required)
- `quantity`: Minimum 1 share, not exceeding owned quantity

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock sale successful",
  "transaction_id": "stock_txn_124",
  "company_symbol": "NTC",
  "quantity": 30,
  "price_per_share": 1250.50,
  "total_amount": 37515.00,
  "profit_loss": 1515.00,
  "new_balance": 824990.00
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "You don't own enough shares to sell"
}
```

---

#### Get Trading Transaction History

**Endpoint:** `GET /trading/transactions`

**Access:** Protected (JWT required)

**Query Parameters:**
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)
- `type` (optional): Filter by transaction type (buy, sell)

**Success Response (200 OK):**
```json
{
  "transactions": [
    {
      "id": "stock_txn_123",
      "user_id": "user_123",
      "company_symbol": "NTC",
      "type": "buy",
      "quantity": 50,
      "price_per_share": 1250.50,
      "total_amount": 62525.00,
      "timestamp": "2026-02-12T10:30:00Z",
      "status": "completed"
    }
  ],
  "limit": 50,
  "offset": 0,
  "total": 45
}
```

---

## Admin Endpoints

**Requires:** `Authorization: Bearer <token>` header + Admin Role

### KYC Management

#### Update User KYC Status

**Endpoint:** `PUT /admin/users/:user_id/kyc`

**Access:** Protected (JWT required) + Admin Role

**URL Parameters:**
- `user_id`: Target user ID

**Request Body:**
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

**Validation:**
- `kyc_status`: "pending", "verified", "rejected", "under_review"
- `role`: "user" or "admin"

**Success Response (200 OK):**
```json
{
  "message": "KYC status updated successfully",
  "user_id": "user_123",
  "kyc_status": "verified",
  "updated_at": "2026-02-12T10:30:00Z"
}
```

---

### Stock Management

#### Seed Stock Data

**Endpoint:** `POST /admin/seed-stocks`

**Access:** Protected (JWT required) + Admin Role

**Request Body:**
```json
{
  "companies": [
    {
      "name": "Nepal Telecom",
      "symbol": "NTC",
      "price": 1250.50,
      "market_cap": "500B",
      "sector": "Telecommunications"
    }
  ]
}
```

**Success Response (200 OK):**
```json
{
  "message": "Stock data seeded successfully",
  "count": 1,
  "companies": [
    {
      "id": "company_1",
      "symbol": "NTC",
      "name": "Nepal Telecom"
    }
  ]
}
```

---

## Data Models

### User Model
```json
{
  "id": "user_123",
  "full_name": "John Doe",
  "email": "user@example.com",
  "phone": "+977981234567",
  "kyc_status": "pending|verified|rejected|under_review",
  "role": "user|admin",
  "created_at": "2026-01-15T08:30:00Z",
  "updated_at": "2026-01-15T08:30:00Z"
}
```

### Virtual Wallet Model
```json
{
  "id": "wallet_123",
  "user_id": "user_123",
  "balance": 1000000.00,
  "total_invested": 150000.00,
  "total_profit_loss": 25000.00,
  "currency": "NPR",
  "created_at": "2026-01-15T08:30:00Z",
  "updated_at": "2026-01-15T08:30:00Z"
}
```

### User Portfolio Model
```json
{
  "id": "portfolio_123",
  "user_id": "user_123",
  "company_id": "company_1",
  "quantity": 100,
  "average_price": 1200.00,
  "total_invested": 120000.00,
  "created_at": "2026-01-15T08:30:00Z",
  "updated_at": "2026-01-15T08:30:00Z"
}
```

### Stock Transaction Model
```json
{
  "id": "stock_txn_123",
  "user_id": "user_123",
  "company_id": "company_1",
  "type": "buy|sell",
  "quantity": 50,
  "price_per_share": 1250.50,
  "total_amount": 62525.00,
  "status": "pending|completed|failed|cancelled",
  "timestamp": "2026-02-12T10:30:00Z"
}
```

---

## Error Handling

### Standard Error Response Format

```json
{
  "error": "Error message describing what went wrong",
  "timestamp": "2026-02-12T10:30:00Z",
  "path": "/api/v1/trading/buy"
}
```

### HTTP Status Codes

| Code | Meaning | Use Case |
|------|---------|----------|
| 200 | OK | Successful request |
| 201 | Created | Resource successfully created |
| 400 | Bad Request | Invalid input or insufficient balance |
| 401 | Unauthorized | Missing or invalid JWT token |
| 403 | Forbidden | Insufficient permissions (not admin) |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Duplicate user email |
| 500 | Server Error | Internal server error |

### Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "invalid credentials" | Wrong email/password | Verify credentials |
| "user already exists" | Email already registered | Use different email |
| "Insufficient balance" | Not enough wallet balance | Top up wallet |
| "You don't own enough shares" | Trying to sell more than owned | Check portfolio |
| "User not authenticated" | Missing JWT token | Add Authorization header |

---

## Examples

### Example 1: Complete Trading Flow

```bash
# 1. Register
curl -X POST https://gold-go-backend.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+977981234567",
    "password": "securepass123",
    "role": "user"
  }'

# 2. Login
TOKEN=$(curl -X POST https://gold-go-backend.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }' | jq -r '.token')

# 3. Check Wallet
curl -X GET https://gold-go-backend.onrender.com/api/v1/trading/wallet \
  -H "Authorization: Bearer $TOKEN"

# 4. Search Stock
curl -X GET https://gold-go-backend.onrender.com/api/v1/stocks/search?q=NTC

# 5. Buy Stock
curl -X POST https://gold-go-backend.onrender.com/api/v1/trading/buy \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NTC",
    "quantity": 50
  }'

# 6. View Portfolio
curl -X GET https://gold-go-backend.onrender.com/api/v1/trading/portfolio \
  -H "Authorization: Bearer $TOKEN"

# 7. Sell Stock
curl -X POST https://gold-go-backend.onrender.com/api/v1/trading/sell \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NTC",
    "quantity": 30
  }'
```

---

## Rate Limiting

Currently no rate limiting is enforced. In production, implement:
- 100 requests per minute per user
- 1000 requests per minute per IP

---

## Support

For API issues or questions:
- **GitHub:** https://github.com/919Umesh/stock_market_sim
- **Issues:** Report via GitHub Issues
