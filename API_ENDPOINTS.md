# 📊 Gold Go - Stock Market Simulation API Documentation

A complete REST API for a stock market simulation platform built with **Go**, **Gin Framework**, and **Appwrite** as the backend database. This platform allows users to trade virtual stocks, manage portfolios, and track market performance.

---

## 🌐 Base URL

```
http://localhost:8080/api/v1
```

For production/deployed instances, replace `localhost:8080` with your actual domain.

---

## 🔐 Authentication

Most endpoints require authentication using **JWT Bearer Token**.

### **How to Authenticate:**
1. Register a new user or login with existing credentials
2. Copy the `token` from the login response
3. Include it in the `Authorization` header for protected endpoints

**Header Format:**
```
Authorization: Bearer <YOUR_JWT_TOKEN>
```

---

## 📋 Table of Contents

1. [Authentication Endpoints](#1-authentication-endpoints)
2. [Stock Market - Public Endpoints](#2-stock-market---public-endpoints)
3. [Wallet Management - Protected](#3-wallet-management---protected)
4. [Trading Operations - Protected](#4-trading-operations---protected)
5. [Admin Endpoints](#5-admin-endpoints)
6. [Error Responses](#error-responses)

---

## 1. Authentication Endpoints

### 1.1 Register New User

**Endpoint:** `POST /auth/register`  
**Auth Required:** ❌ No

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "9800000000",
  "password": "securePassword123",
  "role": "user"
}
```

**Field Validations:**
- `full_name`: Required, 2-100 characters
- `email`: Required, valid email format
- `phone`: Required, 10-15 digits
- `password`: Required, minimum 6 characters
- `role`: Required, must be "user" or "admin"

**Response (201 Created):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "6759abc123def456",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9800000000",
    "role": "user",
    "kyc_status": "pending",
    "created_at": "2026-02-09T10:30:00Z",
    "updated_at": "2026-02-09T10:30:00Z"
  }
}
```

**Possible Errors:**
- `409 Conflict`: User with this email already exists
- `400 Bad Request`: Invalid input data

---

### 1.2 Login

**Endpoint:** `POST /auth/login`  
**Auth Required:** ❌ No

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "6759abc123def456",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9800000000",
    "role": "user",
    "kyc_status": "pending"
  }
}
```

**Possible Errors:**
- `401 Unauthorized`: Invalid credentials

---

### 1.3 Get User Profile

**Endpoint:** `GET /auth/profile`  
**Auth Required:** ✅ Yes

**Response (200 OK):**
```json
{
  "user": {
    "id": "6759abc123def456",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9800000000",
    "role": "user",
    "kyc_status": "verified",
    "created_at": "2026-02-09T10:30:00Z",
    "updated_at": "2026-02-09T10:30:00Z"
  }
}
```

---

### 1.4 Update User Profile

**Endpoint:** `PUT /auth/profile/update`  
**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "full_name": "John Updated Doe",
  "phone": "9801234567"
}
```

**Response (200 OK):**
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "6759abc123def456",
    "full_name": "John Updated Doe",
    "phone": "9801234567",
    "email": "john.doe@example.com"
  }
}
```

---

## 2. Stock Market - Public Endpoints

These endpoints are publicly accessible and don't require authentication.

### 2.1 List All Companies

**Endpoint:** `GET /stocks`  
**Auth Required:** ❌ No

**Query Parameters:**
- `limit` (optional): Number of results, default: 50
- `offset` (optional): Pagination offset, default: 0

**Example Request:**
```
GET /api/v1/stocks?limit=10&offset=0
```

**Response (200 OK):**
```json
{
  "companies": [
    {
      "id": "comp001",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 50000000000,
      "description": "Leading commercial bank in Nepal",
      "founded_year": 1984,
      "employees": 1200,
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

---

### 2.2 Search Companies

**Endpoint:** `GET /stocks/search`  
**Auth Required:** ❌ No

**Query Parameters:**
- `q` (required): Search query (company name or symbol)

**Example Request:**
```
GET /api/v1/stocks/search?q=NABIL
```

**Response (200 OK):**
```json
{
  "companies": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking"
    }
  ]
}
```

---

### 2.3 Get Company Details

**Endpoint:** `GET /stocks/:symbol`  
**Auth Required:** ❌ No

**Example Request:**
```
GET /api/v1/stocks/NABIL
```

**Response (200 OK):**
```json
{
  "company": {
    "id": "comp001",
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

### 2.4 Get Current Stock Price

**Endpoint:** `GET /stocks/:symbol/price`  
**Auth Required:** ❌ No

**Example Request:**
```
GET /api/v1/stocks/NABIL/price
```

**Response (200 OK):**
```json
{
  "price": {
    "company_id": "comp001",
    "symbol": "NABIL",
    "price": 1245.50,
    "open": 1240.00,
    "high": 1250.00,
    "low": 1235.00,
    "volume": 125000,
    "change_percent": 0.44,
    "timestamp": "2026-02-09T15:30:00Z"
  }
}
```

---

### 2.5 Get Price History

**Endpoint:** `GET /stocks/:symbol/history`  
**Auth Required:** ❌ No

**Query Parameters:**
- `timeframe` (optional): "1m", "5m", "15m", "1h", "1d", "1w", "1M" - default: "1d"
- `days` (optional): Number of days to fetch, default: 30

**Example Request:**
```
GET /api/v1/stocks/NABIL/history?timeframe=1d&days=7
```

**Response (200 OK):**
```json
{
  "symbol": "NABIL",
  "timeframe": "1d",
  "prices": [
    {
      "price": 1245.50,
      "open": 1240.00,
      "high": 1250.00,
      "low": 1235.00,
      "volume": 125000,
      "timestamp": "2026-02-09T00:00:00Z"
    }
  ]
}
```

---

### 2.6 Get Market Overview

**Endpoint:** `GET /stocks/market-overview`  
**Auth Required:** ❌ No

**Response (200 OK):**
```json
{
  "total_companies": 23,
  "top_gainers": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "current_price": 1245.50,
      "change_percent": 3.45
    }
  ],
  "top_losers": [
    {
      "symbol": "SCB",
      "name": "Standard Chartered Bank",
      "current_price": 845.00,
      "change_percent": -2.15
    }
  ],
  "most_active": [
    {
      "symbol": "NICA",
      "name": "NIC Asia Bank",
      "volume": 250000
    }
  ]
}
```

---

### 2.7 Get Top Gainers

**Endpoint:** `GET /stocks/top-gainers`  
**Auth Required:** ❌ No

**Query Parameters:**
- `limit` (optional): Number of results, default: 10

**Response (200 OK):**
```json
{
  "gainers": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "current_price": 1245.50,
      "previous_price": 1203.25,
      "change_percent": 3.51
    }
  ]
}
```

---

### 2.8 Get Top Losers

**Endpoint:** `GET /stocks/top-losers`  
**Auth Required:** ❌ No

**Query Parameters:**
- `limit` (optional): Number of results, default: 10

**Response (200 OK):**
```json
{
  "losers": [
    {
      "symbol": "SCB",
      "name": "Standard Chartered Bank",
      "current_price": 845.00,
      "previous_price": 863.50,
      "change_percent": -2.14
    }
  ]
}
```

---

### 2.9 Get Most Active Stocks

**Endpoint:** `GET /stocks/most-active`  
**Auth Required:** ❌ No

**Query Parameters:**
- `limit` (optional): Number of results, default: 10

**Response (200 OK):**
```json
{
  "most_active": [
    {
      "symbol": "NICA",
      "name": "NIC Asia Bank",
      "volume": 250000,
      "current_price": 980.00
    }
  ]
}
```

---

### 2.10 Get Upcoming Events

**Endpoint:** `GET /stocks/:symbol/events`  
**Auth Required:** ❌ No

**Example Request:**
```
GET /api/v1/stocks/NABIL/events
```

**Response (200 OK):**
```json
{
  "events": [
    {
      "company_id": "comp001",
      "event_type": "dividend_announcement",
      "title": "Dividend Declaration",
      "description": "15% cash dividend announced",
      "event_date": "2026-03-15T00:00:00Z",
      "impact": "positive"
    }
  ]
}
```

---

## 3. Wallet Management - Protected

All wallet endpoints require authentication.

### 3.1 Get Wallet Balance

**Endpoint:** `GET /wallet`  
**Auth Required:** ✅ Yes

**Response (200 OK):**
```json
{
  "wallet": {
    "id": "wallet001",
    "user_id": "6759abc123def456",
    "fiat_balance": 50000.00,
    "locked": false,
    "version": 5,
    "created_at": "2026-02-01T10:00:00Z",
    "updated_at": "2026-02-09T14:30:00Z"
  }
}
```

---

### 3.2 Top Up Wallet

**Endpoint:** `POST /wallet/topup`  
**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "amount": 10000.00
}
```

**Validation:**
- `amount`: Required, must be greater than 0

**Response (200 OK):**
```json
{
  "message": "top-up successful",
  "wallet": {
    "id": "wallet001",
    "user_id": "6759abc123def456",
    "fiat_balance": 60000.00,
    "locked": false
  },
  "transaction": {
    "id": "txn001",
    "user_id": "6759abc123def456",
    "type": "credit",
    "amount": 10000.00,
    "status": "completed",
    "reference_id": "topup_uuid-here",
    "created_at": "2026-02-09T15:00:00Z"
  }
}
```

**Possible Errors:**
- `400 Bad Request`: Invalid amount (must be > 0)
- `423 Locked`: Wallet is locked

---

### 3.3 Get Transaction History

**Endpoint:** `GET /transaction`  
**Auth Required:** ✅ Yes

**Response (200 OK):**
```json
{
  "message": "user transactions retrieved successfully",
  "data": [
    {
      "id": "txn001",
      "user_id": "6759abc123def456",
      "type": "credit",
      "amount": 10000.00,
      "status": "completed",
      "reference_id": "topup_uuid",
      "created_at": "2026-02-09T15:00:00Z"
    },
    {
      "id": "txn002",
      "user_id": "6759abc123def456",
      "type": "debit",
      "amount": 5000.00,
      "status": "completed",
      "reference_id": "stock_purchase_uuid",
      "created_at": "2026-02-09T16:00:00Z"
    }
  ]
}
```

---

## 4. Trading Operations - Protected

### 4.1 Get Virtual Trading Wallet

**Endpoint:** `GET /trading/wallet`  
**Auth Required:** ✅ Yes

**Response (200 OK):**
```json
{
  "wallet": {
    "id": "vwallet001",
    "user_id": "6759abc123def456",
    "balance": 45000.00,
    "total_invested": 15000.00,
    "total_profit_loss": 2500.00,
    "created_at": "2026-02-01T10:00:00Z"
  }
}
```

---

### 4.2 Get Portfolio

**Endpoint:** `GET /trading/portfolio`  
**Auth Required:** ✅ Yes

**Response (200 OK):**
```json
{
  "total_value": 47500.00,
  "total_invested": 15000.00,
  "total_profit_loss": 2500.00,
  "profit_loss_percent": 16.67,
  "holdings": [
    {
      "symbol": "NABIL",
      "quantity": 10,
      "average_price": 1200.00,
      "current_price": 1245.50,
      "total_invested": 12000.00,
      "current_value": 12455.00,
      "profit_loss": 455.00,
      "profit_loss_percent": 3.79
    },
    {
      "symbol": "NICA",
      "quantity": 3,
      "average_price": 980.00,
      "current_price": 1015.00,
      "total_invested": 2940.00,
      "current_value": 3045.00,
      "profit_loss": 105.00,
      "profit_loss_percent": 3.57
    }
  ]
}
```

---

### 4.3 Buy Stock

**Endpoint:** `POST /trading/buy`  
**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "symbol": "NABIL",
  "quantity": 5
}
```

**Validations:**
- `symbol`: Required
- `quantity`: Required, minimum 1

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock purchased successfully",
  "transaction": {
    "id": "trade001",
    "user_id": "6759abc123def456",
    "symbol": "NABIL",
    "type": "buy",
    "quantity": 5,
    "price_per_unit": 1245.50,
    "total_amount": 6227.50,
    "status": "completed",
    "created_at": "2026-02-09T16:30:00Z"
  },
  "wallet_balance": 38772.50
}
```

**Possible Errors:**
- `400 Bad Request`: Insufficient funds, invalid quantity
- `404 Not Found`: Stock symbol not found

---

### 4.4 Sell Stock

**Endpoint:** `POST /trading/sell`  
**Auth Required:** ✅ Yes

**Request Body:**
```json
{
  "symbol": "NABIL",
  "quantity": 2
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock sold successfully",
  "transaction": {
    "id": "trade002",
    "user_id": "6759abc123def456",
    "symbol": "NABIL",
    "type": "sell",
    "quantity": 2,
    "price_per_unit": 1245.50,
    "total_amount": 2491.00,
    "status": "completed",
    "created_at": "2026-02-09T17:00:00Z"
  },
  "wallet_balance": 41263.50
}
```

**Possible Errors:**
- `400 Bad Request`: Insufficient stock quantity
- `404 Not Found`: Stock not found in portfolio

---

### 4.5 Get Transaction History

**Endpoint:** `GET /trading/transactions`  
**Auth Required:** ✅ Yes

**Query Parameters:**
- `limit` (optional): Number of results, default: 50
- `offset` (optional): Pagination offset, default: 0

**Example Request:**
```
GET /api/v1/trading/transactions?limit=20&offset=0
```

**Response (200 OK):**
```json
{
  "transactions": [
    {
      "id": "trade001",
      "user_id": "6759abc123def456",
      "symbol": "NABIL",
      "type": "buy",
      "quantity": 5,
      "price_per_unit": 1245.50,
      "total_amount": 6227.50,
      "status": "completed",
      "created_at": "2026-02-09T16:30:00Z"
    }
  ],
  "limit": 20,
  "offset": 0
}
```

---

## 5. Admin Endpoints

Admin endpoints require both authentication and admin role.

### 5.1 Update User KYC Status

**Endpoint:** `PUT /admin/users/:user_id/kyc`  
**Auth Required:** ✅ Yes (Admin only)

**Request Body:**
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

**Valid Values:**
- `kyc_status`: "pending", "verified", "rejected", "under_review"
- `role`: "user", "admin"

**Response (200 OK):**
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": "6759abc123def456",
    "full_name": "John Doe",
    "kyc_status": "verified",
    "role": "user"
  }
}
```

---

### 5.2 Seed Stock Data

**Endpoint:** `POST /admin/seed-stocks`  
**Auth Required:** ✅ Yes (Admin only)

**Description:** Populates the database with initial stock market data including companies, price history, and market events. **This should only be run once** after initial deployment.

**Request Body:** None required

**Response (200 OK):**
```json
{
  "message": "Stock data seeded successfully",
  "companies_created": 23,
  "prices_created": 690,
  "events_created": 46
}
```

**Possible Errors:**
- `400 Bad Request`: Database already seeded (delete via Appwrite console if re-seeding is needed)

---

## Error Responses

All errors follow a consistent format:

### Standard Error Response
```json
{
  "error": "Error type or code",
  "message": "Human-readable error message"
}
```

### Common HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | Success | Request completed successfully |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input data or validation failed |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | User lacks required permissions |
| 404 | Not Found | Requested resource doesn't exist |
| 409 | Conflict | Resource already exists |
| 423 | Locked | Wallet is locked |
| 500 | Internal Server Error | Server-side error occurred |

---

## 🔧 Development & Testing

### Health Check
```
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy"
}
```

### Environment Setup

1. **Clone and Install:**
```bash
git clone <repository-url>
cd gold_go
go mod download
```

2. **Configure Environment Variables:**
Create a `.env` file with:
```env
PORT=8080
JWT_SECRET=your-secret-key-here

# Appwrite Configuration
APPWRITE_ENDPOINT=https://cloud.appwrite.io/v1
APPWRITE_PROJECT_ID=your-project-id
APPWRITE_API_KEY=your-api-key
APPWRITE_DATABASE_ID=your-database-id

# Worker Configuration
WORKER_COUNT=5
QUEUE_SIZE=100
```

3. **Run the Server:**
```bash
go run cmd/main.go
```

---

## 📱 Mobile/Frontend Integration

### Basic Flow

1. **Register/Login** → Get JWT token
2. **Store token** securely (localStorage, secure storage)
3. **Include token** in all protected API requests
4. **Handle errors** gracefully (401 = re-login required)

### Sample Request (JavaScript/Fetch)
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/profile', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${yourJwtToken}`,
    'Content-Type': 'application/json'
  }
});

const data = await response.json();
```

---

## 🛡️ Security Notes

- All passwords are hashed using bcrypt
- JWT tokens are used for stateless authentication
- Admin routes are protected with role-based middleware
- Input validation is enforced on all endpoints
- CORS is enabled for development (configure for production)

---

## 📊 Database Collections (Appwrite)

The API uses the following Appwrite collections:

1. **users** - User accounts and profiles
2. **wallets** - User wallet balances
3. **transactions** - Wallet transaction history
4. **companies** - Stock company information
5. **stock_prices** - Historical stock prices
6. **market_events** - Corporate events and announcements
7. **virtual_wallets** - Trading account balances
8. **user_portfolios** - User stock holdings
9. **stock_transactions** - Buy/sell trade history

---

## 📞 Support

For issues or questions:
- Review error messages carefully
- Check authentication tokens
- Verify request payload format
- Ensure Appwrite is properly configured

---

**Last Updated:** February 9, 2026  
**API Version:** 1.0.0  
**Powered by:** Go 1.24 + Gin + Appwrite
