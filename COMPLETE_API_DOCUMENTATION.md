# Stock Market Simulator - Complete API Documentation

**Base URL:** `http://localhost:8080/api/v1`

**API Version:** v1

**Last Updated:** March 11, 2026

---

## Table of Contents
1. [Health Check](#health-check)
2. [Authentication](#authentication)
3. [Admin Management](#admin-management)
4. [Wallet Management](#wallet-management)
5. [Trading & Orders](#trading--orders)
6. [IPO Management](#ipo-management)
7. [Market Data](#market-data)
8. [Events](#events)
9. [Price Triggers](#price-triggers)

---

## Health Check

### Check API Health
**Endpoint:** `GET /health`

**Authentication:** Not Required

**Response (Status 200):**
```json
{
  "status": "ok",
  "service": "stock-market-simulator"
}
```

---

## Authentication

### Register New User
**Endpoint:** `POST /auth/register`

**Authentication:** Not Required

**Request Body:**
```json
{
  "full_name": "string (required, 2-100 chars)",
  "email": "string (required, valid email)",
  "phone": "string (required, 10-15 digits)",
  "role": "string (required, 3-10 chars) - e.g., 'user' or 'admin'",
  "password": "string (required, min 6 chars)"
}
```

**Response (Status 201):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `201` - User registered successfully
- `400` - Invalid request data
- `409` - User already exists

---

### User Login
**Endpoint:** `POST /auth/login`

**Authentication:** Not Required

**Request Body:**
```json
{
  "email": "string (required, valid email)",
  "password": "string (required)"
}
```

**Response (Status 200):**
```json
{
  "message": "login successful",
  "token": "JWT token string",
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Login successful
- `400` - Invalid request data
- `401` - Invalid credentials

---

### Get User Profile
**Endpoint:** `GET /auth/profile`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `404` - User not found

---

### Update User Profile
**Endpoint:** `PUT /auth/profile/update`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "full_name": "string (optional, 2-100 chars)",
  "phone": "string (optional, 10-15 digits)"
}
```

**Response (Status 200):**
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Profile updated
- `400` - No fields to update or invalid data
- `401` - Unauthorized

---

### Upload Profile Image
**Endpoint:** `POST /auth/profile/image`

**Authentication:** Required (Bearer Token)

**Content-Type:** `multipart/form-data`

**Parameters:**
- `image` (file, required) - Max 5MB

**Response (Status 200):**
```json
{
  "message": "profile image uploaded successfully",
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string (URL to uploaded image)",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Image uploaded
- `400` - Invalid file or missing image
- `401` - Unauthorized

---

## Admin Management

### Get All Users
**Endpoint:** `GET /admin/users`

**Authentication:** Required (Bearer Token, Admin Only)

**Query Parameters:**
- `limit` (int, optional, default: 50) - Number of users to return
- `offset` (int, optional, default: 0) - Number of users to skip

**Response (Status 200):**
```json
{
  "users": [
    {
      "id": "string",
      "full_name": "string",
      "email": "string",
      "phone": "string",
      "kyc_status": "pending|verified|rejected|under_review",
      "role": "string",
      "profile_image_url": "string",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `403` - Admin access required

---

### Update User KYC & Role
**Endpoint:** `PUT /admin/users/{user_id}/kyc`

**Authentication:** Required (Bearer Token, Admin Only)

**URL Parameters:**
- `user_id` (string, required) - The user ID to update

**Request Body:**
```json
{
  "kyc_status": "string (required) - pending|verified|rejected|under_review",
  "role": "string (required) - user|admin"
}
```

**Response (Status 200):**
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": "string",
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "kyc_status": "pending|verified|rejected|under_review",
    "role": "string",
    "profile_image_url": "string",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - KYC updated
- `400` - Invalid request data
- `403` - Admin access required

---

### Create Company
**Endpoint:** `POST /admin/companies`

**Authentication:** Required (Bearer Token, Admin Only)

**Request Body:**
```json
{
  "symbol": "string (required, e.g., 'ACME')",
  "name": "string (required, e.g., 'ACME Corporation')",
  "sector": "string (optional, defaults to 'General')",
  "total_supply": "integer (required, > 0)"
}
```

**Response (Status 201):**
```json
{
  "message": "company created",
  "company": {
    "id": "string",
    "symbol": "string",
    "name": "string",
    "sector": "string",
    "description": "string",
    "total_supply": 1000000,
    "shares_outstanding": "decimal",
    "current_price": "decimal",
    "market_cap": "decimal",
    "eps": "decimal",
    "pe_ratio": "decimal",
    "book_value": "decimal",
    "pbv": "decimal",
    "week_52_high": "decimal",
    "week_52_low": "decimal",
    "avg_120_day": "decimal",
    "yield_1_year": "decimal",
    "listed_date": "ISO8601 date",
    "is_active": true,
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `201` - Company created
- `400` - Invalid request data
- `403` - Admin access required

---

### Launch IPO
**Endpoint:** `POST /admin/ipos`

**Authentication:** Required (Bearer Token, Admin Only)

**Request Body:**
```json
{
  "company_id": "string (required)",
  "price_per_share": "string (required, decimal as string, e.g., '100.50')",
  "total_shares": "integer (required, > 0)",
  "max_per_applicant": "integer (required, > 0)",
  "open_at": "string (required, RFC3339 format, e.g., '2026-03-15T09:30:00Z')",
  "close_at": "string (required, RFC3339 format, e.g., '2026-03-20T15:30:00Z')"
}
```

**Response (Status 201):**
```json
{
  "message": "IPO launched",
  "ipo": {
    "id": "string",
    "company_id": "string",
    "price_per_share": "decimal",
    "total_shares": 100000,
    "allocated_shares": 0,
    "max_per_applicant": 1000,
    "open_at": "ISO8601 timestamp",
    "close_at": "ISO8601 timestamp",
    "status": "pending|open|closed|allocated",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `201` - IPO launched
- `400` - Invalid request data (invalid dates, price format, etc.)
- `403` - Admin access required

---

### Allocate IPO Shares
**Endpoint:** `POST /admin/ipos/{id}/allocate`

**Authentication:** Required (Bearer Token, Admin Only)

**URL Parameters:**
- `id` (string, required) - The IPO ID

**Request Body:** None (POST with empty body)

**Response (Status 200):**
```json
{
  "message": "IPO allocation complete",
  "result": {
    "ipo_id": "string",
    "total_applications": 100,
    "total_shares_allocated": 50000,
    "total_shares_not_allocated": 50000,
    "refunds_processed": 25,
    "timestamp": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Allocation completed
- `400` - Invalid IPO ID or IPO not in allocated state
- `403` - Admin access required

---

### Get IPO Applications
**Endpoint:** `GET /admin/ipos/{id}/applications`

**Authentication:** Required (Bearer Token, Admin Only)

**URL Parameters:**
- `id` (string, required) - The IPO ID

**Response (Status 200):**
```json
{
  "ipo_id": "string",
  "count": 5,
  "applications": [
    {
      "id": "string",
      "ipo_id": "string",
      "user_id": "string",
      "shares_requested": 1000,
      "shares_allocated": 0,
      "amount_paid": "decimal",
      "amount_refunded": "decimal",
      "status": "pending|allocated|not_allocated|refunded",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `403` - Admin access required
- `404` - IPO not found

---

## Wallet Management

### Get All Wallets
**Endpoint:** `GET /wallet/`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "main_wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  },
  "trading_wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "locked_balance": "decimal (reserved for pending buy orders)",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Get Main Wallet
**Endpoint:** `GET /wallet/main`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Get Trading Wallet
**Endpoint:** `GET /wallet/trading`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "locked_balance": "decimal (reserved for pending buy orders)",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Top Up Main Wallet
**Endpoint:** `POST /wallet/topup`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "amount": "string (required, decimal as string, e.g., '1000.50')"
}
```

**Response (Status 200):**
```json
{
  "message": "top-up successful",
  "wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Top-up successful
- `400` - Invalid amount
- `401` - Unauthorized

---

### Transfer Between Wallets
**Endpoint:** `POST /wallet/transfer`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "amount": "string (required, decimal as string, e.g., '500.00')",
  "direction": "string (required) - main_to_trading|trading_to_main"
}
```

**Response (Status 200):**
```json
{
  "message": "transfer successful",
  "main_wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  },
  "trading_wallet": {
    "id": "string",
    "user_id": "string",
    "balance": "decimal",
    "locked_balance": "decimal",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Transfer successful
- `400` - Invalid amount or direction
- `401` - Unauthorized
- `422` - Insufficient balance

---

### Get Transfer History
**Endpoint:** `GET /wallet/transfers`

**Authentication:** Required (Bearer Token)

**Parameters:** None (returns last 50 transfers)

**Response (Status 200):**
```json
{
  "transfers": [
    {
      "id": "string",
      "user_id": "string",
      "amount": "decimal",
      "direction": "main_to_trading|trading_to_main",
      "status": "string",
      "created_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

## Trading & Orders

### Place Buy Order
**Endpoint:** `POST /orders/buy`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "company_id": "string (required)",
  "quantity": "integer (required, > 0)",
  "price": "string (required, decimal as string, e.g., '150.75')",
  "order_type": "string (optional) - limit|market (defaults to 'limit')"
}
```

**Response (Status 201):**
```json
{
  "message": "buy order placed",
  "order": {
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "side": "buy",
    "order_type": "limit|market",
    "price": "decimal",
    "quantity": 100,
    "filled_qty": 0,
    "status": "open|partially_filled|filled|cancelled",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  },
  "matches": 0
}
```

**Status Codes:**
- `201` - Order placed
- `400` - Invalid request data or insufficient balance
- `401` - Unauthorized

---

### Place Sell Order
**Endpoint:** `POST /orders/sell`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "company_id": "string (required)",
  "quantity": "integer (required, > 0)",
  "price": "string (required, decimal as string, e.g., '150.75')"
}
```

**Response (Status 201):**
```json
{
  "message": "sell order placed",
  "order": {
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "side": "sell",
    "order_type": "limit",
    "price": "decimal",
    "quantity": 50,
    "filled_qty": 0,
    "status": "open|partially_filled|filled|cancelled",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  },
  "matches": 0
}
```

**Status Codes:**
- `201` - Order placed
- `400` - Invalid request data or insufficient shares
- `401` - Unauthorized

---

### Cancel Order
**Endpoint:** `PUT /orders/{id}/cancel`

**Authentication:** Required (Bearer Token)

**URL Parameters:**
- `id` (string, required) - The order ID

**Request Body:** None (PUT with empty body)

**Response (Status 200):**
```json
{
  "message": "order cancelled"
}
```

**Status Codes:**
- `200` - Order cancelled
- `400` - Order ID required or order not cancellable
- `401` - Unauthorized

---

### Get Order Book
**Endpoint:** `GET /orders/book/{company_id}`

**Authentication:** Not required

**URL Parameters:**
- `company_id` (string, required) - The company ID

**Response (Status 200):**
```json
{
  "order_book": {
    "symbol": "string",
    "company_id": "string",
    "bids": [
      {
        "price": "decimal",
        "quantity": 1000,
        "orders_count": 5
      }
    ],
    "asks": [
      {
        "price": "decimal",
        "quantity": 1500,
        "orders_count": 8
      }
    ],
    "last_updated": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `400` - Company ID required
- `500` - Server error

---

### Get User Orders
**Endpoint:** `GET /orders/my-orders`

**Authentication:** Required (Bearer Token)

**Parameters:** None (returns last 50 orders)

**Response (Status 200):**
```json
{
  "orders": [
    {
      "id": "string",
      "user_id": "string",
      "company_id": "string",
      "side": "buy|sell",
      "order_type": "limit|market",
      "price": "decimal",
      "quantity": 100,
      "filled_qty": 50,
      "status": "open|partially_filled|filled|cancelled",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Get Portfolio
**Endpoint:** `GET /orders/portfolio`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "portfolio": [
    {
      "id": "string",
      "user_id": "string",
      "company_id": "string",
      "quantity": 500,
      "avg_buy_price": "decimal",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Get User Trades
**Endpoint:** `GET /orders/trades`

**Authentication:** Required (Bearer Token)

**Parameters:** None (returns last 50 trades)

**Response (Status 200):**
```json
{
  "trades": [
    {
      "id": "string",
      "company_id": "string",
      "buy_order_id": "string",
      "sell_order_id": "string",
      "buyer_id": "string",
      "seller_id": "string",
      "price": "decimal",
      "quantity": 100,
      "total_amount": "decimal",
      "created_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Get Company Trade History
**Endpoint:** `GET /market/companies/{id}/trades`

**Authentication:** Not required

**URL Parameters:**
- `id` (string, required) - The company ID

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of matched trades to process

**Description:** Returns the public transaction history for a specific company. Each trade is represented as two records (a buy and a sell) to show both sides of the transaction clearly.

**Response (Status 200):**
```json
{
  "company_id": "string",
  "count": 20,
  "transactions": [
    {
      "id": "trade-uuid-buy",
      "type": "buy",
      "price": "decimal",
      "quantity": 100,
      "created_at": "ISO8601 timestamp"
    },
    {
      "id": "trade-uuid-sell",
      "type": "sell",
      "price": "decimal",
      "quantity": 100,
      "created_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

## IPO Management

### List All IPOs
**Endpoint:** `GET /ipo/`

**Authentication:** Not required

**Parameters:** None

**Response (Status 200):**
```json
{
  "ipos": [
    {
      "id": "string",
      "company_id": "string",
      "price_per_share": "decimal",
      "total_shares": 50000,
      "allocated_shares": 0,
      "max_per_applicant": 1000,
      "open_at": "ISO8601 timestamp",
      "close_at": "ISO8601 timestamp",
      "status": "pending|open|closed|allocated",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ]
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get IPO Details
**Endpoint:** `GET /ipo/{id}`

**Authentication:** Not required

**URL Parameters:**
- `id` (string, required) - The IPO ID

**Response (Status 200):**
```json
{
  "ipo": {
    "id": "string",
    "company_id": "string",
    "price_per_share": "decimal",
    "total_shares": 50000,
    "allocated_shares": 0,
    "max_per_applicant": 1000,
    "open_at": "ISO8601 timestamp",
    "close_at": "ISO8601 timestamp",
    "status": "pending|open|closed|allocated",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `404` - IPO not found

---

### Apply for IPO
**Endpoint:** `POST /ipo/{id}/apply`

**Authentication:** Required (Bearer Token)

**URL Parameters:**
- `id` (string, required) - The IPO ID

**Request Body:**
```json
{
  "shares_requested": "integer (required, > 0)"
}
```

**Response (Status 201):**
```json
{
  "message": "IPO application submitted",
  "application": {
    "id": "string",
    "ipo_id": "string",
    "user_id": "string",
    "shares_requested": 500,
    "shares_allocated": 0,
    "amount_paid": "decimal",
    "amount_refunded": "decimal",
    "status": "pending|allocated|not_allocated|refunded",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `201` - Application submitted
- `400` - Invalid request data or IPO not open
- `401` - Unauthorized

---

## Market Data

### List All Companies
**Endpoint:** `GET /market/companies`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of results
- `offset` (integer, optional, default: 0) - Number of results to skip

**Response (Status 200):**
```json
{
  "data": [
    {
      "id": "string",
      "symbol": "string",
      "name": "string",
      "sector": "string",
      "description": "string",
      "total_supply": 1000000,
      "shares_outstanding": "decimal",
      "current_price": "decimal",
      "market_cap": "decimal",
      "eps": "decimal",
      "pe_ratio": "decimal",
      "book_value": "decimal",
      "pbv": "decimal",
      "week_52_high": "decimal",
      "week_52_low": "decimal",
      "avg_120_day": "decimal",
      "yield_1_year": "decimal",
      "listed_date": "ISO8601 date",
      "is_active": true,
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Recently Listed Companies
**Endpoint:** `GET /market/companies/new`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* Company array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Oldest Listed Companies
**Endpoint:** `GET /market/companies/old`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* Company array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Company Detail
**Endpoint:** `GET /market/companies/{id}`

**Authentication:** Not required

**URL Parameters:**
- `id` (string, required) - The company ID

**Response (Status 200):**
```json
{
  "data": {
    "id": "string",
    "symbol": "string",
    "name": "string",
    "sector": "string",
    "description": "string",
    "total_supply": 1000000,
    "shares_outstanding": "decimal",
    "current_price": "decimal",
    "market_cap": "decimal",
    "eps": "decimal",
    "pe_ratio": "decimal",
    "book_value": "decimal",
    "pbv": "decimal",
    "week_52_high": "decimal",
    "week_52_low": "decimal",
    "avg_120_day": "decimal",
    "yield_1_year": "decimal",
    "listed_date": "ISO8601 date",
    "is_active": true,
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `404` - Company not found

---

### Get Price Prediction
**Endpoint:** `GET /market/companies/{id}/prediction`

**Authentication:** Not required

**URL Parameters:**
- `id` (string, required) - The company ID

**Description:** Predicts the next likely stock price for a specific company using a **Weighted Moving Average (WMA)** algorithm based on the last 5 trades. This is a common academic approach for short-term trend analysis.

**Response (Status 200):**
```json
{
  "data": {
    "company_id": "string",
    "current_price": "decimal",
    "predicted_price": "decimal",
    "expected_change": "decimal",
    "algorithm": "Weighted Moving Average (WMA)",
    "confidence": 0.9,
    "sample_count": 5
  }
}
```

**Status Codes:**
- `200` - Success
- `404` - Company not found

---

### Get Live Trading Data
**Endpoint:** `GET /market/live`

**Authentication:** Not required

**Parameters:** None

**Response (Status 200):**
```json
{
  "data": [
    {
      "symbol": "string",
      "company_id": "string",
      "company_name": "string",
      "sector": "string",
      "ltp": "decimal (Last Traded Price)",
      "change_percent": "decimal",
      "open": "decimal",
      "high": "decimal",
      "low": "decimal",
      "volume": 1000000,
      "previous_close": "decimal",
      "difference": "decimal",
      "turnover": "decimal",
      "last_updated": "ISO8601 timestamp"
    }
  ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Market Index
**Endpoint:** `GET /market/index`

**Authentication:** Not required

**Parameters:** None

**Response (Status 200):**
```json
{
  "data": {
    "index_value": "decimal",
    "change": "decimal",
    "change_percent": "decimal",
    "total_turnover": "decimal",
    "total_volume": 10000000,
    "total_market_cap": "decimal",
    "advances": 25,
    "declines": 10,
    "unchanged": 5,
    "total_companies": 40,
    "previous_close": "decimal",
    "timestamp": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Candlestick Data
**Endpoint:** `GET /market/candlestick`

**Authentication:** Not required

**Query Parameters:**
- `symbol` (string, required) - Stock symbol
- `timeframe` (string, optional, default: "1D") - 1m|5m|15m|1h|1D
- `days` (integer, optional, default: 90) - Number of days of data

**Response (Status 200):**
```json
{
  "data": [
    {
      "timestamp": "ISO8601 timestamp",
      "open": "decimal",
      "high": "decimal",
      "low": "decimal",
      "close": "decimal",
      "volume": 1000000,
      "turnover": "decimal",
      "change_percent": "decimal"
    }
  ],
  "symbol": "string",
  "timeframe": "string",
  "count": 90
}
```

**Status Codes:**
- `200` - Success
- `400` - Symbol required
- `500` - Server error

---

### Get Top Gainers
**Endpoint:** `GET /market/top-gainers`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* LiveTradingData array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Top Losers
**Endpoint:** `GET /market/top-losers`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* LiveTradingData array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Most Active Stocks
**Endpoint:** `GET /market/most-active`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* LiveTradingData array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Top Turnover Stocks
**Endpoint:** `GET /market/top-turnover`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 10) - Number of results

**Response (Status 200):**
```json
{
  "data": [ /* LiveTradingData array */ ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Sector Performance
**Endpoint:** `GET /market/sectors`

**Authentication:** Not required

**Parameters:** None

**Response (Status 200):**
```json
{
  "data": [
    {
      "sector": "string",
      "company_count": 10,
      "avg_change_percent": "decimal",
      "total_turnover": "decimal",
      "total_volume": 5000000,
      "total_market_cap": "decimal"
    }
  ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Companies by Sector
**Endpoint:** `GET /market/sectors/{sector}/companies`

**Authentication:** Not required

**URL Parameters:**
- `sector` (string, required) - Sector name (e.g., "Technology", "Finance")

**Response (Status 200):**
```json
{
  "data": [ /* LiveTradingData array */ ],
  "sector": "string",
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Stream Live Prices (SSE)
**Endpoint:** `GET /market/stream`

**Authentication:** Not required

**Content-Type:** `text/event-stream`

**Description:** Server-Sent Events stream for real-time price updates. Connect and receive price_update events.

**Response Events:**
```
event: price_update
data: {
  "symbol": "string",
  "company_id": "string",
  "company_name": "string",
  "ltp": "decimal",
  "change_percent": "decimal",
  "volume": 1000,
  "timestamp": "ISO8601 timestamp"
}
```

**Status Codes:**
- `200` - Stream connected
- `500` - Server error

---

## Events

### Get All Events
**Endpoint:** `GET /market/events`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of results
- `offset` (integer, optional, default: 0) - Number of results to skip

**Response (Status 200):**
```json
{
  "events": [
    {
      "id": "string",
      "company_id": "string",
      "event_type": "agm|dividend|bonus_share|right_share|quarterly_report|board_meeting|financial_results|stock_split|merger_acquisition|ipo_announcement",
      "title": "string",
      "description": "string",
      "event_date": "ISO8601 timestamp",
      "fiscal_year": "string",
      "status": "upcoming|ongoing|completed|cancelled",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Company Events
**Endpoint:** `GET /market/events/company/{company_id}`

**Authentication:** Not required

**URL Parameters:**
- `company_id` (string, required) - The company ID

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of results

**Response (Status 200):**
```json
{
  "events": [ /* CompanyEvent array */ ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Upcoming Events
**Endpoint:** `GET /market/events/upcoming`

**Authentication:** Not required

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of results

**Response (Status 200):**
```json
{
  "events": [ /* CompanyEvent array with status='upcoming' */ ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

### Get Events by Type
**Endpoint:** `GET /market/events/type/{event_type}`

**Authentication:** Not required

**URL Parameters:**
- `event_type` (string, required) - agm|dividend|bonus_share|right_share|quarterly_report|board_meeting|financial_results|stock_split|merger_acquisition|ipo_announcement

**Query Parameters:**
- `limit` (integer, optional, default: 50) - Number of results

**Response (Status 200):**
```json
{
  "events": [ /* CompanyEvent array filtered by type */ ],
  "count": 50
}
```

**Status Codes:**
- `200` - Success
- `500` - Server error

---

## Price Triggers

### Create Price Trigger
**Endpoint:** `POST /market/triggers`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "company_id": "string (required)",
  "trigger_price": "decimal (required, as decimal object or number)",
  "shares_qty": "integer (required, > 0)",
  "direction": "string (required) - above|below"
}
```

**Response (Status 201):**
```json
{
  "data": {
    "id": "string",
    "user_id": "string",
    "company_id": "string",
    "trigger_price": "decimal",
    "shares_qty": 100,
    "direction": "above|below",
    "status": "active|triggered|cancelled",
    "created_at": "ISO8601 timestamp",
    "updated_at": "ISO8601 timestamp"
  }
}
```

**Status Codes:**
- `201` - Trigger created
- `400` - Invalid request data
- `401` - Unauthorized

---

### Get User Triggers
**Endpoint:** `GET /market/triggers`

**Authentication:** Required (Bearer Token)

**Parameters:** None

**Response (Status 200):**
```json
{
  "data": [
    {
      "id": "string",
      "user_id": "string",
      "company_id": "string",
      "trigger_price": "decimal",
      "shares_qty": 100,
      "direction": "above|below",
      "status": "active|triggered|cancelled",
      "created_at": "ISO8601 timestamp",
      "updated_at": "ISO8601 timestamp"
    }
  ],
  "count": 10
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Server error

---

### Cancel Price Trigger
**Endpoint:** `PUT /market/triggers/{id}/cancel`

**Authentication:** Required (Bearer Token)

**URL Parameters:**
- `id` (string, required) - The trigger ID

**Request Body:** None (PUT with empty body)

**Response (Status 200):**
```json
{
  "message": "trigger cancelled"
}
```

**Status Codes:**
- `200` - Trigger cancelled
- `401` - Unauthorized
- `500` - Server error

---

## Authentication Header Format

For all endpoints requiring authentication, include the JWT token in the Authorization header:

```
Authorization: Bearer <JWT_TOKEN>
```

Example:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Or for authentication/validation errors:

```json
{
  "message": "error message"
}
```

---

## Common Status Codes

| Code | Description |
|------|-------------|
| `200` | Success (GET, PUT requests) |
| `201` | Created (POST requests) |
| `400` | Bad Request (invalid parameters) |
| `401` | Unauthorized (authentication required/failed) |
| `403` | Forbidden (insufficient permissions, e.g., admin-only) |
| `404` | Not Found (resource doesn't exist) |
| `409` | Conflict (resource already exists) |
| `422` | Unprocessable Entity (validation failed, e.g., insufficient balance) |
| `500` | Internal Server Error |

---

## Data Types

- **string**: Text value
- **integer**: Whole number (no decimals)
- **decimal**: Decimal number (use string format in requests: "100.50")
- **ISO8601 timestamp**: Date-time format (e.g., "2026-03-11T13:02:02Z")
- **ISO8601 date**: Date only (e.g., "2026-03-11")
- **boolean**: true or false

---

## Rate Limiting

Currently no rate limiting is enforced, but it is recommended to implement rate limiting for production.

---

## Pagination

For endpoints that support pagination:
- `limit`: Number of results per page (default varies by endpoint)
- `offset`: Number of results to skip (default: 0)

Example:
```
GET /market/companies?limit=25&offset=50
```

---

## Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

---

**End of API Documentation**
