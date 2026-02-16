# Gold Go - API Documentation

Complete API reference for the Gold Go Virtual Stock Market Platform.

**Base URL**: `http://localhost:8081/api/v1` (local) or `https://your-deployment-url/api/v1` (production)

---

## 🔐 Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 📋 Table of Contents

1. [Authentication Endpoints](#authentication-endpoints)
2. [Stock Market Endpoints](#stock-market-endpoints)
3. [Sector Endpoints](#sector-endpoints)
4. [AI Prediction Endpoints](#ai-prediction-endpoints)
5. [User Profile Endpoints](#user-profile-endpoints)
6. [Wallet Endpoints](#wallet-endpoints)
7. [Trading Endpoints](#trading-endpoints)
8. [Admin Endpoints](#admin-endpoints)
9. [Error Responses](#error-responses)

---

## Authentication Endpoints

### Register User
Create a new user account.

**Endpoint**: `POST /auth/register`  
**Authentication**: None (Public)

**Request Body**:
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "9876543210",
  "password": "securepass123",
  "role": "user"
}
```

**Response** (201 Created):
```json
{
  "message": "user registered successfully",
  "user": {
    "$id": "65f8a9b2c3d4e5f6g7h8i9j0",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "kyc_status": "pending",
    "role": "user",
    "$createdAt": "2026-02-16T12:00:00.000Z"
  }
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "password": "securepass123",
    "role": "user"
  }'
```

---

### Login
Authenticate and receive JWT token.

**Endpoint**: `POST /auth/login`  
**Authentication**: None (Public)

**Request Body**:
```json
{
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response** (200 OK):
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "$id": "65f8a9b2c3d4e5f6g7h8i9j0",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "kyc_status": "pending",
    "role": "user"
  }
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

---

## Stock Market Endpoints

### List All Companies
Get a paginated list of all companies.

**Endpoint**: `GET /stocks?limit=25&offset=0`  
**Authentication**: None (Public)

**Query Parameters**:
- `limit` (optional): Number of results (default: 25)
- `offset` (optional): Pagination offset (default: 0)

**Response** (200 OK):
```json
{
  "companies": [
    {
      "$id": "company_id_1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500,
      "is_active": true,
      "$createdAt": "2026-02-16T11:33:50.000Z"
    }
  ],
  "total": 25,
  "limit": 25,
  "offset": 0
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks?limit=10&offset=0
```

---

### Search Companies
Search for companies by name or symbol.

**Endpoint**: `GET /stocks/search?q=nabil`  
**Authentication**: None (Public)

**Query Parameters**:
- `q` (required): Search query (name or symbol)

**Response** (200 OK):
```json
{
  "companies": [
    {
      "$id": "company_id_1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "is_active": true
    }
  ]
}
```

**cURL Example**:
```bash
curl "http://localhost:8081/api/v1/stocks/search?q=bank"
```

---

### Get Company Details
Get detailed information about a specific company.

**Endpoint**: `GET /stocks/:symbol`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "company": {
    "$id": "company_id_1",
    "symbol": "NABIL",
    "name": "Nabil Bank Limited",
    "sector": "Banking",
    "market_cap": 180000000000,
    "description": "Leading private sector bank in Nepal",
    "founded_year": 1984,
    "employees": 3500,
    "is_active": true,
    "$createdAt": "2026-02-16T11:33:50.000Z"
  }
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/NABIL
```

---

### Get Current Stock Price
Get the latest stock price for a company.

**Endpoint**: `GET /stocks/:symbol/price`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "price": {
    "$id": "price_id_1",
    "company_id": "company_id_1",
    "open_price": 1250.50,
    "high_price": 1275.00,
    "low_price": 1240.25,
    "close_price": 1268.75,
    "volume": 2500000,
    "timestamp": "2026-02-16T00:00:00Z",
    "timeframe": "1D"
  }
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/NABIL/price
```

---

### Get Price History
Get historical price data for a company.

**Endpoint**: `GET /stocks/:symbol/history?days=30`  
**Authentication**: None (Public)

**Query Parameters**:
- `days` (optional): Number of days (default: 30, max: 365)

**Response** (200 OK):
```json
{
  "symbol": "NABIL",
  "prices": [
    {
      "$id": "price_id_1",
      "company_id": "company_id_1",
      "open_price": 1250.50,
      "high_price": 1275.00,
      "low_price": 1240.25,
      "close_price": 1268.75,
      "volume": 2500000,
      "timestamp": "2026-02-16T00:00:00Z",
      "timeframe": "1D"
    }
  ]
}
```

**cURL Example**:
```bash
curl "http://localhost:8081/api/v1/stocks/NABIL/history?days=7"
```

---

### Get Market Overview
Get general market statistics.

**Endpoint**: `GET /stocks/market-overview`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "total_companies": 25,
  "total_market_cap": 2500000000000,
  "avg_volume": 1500000,
  "top_gainer": {
    "symbol": "NTC",
    "change_percent": 5.2
  },
  "top_loser": {
    "symbol": "HYDRO1",
    "change_percent": -3.1
  }
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/market-overview
```

---

### Get Top Gainers
Get list of top gaining stocks.

**Endpoint**: `GET /stocks/top-gainers?limit=10`  
**Authentication**: None (Public)

**Query Parameters**:
- `limit` (optional): Number of results (default: 10)

**Response** (200 OK):
```json
{
  "gainers": [
    {
      "symbol": "NTC",
      "name": "Nepal Telecom",
      "current_price": 1850.00,
      "change_percent": 5.2,
      "volume": 3200000
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/top-gainers?limit=5
```

---

### Get Top Losers
Get list of top losing stocks.

**Endpoint**: `GET /stocks/top-losers?limit=10`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "losers": [
    {
      "symbol": "HYDRO1",
      "name": "Upper Tamakoshi Hydropower",
      "current_price": 580.00,
      "change_percent": -3.1,
      "volume": 1800000
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/top-losers?limit=5
```

---

### Get Most Active Stocks
Get most traded stocks by volume.

**Endpoint**: `GET /stocks/most-active?limit=10`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "most_active": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "current_price": 1268.75,
      "volume": 4500000,
      "change_percent": 2.3
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/most-active
```

---

### Get Upcoming Events
Get upcoming market events for a company.

**Endpoint**: `GET /stocks/:symbol/events`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "events": [
    {
      "$id": "event_id_1",
      "company_id": "company_id_1",
      "event_type": "earnings",
      "title": "Q3 Financial Results",
      "description": "Strong quarterly earnings announcement with profit growth",
      "impact_percentage": 3.5,
      "event_date": "2026-02-20T00:00:00Z",
      "image_url": "https://placehold.co/600x400?text=Q3+Financial+Results"
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/stocks/NABIL/events
```

---

## Sector Endpoints

### Get All Sectors
List all available sectors.

**Endpoint**: `GET /sectors`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "sectors": [
    {
      "name": "Banking",
      "company_count": 5,
      "total_market_cap": 700000000000
    },
    {
      "name": "Information Technology",
      "company_count": 5,
      "total_market_cap": 450000000000
    },
    {
      "name": "Hydropower",
      "company_count": 5,
      "total_market_cap": 310000000000
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/sectors
```

---

### Get Companies by Sector
Get all companies in a specific sector.

**Endpoint**: `GET /sectors/:sector/companies`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "sector": "Banking",
  "companies": [
    {
      "$id": "company_id_1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "market_cap": 180000000000,
      "is_active": true
    }
  ]
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/sectors/Banking/companies
```

---

### Get Sector Statistics
Get performance statistics for a sector.

**Endpoint**: `GET /sectors/:sector/stats`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "sector": "Banking",
  "total_companies": 5,
  "total_market_cap": 700000000000,
  "avg_change_percent": 2.1,
  "total_volume": 12000000
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/sectors/Banking/stats
```

---

## AI Prediction Endpoints

### Get Stock Price Prediction
Get AI-based price prediction for a stock (uses Linear Regression on 30 days of historical data).

**Endpoint**: `GET /prediction/:symbol`  
**Authentication**: None (Public)

**Response** (200 OK):
```json
{
  "symbol": "NABIL",
  "predicted_price": 1285.50,
  "rmse": 12.34,
  "mae": 8.76,
  "algorithm": "Linear Regression (Simple)",
  "datapoints": 30
}
```

**cURL Example**:
```bash
curl http://localhost:8081/api/v1/prediction/NABIL
```

---

## User Profile Endpoints

### Get User Profile
Get current user's profile information.

**Endpoint**: `GET /auth/profile`  
**Authentication**: Required (JWT)

**Response** (200 OK):
```json
{
  "user": {
    "$id": "user_id_1",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "file_id_123",
    "$createdAt": "2026-02-15T10:00:00.000Z"
  }
}
```

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/auth/profile
```

---

### Update Profile
Update user profile information.

**Endpoint**: `PUT /auth/profile/update`  
**Authentication**: Required (JWT)

**Request Body**:
```json
{
  "full_name": "John Updated Doe",
  "phone": "9876543211"
}
```

**Response** (200 OK):
```json
{
  "message": "profile updated successfully",
  "user": {
    "$id": "user_id_1",
    "full_name": "John Updated Doe",
    "email": "john@example.com",
    "phone": "9876543211",
    "kyc_status": "verified",
    "role": "user"
  }
}
```

**cURL Example**:
```bash
curl -X PUT -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"full_name": "John Updated Doe", "phone": "9876543211"}' \
  http://localhost:8081/api/v1/auth/profile/update
```

---

### Upload Profile Image
Upload a profile picture (multipart form data).

**Endpoint**: `POST /auth/profile/image`  
**Authentication**: Required (JWT)

**Request**: Multipart form with `image` field  
**Max Size**: 5MB  
**Allowed Types**: jpg, jpeg, png, gif

**Response** (200 OK):
```json
{
  "message": "profile image uploaded successfully",
  "user": {
    "$id": "user_id_1",
    "full_name": "John Doe",
    "email": "john@example.com",
    "profile_image_id": "new_file_id_456"
  }
}
```

**cURL Example**:
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -F "image=@/path/to/profile.jpg" \
  http://localhost:8081/api/v1/auth/profile/image
```

---

## Wallet Endpoints

### Get Wallet Balance
Get user's virtual wallet information.

**Endpoint**: `GET /wallet`  
**Authentication**: Required (JWT)

**Response** (200 OK):
```json
{
  "wallet": {
    "$id": "wallet_id_1",
    "user_id": "user_id_1",
    "balance": 950000.00,
    "total_invested": 50000.00,
    "total_profit_loss": 2500.00,
    "locked": false,
    "$createdAt": "2026-02-15T10:05:00.000Z"
  }
}
```

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/wallet
```

---

### Top Up Wallet
Add virtual currency to wallet.

**Endpoint**: `POST /wallet/topup`  
**Authentication**: Required (JWT)

**Request Body**:
```json
{
  "amount": 100000
}
```

**Response** (200 OK):
```json
{
  "message": "wallet topped up successfully",
  "wallet": {
    "$id": "wallet_id_1",
    "user_id": "user_id_1",
    "balance": 1050000.00,
    "total_invested": 50000.00,
    "total_profit_loss": 2500.00
  },
  "transaction": {
    "$id": "txn_id_1",
    "user_id": "user_id_1",
    "type": "topup",
    "amount": 100000,
    "status": "completed",
    "reference_id": "TOPUP_1234567890"
  }
}
```

**cURL Example**:
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount": 100000}' \
  http://localhost:8081/api/v1/wallet/topup
```

---

### Get Wallet Transactions
Get wallet transaction history (top-ups, etc.).

**Endpoint**: `GET /transaction?limit=50&offset=0`  
**Authentication**: Required (JWT)

**Response** (200 OK):
```json
{
  "transactions": [
    {
      "$id": "txn_id_1",
      "user_id": "user_id_1",
      "type": "topup",
      "amount": 100000,
      "status": "completed",
      "reference_id": "TOPUP_1234567890",
      "$createdAt": "2026-02-16T10:00:00.000Z"
    }
  ],
  "limit": 50,
  "offset": 0
}
```

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8081/api/v1/transaction?limit=20"
```

---

## Trading Endpoints

### Get Trading Wallet
Get wallet information (same as /wallet but under trading context).

**Endpoint**: `GET /trading/wallet`  
**Authentication**: Required (JWT)

**Response**: Same as `GET /wallet`

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/trading/wallet
```

---

### Get Portfolio
Get user's stock portfolio holdings.

**Endpoint**: `GET /trading/portfolio`  
**Authentication**: Required (JWT)

**Response** (200 OK):
```json
{
  "portfolio": [
    {
      "$id": "portfolio_id_1",
      "user_id": "user_id_1",
      "company_id": "company_id_1",
      "quantity": 50,
      "average_price": 1200.00,
      "total_invested": 60000.00,
      "company": {
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "current_price": 1268.75
      },
      "current_value": 63437.50,
      "profit_loss": 3437.50,
      "profit_loss_percent": 5.73
    }
  ],
  "total_invested": 60000.00,
  "current_value": 63437.50,
  "total_profit_loss": 3437.50
}
```

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/trading/portfolio
```

---

### Buy Stock
Purchase shares of a company.

**Endpoint**: `POST /trading/buy`  
**Authentication**: Required (JWT)

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
  "transaction": {
    "$id": "stock_txn_id_1",
    "user_id": "user_id_1",
    "company_id": "company_id_1",
    "type": "buy",
    "quantity": 10,
    "price_per_share": 1268.75,
    "total_amount": 12687.50,
    "status": "completed",
    "$createdAt": "2026-02-16T12:00:00.000Z"
  },
  "wallet": {
    "balance": 937312.50,
    "total_invested": 62687.50
  }
}
```

**Error Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "Insufficient balance",
  "required": 12687.50,
  "available": 5000.00
}
```

**cURL Example**:
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "NABIL", "quantity": 10}' \
  http://localhost:8081/api/v1/trading/buy
```

---

### Sell Stock
Sell shares of a company.

**Endpoint**: `POST /trading/sell`  
**Authentication**: Required (JWT)

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
  "transaction": {
    "$id": "stock_txn_id_2",
    "user_id": "user_id_1",
    "company_id": "company_id_1",
    "type": "sell",
    "quantity": 5,
    "price_per_share": 1280.00,
    "total_amount": 6400.00,
    "status": "completed",
    "$createdAt": "2026-02-16T12:30:00.000Z"
  },
  "wallet": {
    "balance": 943712.50,
    "total_invested": 56687.50
  },
  "profit_loss": 343.75
}
```

**Error Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "Insufficient shares",
  "requested": 10,
  "available": 5
}
```

**cURL Example**:
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "NABIL", "quantity": 5}' \
  http://localhost:8081/api/v1/trading/sell
```

---

### Get Trading History
Get stock transaction history (buy/sell records).

**Endpoint**: `GET /trading/transactions?limit=50&offset=0`  
**Authentication**: Required (JWT)

**Response** (200 OK):
```json
{
  "transactions": [
    {
      "$id": "stock_txn_id_1",
      "user_id": "user_id_1",
      "company_id": "company_id_1",
      "type": "buy",
      "quantity": 10,
      "price_per_share": 1268.75,
      "total_amount": 12687.50,
      "status": "completed",
      "reference_id": "TXN_NABIL_1708084800000",
      "$createdAt": "2026-02-16T12:00:00.000Z",
      "company": {
        "symbol": "NABIL",
        "name": "Nabil Bank Limited"
      }
    }
  ],
  "limit": 50,
  "offset": 0
}
```

**cURL Example**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8081/api/v1/trading/transactions?limit=20"
```

---

## Admin Endpoints

### Update User KYC Status
Update a user's KYC status and role (admin only).

**Endpoint**: `PUT /admin/users/:user_id/kyc`  
**Authentication**: Required (JWT + Admin Role)

**Request Body**:
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

**Valid KYC Statuses**: `pending`, `verified`, `rejected`, `under_review`  
**Valid Roles**: `user`, `admin`

**Response** (200 OK):
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "$id": "user_id_1",
    "full_name": "John Doe",
    "email": "john@example.com",
    "kyc_status": "verified",
    "role": "user"
  }
}
```

**cURL Example**:
```bash
curl -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"kyc_status": "verified", "role": "user"}' \
  http://localhost:8081/api/v1/admin/users/user_id_123/kyc
```

---

### Seed Stock Data
Trigger stock data seeding (creates 25 companies, 30 days of prices, 10 events).

**Endpoint**: `POST /admin/seed-stocks`  
**Authentication**: Required (JWT + Admin Role)

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Database seeded successfully",
  "companies_created": 25,
  "prices_created": 750,
  "events_created": 10
}
```

**cURL Example**:
```bash
curl -X POST -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/seed-stocks
```

---

## Error Responses

All error responses follow a consistent format:

### 400 Bad Request
```json
{
  "error": "Validation failed",
  "details": "quantity must be at least 1"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing token"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden",
  "message": "Admin access required"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Company with symbol 'INVALID' not found"
}
```

### 409 Conflict
```json
{
  "error": "Conflict",
  "message": "User already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred"
}
```

---

## Testing Workflow Example

Here's a complete workflow to test the API:

```bash
# 1. Register a new user
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "9876543210",
    "password": "test123",
    "role": "user"
  }'

# 2. Login and save token
TOKEN=$(curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123"
  }' | jq -r '.token')

# 3. View profile
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/auth/profile

# 4. Check wallet (should have 1,000,000 NPR initial balance)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/trading/wallet

# 5. Browse stocks
curl http://localhost:8081/api/v1/stocks?limit=5

# 6. Get current price for NABIL
curl http://localhost:8081/api/v1/stocks/NABIL/price

# 7. Buy 10 shares of NABIL
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "NABIL", "quantity": 10}' \
  http://localhost:8081/api/v1/trading/buy

# 8. Check portfolio
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/trading/portfolio

# 9. Sell 5 shares
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "NABIL", "quantity": 5}' \
  http://localhost:8081/api/v1/trading/sell

# 10. View transaction history
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/v1/trading/transactions

# 11. Get AI prediction
curl http://localhost:8081/api/v1/prediction/NABIL
```

---

## Data Volumes

After running the seed script (`go run scripts/setup_appwrite.go`):

- **Companies**: 25 Nepalese companies across 6 sectors
- **Stock Prices**: ~750 records (30 days × 25 companies)
- **Market Events**: 250 events (10 per company)
- **Test Transactions**: 750 stock transactions (30 per company)

---

## Notes

- All monetary values are in NPR (Nepalese Rupees)
- Initial wallet balance: 1,000,000 NPR
- Stock prices are simulated for educational purposes
- ML predictions use Linear Regression on 30 days of historical data
- Timestamps are in ISO 8601 format (UTC)

---

**Last Updated**: February 16, 2026  
**API Version**: v1
