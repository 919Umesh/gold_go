# Gold Go - Stock Market Simulation API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Table of Contents
1. [Authentication](#authentication)
2. [Stock Management](#stock-management)
3. [Sector Information](#sector-information)
4. [Stock Predictions](#stock-predictions)
5. [Wallet Management](#wallet-management)
6. [Trading Operations](#trading-operations)
7. [Admin Operations](#admin-operations)
8. [Health Check](#health-check)

---

## Authentication

### 1. Register User
**POST** `/api/v1/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "9841234567",
  "role": "user",
  "password": "securePassword123"
}
```

**Validation Rules:**
- `full_name`: Required, min 2, max 100 characters
- `email`: Required, valid email format
- `phone`: Required, min 10, max 15 characters
- `role`: Required, min 3, max 10 characters (e.g., "user", "admin")
- `password`: Required, minimum 6 characters

**Response (201 Created):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "uuid-string",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9841234567",
    "kyc_status": "pending",
    "role": "user",
    "profile_image_id": "",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input data
- `409 Conflict`: User already exists
- `500 Internal Server Error`: Registration failed

---

### 2. Login
**POST** `/api/v1/auth/login`

Authenticate a user and receive a JWT token.

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
    "id": "uuid-string",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9841234567",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Invalid credentials

---

### 3. Get Profile
**GET** `/api/v1/auth/profile`

Get the authenticated user's profile.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "uuid-string",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9841234567",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "image-id",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `404 Not Found`: User not found
- `500 Internal Server Error`: Invalid user ID type

---

### 4. Update Profile
**PUT** `/api/v1/auth/profile/update`

Update the authenticated user's profile information.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "phone": "9841234568"
}
```

**Validation Rules:**
- `full_name`: Optional, min 2, max 100 characters
- `phone`: Optional, min 10, max 15 characters

**Response (200 OK):**
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "uuid-string",
    "full_name": "John Doe Updated",
    "email": "john.doe@example.com",
    "phone": "9841234568",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T12:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: No fields to update or invalid data
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Profile update failed

---

### 5. Upload Profile Image
**POST** `/api/v1/auth/profile/image`

Upload a profile image for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Request Body (Form Data):**
- `image`: Image file (max 5MB)

**Response (200 OK):**
```json
{
  "message": "profile image uploaded successfully",
  "user": {
    "id": "uuid-string",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9841234567",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "uploaded-image-id",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T12:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Image file is required or image too large (max 5MB)
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to upload image

---

## Stock Management

### 6. List Companies
**GET** `/api/v1/stocks`

Get a paginated list of all companies.

**Query Parameters:**
- `limit` (optional): Number of results per page (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example:**
```
GET /api/v1/stocks?limit=20&offset=0
```

**Response (200 OK):**
```json
{
  "companies": [
    {
      "id": "uuid-string",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": "180000000000",
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500,
      "total_shares": 1800000000,
      "available_shares": 1500000000,
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-02-17T00:00:00Z"
    }
  ],
  "limit": 20,
  "offset": 0,
  "count": 20,
  "total": 150
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch companies

---

### 7. Search Companies
**GET** `/api/v1/stocks/search`

Search for companies by name or symbol.

**Query Parameters:**
- `q` (required): Search query string

**Example:**
```
GET /api/v1/stocks/search?q=nabil
```

**Response (200 OK):**
```json
{
  "companies": [
    {
      "id": "uuid-string",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": "180000000000",
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500,
      "total_shares": 1800000000,
      "available_shares": 1500000000,
      "is_active": true,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-02-17T00:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Search query required
- `500 Internal Server Error`: Search failed

---

### 8. Get Company Details
**GET** `/api/v1/stocks/:symbol`

Get detailed information about a specific company.

**Path Parameters:**
- `symbol`: Company stock symbol (e.g., NABIL)

**Example:**
```
GET /api/v1/stocks/NABIL
```

**Response (200 OK):**
```json
{
  "company": {
    "id": "uuid-string",
    "symbol": "NABIL",
    "name": "Nabil Bank Limited",
    "sector": "Banking",
    "market_cap": "180000000000",
    "description": "Leading private sector bank in Nepal",
    "founded_year": 1984,
    "employees": 3500,
    "total_shares": 1800000000,
    "available_shares": 1500000000,
    "is_active": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-02-17T00:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Symbol parameter required
- `404 Not Found`: Company not found

---

### 9. Get Current Stock Price
**GET** `/api/v1/stocks/:symbol/price`

Get the current price information for a stock.

**Path Parameters:**
- `symbol`: Company stock symbol

**Example:**
```
GET /api/v1/stocks/NABIL/price
```

**Response (200 OK):**
```json
{
  "price": {
    "id": "uuid-string",
    "company_id": "company-uuid",
    "open_price": "950.50",
    "high_price": "975.00",
    "low_price": "945.00",
    "close_price": "970.00",
    "volume": 125000,
    "timestamp": "2026-02-17T16:00:00Z",
    "timeframe": "1D",
    "created_at": "2026-02-17T16:00:00Z",
    "company": {
      "id": "company-uuid",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Symbol parameter required
- `404 Not Found`: Price not found

---

### 10. Get Price History
**GET** `/api/v1/stocks/:symbol/history`

Get historical price data for a stock.

**Path Parameters:**
- `symbol`: Company stock symbol

**Query Parameters:**
- `timeframe` (optional): Time interval (default: "1D")
- `days` (optional): Number of days to retrieve (default: 30, max: 365)

**Example:**
```
GET /api/v1/stocks/NABIL/history?timeframe=1D&days=30
```

**Response (200 OK):**
```json
{
  "symbol": "NABIL",
  "timeframe": "1D",
  "prices": [
    {
      "id": "uuid-string",
      "company_id": "company-uuid",
      "open_price": "950.50",
      "high_price": "975.00",
      "low_price": "945.00",
      "close_price": "970.00",
      "volume": 125000,
      "timestamp": "2026-02-17T00:00:00Z",
      "timeframe": "1D",
      "created_at": "2026-02-17T16:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Symbol parameter required
- `404 Not Found`: Price history not found

---

### 11. Get Market Overview
**GET** `/api/v1/stocks/market-overview`

Get overall market statistics and overview.

**Response (200 OK):**
```json
{
  "total_companies": 150,
  "total_market_cap": "5000000000000",
  "total_volume": 5000000,
  "total_trades": 12500,
  "advancing": 85,
  "declining": 50,
  "unchanged": 15,
  "avg_change_percentage": "1.25"
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch market overview

---

### 12. Get Top Gainers
**GET** `/api/v1/stocks/top-gainers`

Get the list of stocks with the highest price gains.

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Example:**
```
GET /api/v1/stocks/top-gainers?limit=5
```

**Response (200 OK):**
```json
{
  "gainers": [
    {
      "company": {
        "id": "uuid-string",
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "sector": "Banking"
      },
      "current_price": "970.00",
      "change": "45.50",
      "change_percentage": "4.92"
    }
  ]
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch top gainers

---

### 13. Get Top Losers
**GET** `/api/v1/stocks/top-losers`

Get the list of stocks with the highest price losses.

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Example:**
```
GET /api/v1/stocks/top-losers?limit=5
```

**Response (200 OK):**
```json
{
  "losers": [
    {
      "company": {
        "id": "uuid-string",
        "symbol": "EXAMPLE",
        "name": "Example Company",
        "sector": "Technology"
      },
      "current_price": "450.00",
      "change": "-25.50",
      "change_percentage": "-5.36"
    }
  ]
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch top losers

---

### 14. Get Most Active Stocks
**GET** `/api/v1/stocks/most-active`

Get stocks with the highest trading volume.

**Query Parameters:**
- `limit` (optional): Number of results (default: 10)

**Example:**
```
GET /api/v1/stocks/most-active?limit=10
```

**Response (200 OK):**
```json
{
  "active": [
    {
      "company": {
        "id": "uuid-string",
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "sector": "Banking"
      },
      "current_price": "970.00",
      "volume": 250000,
      "change_percentage": "2.50"
    }
  ]
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch most active

---

### 15. Get Upcoming Events
**GET** `/api/v1/stocks/:symbol/events`

Get upcoming market events for a specific company.

**Path Parameters:**
- `symbol`: Company stock symbol

**Example:**
```
GET /api/v1/stocks/NABIL/events
```

**Response (200 OK):**
```json
{
  "events": [
    {
      "id": "uuid-string",
      "company_id": "company-uuid",
      "event_type": "earnings",
      "title": "Q1 Financial Results",
      "description": "Strong quarterly earnings with revenue growth",
      "impact_percentage": 3.5,
      "event_date": "2026-03-15T00:00:00Z",
      "created_at": "2026-02-10T00:00:00Z"
    },
    {
      "id": "uuid-string",
      "company_id": "company-uuid",
      "event_type": "dividend",
      "title": "Annual Dividend Announcement",
      "description": "Annual dividend payout to shareholders",
      "impact_percentage": 2.0,
      "event_date": "2026-04-01T00:00:00Z",
      "created_at": "2026-02-12T00:00:00Z"
    }
  ]
}
```

**Event Types:**
- `earnings`: Earnings reports
- `news`: General news
- `dividend`: Dividend announcements
- `merger`: Merger/acquisition news
- `ipo`: IPO/FPO announcements
- `split`: Stock split announcements

**Error Responses:**
- `404 Not Found`: Events not available

---

## Sector Information

### 16. Get All Sectors
**GET** `/api/v1/sectors`

Get a list of all available sectors.

**Response (200 OK):**
```json
{
  "message": "sectors retrieved successfully",
  "sectors": [
    "Banking",
    "Information Technology",
    "Hydropower",
    "Insurance",
    "Pharma",
    "Manufacturing",
    "Real Estate"
  ],
  "count": 7
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch sectors

---

### 17. Get Companies by Sector
**GET** `/api/v1/sectors/:sector/companies`

Get all companies in a specific sector.

**Path Parameters:**
- `sector`: Sector name (e.g., Banking, Technology)

**Query Parameters:**
- `limit` (optional): Number of results per page (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example:**
```
GET /api/v1/sectors/Banking/companies?limit=10&offset=0
```

**Response (200 OK):**
```json
{
  "sector": "Banking",
  "companies": [
    {
      "id": "uuid-string",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": "180000000000",
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500,
      "total_shares": 1800000000,
      "available_shares": 1500000000,
      "is_active": true
    }
  ],
  "limit": 10,
  "offset": 0,
  "count": 10,
  "total": 25
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch companies

---

### 18. Get Sector Statistics
**GET** `/api/v1/sectors/:sector/stats`

Get statistical information about a specific sector.

**Path Parameters:**
- `sector`: Sector name

**Example:**
```
GET /api/v1/sectors/Banking/stats
```

**Response (200 OK):**
```json
{
  "message": "sector statistics",
  "sector": "Banking",
  "statistics": {
    "company_count": 25,
    "total_market_cap": "2500000000000",
    "avg_market_cap": "100000000000",
    "total_employees": 45000,
    "avg_employees": 1800,
    "avg_founded_year": 1995
  },
  "top_5_companies": [
    {
      "id": "uuid-string",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": "180000000000",
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500
    }
  ]
}
```

**Error Responses:**
- `404 Not Found`: No companies found in this sector
- `500 Internal Server Error`: Failed to fetch sector stats

---

## Stock Predictions

### 19. Get Stock Price Prediction
**GET** `/api/v1/prediction/:symbol`

Get AI/ML-based price prediction for a stock.

**Path Parameters:**
- `symbol`: Company stock symbol

**Example:**
```
GET /api/v1/prediction/NABIL
```

**Response (200 OK):**
```json
{
  "prediction": {
    "id": "uuid-string",
    "company_id": "company-uuid",
    "predicted_price": 985.50,
    "confidence_score": 0.85,
    "prediction_date": "2026-02-17T00:00:00Z",
    "target_date": "2026-02-18T00:00:00Z",
    "model_used": "linear_regression",
    "created_at": "2026-02-17T08:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Symbol is required
- `500 Internal Server Error`: Failed to predict price

---

## Wallet Management

### 20. Get Wallet
**GET** `/api/v1/wallet`

Get the authenticated user's wallet information.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "wallet": {
    "id": "uuid-string",
    "user_id": "user-uuid",
    "fiat_balance": "50000.00",
    "locked_amount": "0.00",
    "is_active": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-02-17T00:00:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to fetch wallet

---

### 21. Top Up Wallet
**POST** `/api/v1/wallet/topup`

Add funds to the user's wallet.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "amount": 10000.00
}
```

**Validation Rules:**
- `amount`: Required, must be greater than 0

**Response (200 OK):**
```json
{
  "message": "top-up successful",
  "wallet": {
    "id": "uuid-string",
    "user_id": "user-uuid",
    "fiat_balance": "60000.00",
    "locked_amount": "0.00",
    "is_active": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  },
  "transaction": {
    "id": "transaction-uuid",
    "user_id": "user-uuid",
    "type": "topup",
    "amount": "10000.00",
    "status": "success",
    "reference_id": "topup_xxxxx-xxxx-xxxx",
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid amount
- `401 Unauthorized`: User not authenticated
- `423 Locked`: Wallet is locked
- `500 Internal Server Error`: Top-up failed

---

### 22. Get User Transactions
**GET** `/api/v1/transaction`

Get all wallet transactions for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "message": "user transactions retrieved successfully",
  "data": [
    {
      "id": "transaction-uuid",
      "user_id": "user-uuid",
      "type": "topup",
      "amount": "10000.00",
      "status": "success",
      "reference_id": "topup_xxxxx-xxxx-xxxx",
      "created_at": "2026-02-17T10:00:00Z",
      "updated_at": "2026-02-17T10:00:00Z"
    },
    {
      "id": "transaction-uuid-2",
      "user_id": "user-uuid",
      "type": "topup",
      "amount": "50000.00",
      "status": "success",
      "reference_id": "topup_yyyyy-yyyy-yyyy",
      "created_at": "2026-01-15T12:00:00Z",
      "updated_at": "2026-01-15T12:00:00Z"
    }
  ]
}
```

**Transaction Types:**
- `topup`: Wallet top-up
- `refund`: Refund to wallet

**Transaction Status:**
- `pending`: Transaction pending
- `success`: Transaction completed successfully
- `failed`: Transaction failed

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to fetch transactions

---

## Trading Operations

### 23. Get Trading Wallet
**GET** `/api/v1/trading/wallet`

Get the trading wallet for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "wallet": {
    "id": "uuid-string",
    "user_id": "user-uuid",
    "balance": "45000.00",
    "total_invested": "5000.00",
    "total_profit_loss": "250.00",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-02-17T00:00:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to get wallet

---

### 24. Get Portfolio
**GET** `/api/v1/trading/portfolio`

Get the user's stock portfolio with current values.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "portfolio": [
    {
      "id": "portfolio-uuid",
      "user_id": "user-uuid",
      "company_id": "company-uuid",
      "quantity": 10,
      "average_price": "950.00",
      "total_invested": "9500.00",
      "current_price": "970.00",
      "current_value": "9700.00",
      "profit_loss": "200.00",
      "profit_loss_percentage": "2.11",
      "created_at": "2026-02-10T00:00:00Z",
      "updated_at": "2026-02-17T00:00:00Z",
      "company": {
        "id": "company-uuid",
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "sector": "Banking"
      }
    }
  ],
  "summary": {
    "total_invested": "9500.00",
    "current_value": "9700.00",
    "total_profit_loss": "200.00",
    "total_profit_loss_percentage": "2.11"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to get portfolio

---

### 25. Buy Stock
**POST** `/api/v1/trading/buy`

Purchase stocks for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "symbol": "NABIL",
  "quantity": 10
}
```

**Validation Rules:**
- `symbol`: Required
- `quantity`: Required, minimum 1

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock purchased successfully",
  "transaction": {
    "id": "transaction-uuid",
    "user_id": "user-uuid",
    "company_id": "company-uuid",
    "type": "buy",
    "quantity": 10,
    "price_per_share": "970.00",
    "total_amount": "9700.00",
    "status": "completed",
    "reference_id": "BUY_xxxxx-xxxx-xxxx",
    "created_at": "2026-02-17T10:30:00Z"
  },
  "portfolio": {
    "id": "portfolio-uuid",
    "user_id": "user-uuid",
    "company_id": "company-uuid",
    "quantity": 10,
    "average_price": "970.00",
    "total_invested": "9700.00",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  },
  "wallet_balance": "40300.00"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input or insufficient funds/shares
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Purchase failed

---

### 26. Sell Stock
**POST** `/api/v1/trading/sell`

Sell stocks from the user's portfolio.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "symbol": "NABIL",
  "quantity": 5
}
```

**Validation Rules:**
- `symbol`: Required
- `quantity`: Required, minimum 1

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock sold successfully",
  "transaction": {
    "id": "transaction-uuid",
    "user_id": "user-uuid",
    "company_id": "company-uuid",
    "type": "sell",
    "quantity": 5,
    "price_per_share": "980.00",
    "total_amount": "4900.00",
    "status": "completed",
    "reference_id": "SELL_xxxxx-xxxx-xxxx",
    "created_at": "2026-02-17T11:00:00Z"
  },
  "portfolio": {
    "id": "portfolio-uuid",
    "user_id": "user-uuid",
    "company_id": "company-uuid",
    "quantity": 5,
    "average_price": "970.00",
    "total_invested": "4850.00",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T11:00:00Z"
  },
  "wallet_balance": "45200.00",
  "profit_loss": "50.00"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input or insufficient stocks
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Sale failed

---

### 27. Get Transaction History
**GET** `/api/v1/trading/transactions`

Get the trading transaction history for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `limit` (optional): Number of results per page (default: 50)
- `offset` (optional): Number of results to skip (default: 0)

**Example:**
```
GET /api/v1/trading/transactions?limit=20&offset=0
```

**Response (200 OK):**
```json
{
  "transactions": [
    {
      "id": "transaction-uuid",
      "user_id": "user-uuid",
      "company_id": "company-uuid",
      "type": "buy",
      "quantity": 10,
      "price_per_share": "970.00",
      "total_amount": "9700.00",
      "status": "completed",
      "reference_id": "BUY_xxxxx-xxxx-xxxx",
      "created_at": "2026-02-17T10:30:00Z",
      "company": {
        "id": "company-uuid",
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "sector": "Banking"
      }
    },
    {
      "id": "transaction-uuid-2",
      "user_id": "user-uuid",
      "company_id": "company-uuid",
      "type": "sell",
      "quantity": 5,
      "price_per_share": "980.00",
      "total_amount": "4900.00",
      "status": "completed",
      "reference_id": "SELL_xxxxx-xxxx-xxxx",
      "created_at": "2026-02-17T11:00:00Z",
      "company": {
        "id": "company-uuid",
        "symbol": "NABIL",
        "name": "Nabil Bank Limited",
        "sector": "Banking"
      }
    }
  ],
  "limit": 20,
  "offset": 0
}
```

**Transaction Types:**
- `buy`: Stock purchase
- `sell`: Stock sale

**Transaction Status:**
- `pending`: Transaction pending
- `completed`: Transaction completed
- `failed`: Transaction failed
- `cancelled`: Transaction cancelled

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to get transactions

---

## Admin Operations

### 28. Update User KYC Status
**PUT** `/api/v1/admin/users/:user_id/kyc`

Update the KYC status and role of a user (Admin only).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `user_id`: User ID to update

**Request Body:**
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

**Validation Rules:**
- `kyc_status`: Required, must be one of: pending, verified, rejected, under_review
- `role`: Required, must be one of: user, admin

**Response (200 OK):**
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": "user-uuid",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "9841234567",
    "kyc_status": "verified",
    "role": "user",
    "profile_image_id": "",
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T14:00:00Z"
  }
}
```

**KYC Status Values:**
- `pending`: KYC pending review
- `verified`: KYC verified
- `rejected`: KYC rejected
- `under_review`: KYC under review

**Error Responses:**
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: KYC update failed

---

### 29. Seed Stock Data
**POST** `/api/v1/admin/seed-stocks`

Seed initial stock data into the database (Admin only).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{}
```

**Response (200 OK):**
```json
{
  "message": "Stock data seeded successfully",
  "companies_added": 25,
  "prices_added": 750,
  "events_added": 125
}
```

**Error Responses:**
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Seeding failed

---

## Health Check

### 30. Health Check
**GET** `/health`

Check if the API server is running.

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

## Error Response Format

All error responses follow this general format:

```json
{
  "error": "Error message description"
}
```

Or for validation errors:

```json
{
  "error": "Field validation error details"
}
```

---

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

To obtain a token, use the `/api/v1/auth/login` endpoint.

---

## Data Types

### Decimal Values
All monetary values and prices are returned as string representations of decimal numbers to maintain precision:
```json
{
  "price": "970.50",
  "balance": "50000.00"
}
```

### Dates
All dates follow ISO 8601 format:
```json
{
  "created_at": "2026-02-17T10:30:00Z"
}
```

---

## Rate Limiting

Currently, there are no rate limits implemented. This may change in future versions.

---

## Pagination

Endpoints that return lists support pagination through `limit` and `offset` query parameters:
- `limit`: Number of items to return (default and maximum values vary by endpoint)
- `offset`: Number of items to skip

Example:
```
GET /api/v1/stocks?limit=20&offset=40
```

---

## Support

For questions or issues, please contact the development team or refer to the project repository.

---

## Changelog

**Version 1.0.0** (Current)
- Initial API release
- Complete authentication system
- Stock management and trading features
- Wallet management
- Admin operations
- Market predictions using ML

---

*Last Updated: February 17, 2026*
