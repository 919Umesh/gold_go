# Share Market Simulator — Comprehensive API Documentation (v2.0)

This document provides exhaustive details for every endpoint, including URLs, authentications, required parameters, and example JSON request/response payloads.

**Base URL**: `http://localhost:8080/api/v1`  
**Standard Headers**: 
- `Content-Type: application/json`
- `Authorization: Bearer <JWT_TOKEN>` (for protected routes)

**Authentication**: Bearer token from login endpoint required for protected endpoints.

---

## 1. Health & System

### 1.1 Health Check
`GET /health`

**Response (200 OK):**
```json
{
  "status": "ok",
  "version": "2.0"
}
```

---

### 1.2 Swagger Documentation
`GET /swagger/*any`

Access the interactive Swagger UI at: `http://localhost:8080/swagger/index.html`

---

## 2. Authentication (Public)

### 2.1 Register User
**Endpoint**: `POST /auth/register`

**Description**: Create a new user account with email and password

**Request Headers**:
```
Content-Type: application/json
```

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `full_name` | String | Yes | Min 2, Max 100 chars | User's full name |
| `email` | String | Yes | Valid email format | Must be unique |
| `phone` | String | Yes | 10-15 digits | Phone number |
| `password` | String | Yes | Min 6 characters | User's password |
| `role` | String | Yes | `user` or `admin` | User role (typically `user`) |

**Example Request**:
```json
{
  "full_name": "Sita Shrestha",
  "email": "sita@example.com",
  "phone": "9841234567",
  "password": "mypassword123",
  "role": "user"
}
```

**Response (201 Created)**:
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha",
    "email": "sita@example.com",
    "phone": "9841234567",
    "role": "user",
    "kyc_status": "pending"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input (missing fields, invalid email, short password)
- **409 Conflict**: User with this email already exists

---

### 2.2 Login
**Endpoint**: `POST /auth/login`

**Description**: Authenticate user and generate JWT token

**Request Body Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `email` | String | Yes | Registered email |
| `password` | String | Yes | User's password |

**Example Request**:
```json
{
  "email": "sita@example.com",
  "password": "mypassword123"
}
```

**Response (200 OK)**:
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3ZjhiOWUxMC1hMWIyLWMzZDQtZTVmNi03ZzhoOWkwajFrMmwiLCJleHAiOjE3NDE1OTM4NTd9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha",
    "email": "sita@example.com",
    "role": "user"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Missing email or password
- **401 Unauthorized**: Invalid credentials

---

### 2.3 Get Profile
**Endpoint**: `GET /auth/profile` ⚠️ **Auth Required**

**Description**: Get current authenticated user's profile information

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha",
    "email": "sita@example.com",
    "phone": "9841234567",
    "role": "user",
    "kyc_status": "pending",
    "profile_image_url": "https://cdn.example.com/profiles/user123.jpg"
  }
}
```

**Error Responses**:
- **401 Unauthorized**: Missing or invalid JWT token
- **404 Not Found**: User not found

---

### 2.4 Update Profile
**Endpoint**: `PUT /auth/profile/update` ⚠️ **Auth Required**

**Description**: Update user's full name or phone number

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `full_name` | String | No | Min 2, Max 100 chars | Updated full name |
| `phone` | String | No | 10-15 digits | Updated phone number |

**Example Request**:
```json
{
  "full_name": "Sita Shrestha Singh",
  "phone": "9851234567"
}
```

**Response (200 OK)**:
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha Singh",
    "email": "sita@example.com",
    "phone": "9851234567",
    "role": "user",
    "kyc_status": "pending"
  }
}
```

**Error Responses**:
- **400 Bad Request**: No fields to update or invalid values
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Profile update failed

---

### 2.5 Upload Profile Image
**Endpoint**: `POST /auth/profile/image` ⚠️ **Auth Required**

**Description**: Upload a profile image (max 5MB)

**Request Headers**:
```
Content-Type: multipart/form-data
Authorization: Bearer <JWT_TOKEN>
```

**Form Data Parameters**:
| Parameter | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `image` | File | Yes | Max 5MB (JPG, PNG) | Profile image file |

**Response (200 OK)**:
```json
{
  "message": "profile image uploaded successfully",
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha",
    "email": "sita@example.com",
    "profile_image_url": "https://cdn.example.com/profiles/user123_new.jpg"
  }
}
```

**Error Responses**:
- **400 Bad Request**: No image file or file too large
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Upload failed

---

## 3. Wallet Management (Two-Wallet System)

The system uses two wallet types:
- **Main Wallet**: For deposits and general fund management
- **Trading Wallet**: For active trading with locked balance for pending orders

---

### 3.1 Get All Wallets (Both Main & Trading)
**Endpoint**: `GET /wallet` ⚠️ **Auth Required**

**Description**: Get both main and trading wallet information for the current user

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "main_wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "15000.500000",
    "updated_at": "2026-03-09T11:45:00Z"
  },
  "trading_wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "5000.000000",
    "locked_balance": "1200.000000",
    "available_balance": "3800.000000",
    "updated_at": "2026-03-09T11:45:00Z"
  }
}
```

**Error Responses**:
- **401 Unauthorized**: Missing or invalid JWT token
- **500 Internal Server Error**: Database error

---

### 3.2 Get Main Wallet
**Endpoint**: `GET /wallet/main` ⚠️ **Auth Required**

**Description**: Get main wallet balance only (for deposits)

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "15000.500000",
    "updated_at": "2026-03-09T11:45:00Z"
  }
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Failed to fetch wallet

---

### 3.3 Get Trading Wallet
**Endpoint**: `GET /wallet/trading` ⚠️ **Auth Required**

**Description**: Get trading wallet with available and locked balance breakdown

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "5000.000000",
    "locked_balance": "1200.000000",
    "available_balance": "3800.000000",
    "updated_at": "2026-03-09T11:45:00Z"
  }
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Failed to fetch wallet

---

### 3.4 Top-Up Main Wallet (Deposit)
**Endpoint**: `POST /wallet/topup` ⚠️ **Auth Required**

**Description**: Add funds to main wallet (simulated deposit)

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `amount` | String | Yes | Positive decimal | Amount to deposit |

**Example Request**:
```json
{
  "amount": "1000.00"
}
```

**Response (200 OK)**:
```json
{
  "message": "top-up successful",
  "wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "16000.500000",
    "updated_at": "2026-03-09T11:50:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid amount (not numeric, negative, or zero)
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Top-up failed

---

### 3.5 Get Main Wallet Only
**Endpoint**: `GET /wallet/main` ⚠️ **Auth Required**

**Description**: Get main wallet balance

**Response**:
```json
{
  "wallet": {
    "user_id": "uuid",
    "balance": "15000.000000",
    "updated_at": "2026-03-09T11:45:00Z"
  }
}
```

---

### 3.6 Transfer Between Wallets
**Endpoint**: `POST /wallet/transfer` ⚠️ **Auth Required**

**Description**: Transfer funds between main and trading wallet

**Request Body Parameters**:
| Parameter | Type | Required | Values | Description |
|---|---|---|---|---|
| `amount` | String | Yes | Positive decimal | Amount to transfer |
| `direction` | String | Yes | `main_to_trading` or `trading_to_main` | Transfer direction |

**Example Request (Main to Trading)**:
```json
{
  "amount": "2000.00",
  "direction": "main_to_trading"
}
```

**Example Request (Trading to Main)**:
```json
{
  "amount": "500.00",
  "direction": "trading_to_main"
}
```

**Response (200 OK)**:
```json
{
  "message": "transfer successful",
  "main_wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "13000.500000"
  },
  "trading_wallet": {
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "balance": "7000.000000",
    "locked_balance": "1200.000000",
    "available_balance": "5800.000000"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid amount or invalid direction
- **401 Unauthorized**: Not authenticated
- **422 Unprocessable Entity**: Insufficient balance

---

### 3.7 Get Transfer History
**Endpoint**: `GET /wallet/transfers` ⚠️ **Auth Required**

**Description**: Get history of wallet transfers (up to 50 recent transfers)

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "transfers": [
    {
      "id": "transfer-uuid-1",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "amount": "2000.000000",
      "direction": "main_to_trading",
      "from_balance": "14000.500000",
      "to_balance": "5000.000000",
      "created_at": "2026-03-09T11:45:00Z"
    },
    {
      "id": "transfer-uuid-2",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "amount": "500.000000",
      "direction": "trading_to_main",
      "from_balance": "5000.000000",
      "to_balance": "14500.500000",
      "created_at": "2026-03-09T10:30:00Z"
    }
  ]
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Failed to fetch transfers

---

## 4. IPO Management

### 4.1 List All IPOs (Public)
**Endpoint**: `GET /ipos`

**Description**: Get list of all IPOs (active, closed, and allocated)

**Query Parameters**: None

**Response (200 OK)**:
```json
{
  "ipos": [
    {
      "id": "ipo-uuid-1",
      "company_id": "comp-uuid-1",
      "company_name": "Tech Bank Ltd",
      "symbol": "TECHBANK",
      "price_per_share": "100.000000",
      "total_shares": 10000,
      "shares_remaining": 5000,
      "max_per_applicant": 50,
      "status": "open",
      "total_applications": 250,
      "open_at": "2026-03-01T00:00:00Z",
      "close_at": "2026-03-31T23:59:59Z"
    },
    {
      "id": "ipo-uuid-2",
      "company_id": "comp-uuid-2",
      "company_name": "Finance Corp",
      "symbol": "FINCORP",
      "price_per_share": "150.000000",
      "total_shares": 5000,
      "shares_remaining": 0,
      "max_per_applicant": 25,
      "status": "allocated",
      "total_applications": 120,
      "open_at": "2026-02-01T00:00:00Z",
      "close_at": "2026-02-28T23:59:59Z"
    }
  ]
}
```

**Error Responses**:
- **500 Internal Server Error**: Failed to fetch IPOs

---

### 4.2 Get IPO Details
**Endpoint**: `GET /ipos/{id}` 

**Description**: Get detailed information about a specific IPO

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | String | Yes | IPO ID (UUID) |

**Response (200 OK)**:
```json
{
  "ipo": {
    "id": "ipo-uuid-1",
    "company_id": "comp-uuid-1",
    "company_name": "Tech Bank Ltd",
    "symbol": "TECHBANK",
    "sector": "Banking",
    "price_per_share": "100.000000",
    "total_shares": 10000,
    "shares_remaining": 5000,
    "max_per_applicant": 50,
    "status": "open",
    "total_applications": 250,
    "total_applied_shares": 8500,
    "open_at": "2026-03-01T00:00:00Z",
    "close_at": "2026-03-31T23:59:59Z",
    "created_at": "2026-02-28T10:00:00Z"
  }
}
```

**Error Responses**:
- **404 Not Found**: IPO not found with given ID
- **500 Internal Server Error**: Database error

---

### 4.3 Apply for IPO
**Endpoint**: `POST /ipos/{id}/apply` ⚠️ **Auth Required**

**Description**: Submit application to buy shares in an open IPO

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | String | Yes | IPO ID (UUID) |

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `shares_requested` | Integer | Yes | > 0, ≤ max_per_applicant | Number of shares to request |

**Example Request**:
```json
{
  "shares_requested": 50
}
```

**Response (201 Created)**:
```json
{
  "message": "IPO application submitted",
  "application": {
    "id": "app-uuid-1",
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "ipo_id": "ipo-uuid-1",
    "shares_requested": 50,
    "shares_allocated": null,
    "status": "pending",
    "total_cost": "5000.000000",
    "created_at": "2026-03-10T15:30:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid shares_requested or IPO not found
- **401 Unauthorized**: Not authenticated
- **422 Unprocessable Entity**: Already applied for this IPO

---

### 4.4 Create Company (Admin Only)
**Endpoint**: `POST /admin/companies` ⚠️ **Auth Required (Admin)**

**Description**: Register a new company in the system

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `symbol` | String | Yes | Unique, 2-8 chars | Stock symbol (e.g., NABIL) |
| `name` | String | Yes | 2-100 chars | Company name |
| `sector` | String | No | Max 50 chars | Industry sector (default: "General") |
| `total_supply` | Integer | Yes | > 0 | Total shares issued |

**Example Request**:
```json
{
  "symbol": "TECHBANK",
  "name": "Tech Bank Ltd",
  "sector": "Banking",
  "total_supply": 1000000
}
```

**Response (201 Created)**:
```json
{
  "message": "company created",
  "company": {
    "id": "comp-uuid-1",
    "symbol": "TECHBANK",
    "name": "Tech Bank Ltd",
    "sector": "Banking",
    "total_supply": 1000000,
    "created_at": "2026-03-10T15:30:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input or symbol already exists
- **403 Forbidden**: Admin access required

---

### 4.5 Launch IPO (Admin Only)
**Endpoint**: `POST /admin/ipos` ⚠️ **Auth Required (Admin)**

**Description**: Start an IPO for a specific company

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `company_id` | String | Yes | Valid company ID | Company UUID |
| `price_per_share` | String | Yes | Positive decimal | IPO share price |
| `total_shares` | Integer | Yes | > 0 | Shares offered in IPO |
| `max_per_applicant` | Integer | Yes | > 0 | Max shares per applicant |
| `open_at` | String | Yes | RFC3339 format | IPO opening time |
| `close_at` | String | Yes | RFC3339 format | IPO closing time |

**Example Request**:
```json
{
  "company_id": "comp-uuid-1",
  "price_per_share": "100.00",
  "total_shares": 10000,
  "max_per_applicant": 50,
  "open_at": "2026-03-01T00:00:00Z",
  "close_at": "2026-03-31T23:59:59Z"
}
```

**Response (201 Created)**:
```json
{
  "message": "IPO launched",
  "ipo": {
    "id": "ipo-uuid-1",
    "company_id": "comp-uuid-1",
    "company_name": "Tech Bank Ltd",
    "price_per_share": "100.000000",
    "total_shares": 10000,
    "shares_remaining": 10000,
    "max_per_applicant": 50,
    "status": "open",
    "open_at": "2026-03-01T00:00:00Z",
    "close_at": "2026-03-31T23:59:59Z",
    "created_at": "2026-03-10T15:30:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input or date format
- **403 Forbidden**: Admin access required

---

### 4.6 Allocate IPO (Admin Only - Lottery Based)
**Endpoint**: `POST /admin/ipos/{id}/allocate` ⚠️ **Auth Required (Admin)**

**Description**: Execute lottery-based IPO allocation after close. Refunds unsuccessful bids proportionally.

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | String | Yes | IPO ID (UUID) |

**Response (200 OK)**:
```json
{
  "message": "IPO allocation complete",
  "result": {
    "ipo_id": "ipo-uuid-1",
    "company_name": "Tech Bank Ltd",
    "total_shares_offered": 10000,
    "total_shares_applied": 15000,
    "total_allocated": 10000,
    "unsuccessful_shares": 5000,
    "successful_applicants": 180,
    "unsuccessful_applicants": 70,
    "total_refunded_amount": "500000.000000",
    "allocation_date": "2026-03-10T15:30:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: IPO not found or already allocated
- **403 Forbidden**: Admin access required

---

## 5. Trading & Orders

### 5.1 Get Order Book (Public)
**Endpoint**: `GET /orderbook/{company_id}`

**Description**: View current buy (bid) and sell (ask) orders for a company

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `company_id` | String | Yes | Company ID (UUID) |

**Response (200 OK)**:
```json
{
  "order_book": {
    "company_id": "comp-uuid-1",
    "company_symbol": "NABIL",
    "bids": [
      {
        "price": "105.50",
        "quantity": 1000,
        "orders": 5,
        "total_value": "105500.000000"
      },
      {
        "price": "105.00",
        "quantity": 500,
        "orders": 2,
        "total_value": "52500.000000"
      }
    ],
    "asks": [
      {
        "price": "106.00",
        "quantity": 200,
        "orders": 1,
        "total_value": "21200.000000"
      },
      {
        "price": "106.50",
        "quantity": 800,
        "orders": 3,
        "total_value": "85200.000000"
      }
    ],
    "last_traded_price": "105.75",
    "last_traded_quantity": 100,
    "last_traded_at": "2026-03-10T15:25:30Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid company_id format
- **500 Internal Server Error**: Database error

---

### 5.2 Place Buy Order
**Endpoint**: `POST /orders/buy` ⚠️ **Auth Required**

**Description**: Place a new buy order (limit or market type)

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `company_id` | String | Yes | Valid UUID | Company to buy |
| `quantity` | Integer | Yes | > 0 | Number of shares |
| `price` | String | Yes | Positive decimal | Price per share (for limit) |
| `order_type` | String | No | `limit` or `market` | Order type (default: `limit`) |

**Example Request (Limit Order)**:
```json
{
  "company_id": "comp-uuid-1",
  "quantity": 10,
  "price": "105.50",
  "order_type": "limit"
}
```

**Example Request (Market Order)**:
```json
{
  "company_id": "comp-uuid-1",
  "quantity": 10,
  "price": "0",
  "order_type": "market"
}
```

**Response (201 Created)**:
```json
{
  "message": "buy order placed",
  "order": {
    "id": "order-uuid-1",
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "company_id": "comp-uuid-1",
    "quantity": 10,
    "remaining_quantity": 3,
    "price": "105.500000",
    "total_cost": "1055.000000",
    "status": "open",
    "side": "buy",
    "order_type": "limit",
    "created_at": "2026-03-10T15:30:00Z"
  },
  "matches": 7
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input or insufficient balance
- **401 Unauthorized**: Not authenticated
- **422 Unprocessable Entity**: Insufficient trading wallet balance

---

### 5.3 Place Sell Order
**Endpoint**: `POST /orders/sell` ⚠️ **Auth Required**

**Description**: Place a new sell limit order

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `company_id` | String | Yes | Valid UUID | Company to sell |
| `quantity` | Integer | Yes | > 0 | Number of shares |
| `price` | String | Yes | Positive decimal | Price per share |

**Example Request**:
```json
{
  "company_id": "comp-uuid-1",
  "quantity": 5,
  "price": "110.00"
}
```

**Response (201 Created)**:
```json
{
  "message": "sell order placed",
  "order": {
    "id": "order-uuid-2",
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "company_id": "comp-uuid-1",
    "quantity": 5,
    "remaining_quantity": 5,
    "price": "110.000000",
    "total_value": "550.000000",
    "status": "open",
    "side": "sell",
    "created_at": "2026-03-10T15:30:00Z"
  },
  "matches": 0
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input or insufficient holdings
- **401 Unauthorized**: Not authenticated
- **422 Unprocessable Entity**: Insufficient shares to sell

---

### 5.4 Cancel Order
**Endpoint**: `DELETE /orders/{id}` ⚠️ **Auth Required**

**Description**: Cancel an open buy or sell order

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | String | Yes | Order ID (UUID) |

**Response (200 OK)**:
```json
{
  "message": "order cancelled"
}
```

**Error Responses**:
- **400 Bad Request**: Order not found or already executed
- **401 Unauthorized**: Not authenticated or order doesn't belong to user

---

### 5.5 Get User's Orders
**Endpoint**: `GET /orders/my` ⚠️ **Auth Required**

**Description**: Get history and current open orders for authenticated user (up to 50 recent)

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "orders": [
    {
      "id": "order-uuid-1",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "quantity": 10,
      "remaining_quantity": 3,
      "executed_quantity": 7,
      "price": "105.500000",
      "total_cost": "1055.000000",
      "status": "partially_executed",
      "side": "buy",
      "order_type": "limit",
      "created_at": "2026-03-10T15:30:00Z"
    },
    {
      "id": "order-uuid-2",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "quantity": 5,
      "remaining_quantity": 5,
      "executed_quantity": 0,
      "price": "110.000000",
      "total_value": "550.000000",
      "status": "open",
      "side": "sell",
      "created_at": "2026-03-10T15:30:00Z"
    }
  ]
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Database error

---

### 5.6 Get User's Portfolio (Holdings)
**Endpoint**: `GET /portfolio` ⚠️ **Auth Required**

**Description**: Get all current stock holdings with cost basis and current value

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "portfolio": [
    {
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "company_name": "Nabil Bank Ltd",
      "quantity": 50,
      "avg_buy_price": "102.500000",
      "current_price": "105.750000",
      "total_cost_basis": "5125.000000",
      "current_value": "5287.500000",
      "unrealized_gain_loss": "162.500000",
      "unrealized_gain_loss_percent": "3.17"
    },
    {
      "company_id": "comp-uuid-2",
      "company_symbol": "EVEREST",
      "company_name": "Everest Bank Ltd",
      "quantity": 25,
      "avg_buy_price": "500.000000",
      "current_price": "510.000000",
      "total_cost_basis": "12500.000000",
      "current_value": "12750.000000",
      "unrealized_gain_loss": "250.000000",
      "unrealized_gain_loss_percent": "2.00"
    }
  ],
  "portfolio_summary": {
    "total_holdings": 2,
    "total_shares": 75,
    "total_cost_basis": "17625.000000",
    "total_current_value": "18037.500000",
    "total_unrealized_gain_loss": "412.500000",
    "total_unrealized_gain_loss_percent": "2.34"
  }
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Failed to fetch portfolio

---

### 5.7 Get User's Trades (Executed Trades)
**Endpoint**: `GET /trades` ⚠️ **Auth Required**

**Description**: Get history of all executed trades for the user (up to 50 recent)

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "trades": [
    {
      "id": "trade-uuid-1",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "side": "buy",
      "quantity": 5,
      "price": "105.500000",
      "total_value": "527.500000",
      "commission": "2.637500",
      "net_value": "530.137500",
      "matched_with_order_id": "order-uuid-3",
      "executed_at": "2026-03-10T15:30:00Z"
    },
    {
      "id": "trade-uuid-2",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "side": "buy",
      "quantity": 2,
      "price": "105.100000",
      "total_value": "210.200000",
      "commission": "1.051000",
      "net_value": "211.251000",
      "matched_with_order_id": "order-uuid-4",
      "executed_at": "2026-03-10T15:28:30Z"
    },
    {
      "id": "trade-uuid-3",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-2",
      "company_symbol": "EVEREST",
      "side": "sell",
      "quantity": 10,
      "price": "510.000000",
      "total_value": "5100.000000",
      "commission": "25.500000",
      "net_value": "5074.500000",
      "matched_with_order_id": "order-uuid-5",
      "executed_at": "2026-03-09T14:20:00Z"
    }
  ]
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Database error

---

## 6. Market Data & Live Trading

### 6.1 List All Companies (Public)
**Endpoint**: `GET /market/companies`

**Description**: Get list of all trading companies with summary data

**Response (200 OK)**:
```json
{
  "companies": [
    {
      "symbol": "NABIL",
      "company_name": "Nabil Bank Ltd",
      "sector": "Banking",
      "current_price": "106.000000",
      "change_percent": "0.45",
      "volume": 25000,
      "turnover": "2650000.000000"
    },
    {
      "symbol": "EVEREST",
      "company_name": "Everest Bank Ltd",
      "sector": "Banking",
      "current_price": "510.000000",
      "change_percent": "-0.98",
      "volume": 12500,
      "turnover": "6375000.000000"
    }
  ],
  "count": 2
}
```

**Error Responses**:
- **500 Internal Server Error**: Failed to fetch companies

---

### 6.2 Get Company Details (Public)
**Endpoint**: `GET /market/companies/{symbol}`

**Description**: Get detailed information for a specific company by symbol

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `symbol` | String | Yes | Company stock symbol (e.g., NABIL, EVEREST) |

**Response (200 OK)**:
```json
{
  "company": {
    "symbol": "NABIL",
    "company_name": "Nabil Bank Ltd",
    "sector": "Banking",
    "current_price": "106.000000",
    "previous_close": "105.500000",
    "open_price": "105.500000",
    "high": "107.000000",
    "low": "104.500000",
    "change": "0.500000",
    "change_percent": "0.47",
    "volume": 25000,
    "turnover": "2650000.000000",
    "ltp_updated_at": "2026-03-10T15:35:00Z"
  }
}
```

**Error Responses**:
- **404 Not Found**: Company with symbol not found
- **500 Internal Server Error**: Database error

---

### 6.3 Get Live Trading Data
**Endpoint**: `GET /market/live`

**Description**: Get live price and volume data for all companies trading

**Response (200 OK)**:
```json
{
  "data": [
    {
      "symbol": "NABIL",
      "company_name": "Nabil Bank Ltd",
      "sector": "Banking",
      "ltp": "106.000000",
      "previous_close": "105.500000",
      "open_price": "105.500000",
      "high": "107.000000",
      "low": "104.500000",
      "change": "0.500000",
      "change_percent": "0.47",
      "volume": 25000,
      "turnover": "2650000.000000"
    },
    {
      "symbol": "EVEREST",
      "company_name": "Everest Bank Ltd",
      "sector": "Banking",
      "ltp": "510.000000",
      "previous_close": "515.000000",
      "open_price": "515.000000",
      "high": "516.000000",
      "low": "509.000000",
      "change": "-5.000000",
      "change_percent": "-0.97",
      "volume": 12500,
      "turnover": "6375000.000000"
    }
  ]
}
```

**Error Responses**:
- **500 Internal Server Error**: Failed to fetch market data

---

### 6.4 Get Market Index
**Endpoint**: `GET /market/index`

**Description**: Get overall market indicators (advances, declines, market cap)

**Response (200 OK)**:
```json
{
  "index": {
    "index_value": "2054.320000",
    "previous_index": "2042.100000",
    "change": "12.220000",
    "change_percent": "0.60",
    "advances": 140,
    "declines": 80,
    "unchanged": 12,
    "total_companies": 232,
    "total_market_cap": "35000000000.000000",
    "total_volume": 150000,
    "total_turnover": "18500000.000000",
    "last_updated": "2026-03-10T15:35:00Z"
  }
}
```

**Error Responses**:
- **500 Internal Server Error**: Failed to calculate market index

---

### 6.5 Get Candlestick Data (OHLCV)
**Endpoint**: `GET /market/companies/{symbol}/candles`

**Description**: Get OHLCV candlestick data for technical analysis

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `symbol` | String | Yes | Company stock symbol |

**Query Parameters**:
| Parameter | Type | Default | Description |
|---|---|---|---|
| `timeframe` | String | `1D` | Candle period (`1D`, `1W`, `1M`, etc.) |
| `days` | Integer | `30` | Number of days of data |

**Example URL**: `GET /market/companies/NABIL/candles?timeframe=1D&days=30`

**Response (200 OK)**:
```json
{
  "candles": [
    {
      "timestamp": "2026-02-09T00:00:00Z",
      "open": "104.000000",
      "high": "105.500000",
      "low": "103.500000",
      "close": "105.000000",
      "volume": 18500
    },
    {
      "timestamp": "2026-02-10T00:00:00Z",
      "open": "105.000000",
      "high": "106.500000",
      "low": "104.800000",
      "close": "105.500000",
      "volume": 22000
    },
    {
      "timestamp": "2026-03-10T00:00:00Z",
      "open": "105.500000",
      "high": "107.000000",
      "low": "104.500000",
      "close": "106.000000",
      "volume": 25000
    }
  ]
}
```

**Error Responses**:
- **400 Bad Request**: Invalid symbol or timeframe
- **500 Internal Server Error**: Failed to fetch candlestick data

---

### 6.6 Stream Live Prices (SSE - Server-Sent Events)
**Endpoint**: `GET /market/stream`

**Description**: Establish real-time WebSocket-style connection for live price updates using Server-Sent Events

**Request Headers**:
```
Accept: text/event-stream
```

**Response Stream (200 OK - Streaming)**:
```
data: {"symbol":"NABIL","price":"106.000000","volume":25000,"timestamp":"2026-03-10T15:35:00Z","change":"0.47"}

data: {"symbol":"EVEREST","price":"510.000000","volume":12500,"timestamp":"2026-03-10T15:35:05Z","change":"-0.97"}

data: {"symbol":"NABIL","price":"106.050000","volume":25150,"timestamp":"2026-03-10T15:35:10Z","change":"0.52"}
```

**Connection Notes**:
- Keep connection alive to continuously receive updates
- Reconnect if disconnected
- Updates fire whenever price changes occur

**Error Responses**:
- **500 Internal Server Error**: Streaming service unavailable

---

## 7. Price Triggers (Auto-Trading)

### 7.1 Create Price Trigger
**Endpoint**: `POST /triggers` ⚠️ **Auth Required**

**Description**: Create an automated sell trigger that executes when stock price reaches a target

**Request Body Parameters**:
| Parameter | Type | Required | Validation | Description |
|---|---|---|---|---|
| `company_id` | String | Yes | Valid UUID | Company to watch |
| `trigger_price` | String | Yes | Positive decimal | Target price to trigger |
| `shares_qty` | Integer | Yes | > 0 | Shares to auto-sell |
| `direction` | String | Yes | `above` or `below` | Trigger when price goes above or below |

**Example Request (Auto-sell above 150)**:
```json
{
  "company_id": "comp-uuid-1",
  "trigger_price": "150.00",
  "shares_qty": 10,
  "direction": "above"
}
```

**Example Request (Auto-sell below 50)**:
```json
{
  "company_id": "comp-uuid-2",
  "trigger_price": "50.00",
  "shares_qty": 5,
  "direction": "below"
}
```

**Response (201 Created)**:
```json
{
  "message": "price trigger created",
  "trigger": {
    "id": "trigger-uuid-1",
    "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "company_id": "comp-uuid-1",
    "company_symbol": "NABIL",
    "trigger_price": "150.000000",
    "shares_qty": 10,
    "direction": "above",
    "status": "active",
    "created_at": "2026-03-10T15:30:00Z"
  }
}
```

**How It Works**:
- Trigger remains active until executed or cancelled
- When LTP crosses the trigger price, a sell order is automatically placed
- After execution, trigger status changes to `executed`
- Shares up to `shares_qty` are sold at market price

**Error Responses**:
- **400 Bad Request**: Invalid input or company not found
- **401 Unauthorized**: Not authenticated

---

### 7.2 Cancel Price Trigger
**Endpoint**: `DELETE /triggers/{id}` ⚠️ **Auth Required**

**Description**: Cancel an active price trigger

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | String | Yes | Trigger ID (UUID) |

**Response (200 OK)**:
```json
{
  "message": "trigger cancelled"
}
```

**Error Responses**:
- **400 Bad Request**: Trigger not found or already executed
- **401 Unauthorized**: Not authenticated or trigger doesn't belong to user

---

### 7.3 Get User's Triggers
**Endpoint**: `GET /triggers` ⚠️ **Auth Required**

**Description**: Get all active and historical price triggers for the user

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200 OK)**:
```json
{
  "triggers": [
    {
      "id": "trigger-uuid-1",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-1",
      "company_symbol": "NABIL",
      "trigger_price": "150.000000",
      "shares_qty": 10,
      "direction": "above",
      "status": "active",
      "created_at": "2026-03-10T15:30:00Z"
    },
    {
      "id": "trigger-uuid-2",
      "user_id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
      "company_id": "comp-uuid-2",
      "company_symbol": "EVEREST",
      "trigger_price": "50.000000",
      "shares_qty": 5,
      "direction": "below",
      "status": "executed",
      "execution_price": "49.500000",
      "execution_time": "2026-03-09T14:15:00Z",
      "created_at": "2026-03-08T10:00:00Z"
    }
  ]
}
```

**Error Responses**:
- **401 Unauthorized**: Not authenticated
- **500 Internal Server Error**: Failed to fetch triggers

---

## 8. Admin Operations

### 8.1 Update User KYC and Role (Admin Only)
**Endpoint**: `PUT /admin/users/{user_id}/kyc` ⚠️ **Auth Required (Admin)**

**Description**: Update KYC verification status and user role

**Path Parameters**:
| Parameter | Type | Required | Description |
|---|---|---|---|
| `user_id` | String | Yes | User ID (UUID) |

**Request Body Parameters**:
| Parameter | Type | Required | Values | Description |
|---|---|---|---|---|
| `kyc_status` | String | Yes | `pending`, `verified`, `rejected`, `under_review` | KYC status |
| `role` | String | Yes | `user` or `admin` | User role |

**Example Request**:
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

**Response (200 OK)**:
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": "7f8b9e10-a1b2-c3d4-e5f6-7g8h9i0j1k2l",
    "full_name": "Sita Shrestha",
    "email": "sita@example.com",
    "role": "user",
    "kyc_status": "verified",
    "updated_at": "2026-03-10T15:30:00Z"
  }
}
```

**Error Responses**:
- **400 Bad Request**: Invalid input or user not found
- **403 Forbidden**: Admin access required

---

## 9. Error Codes & Status Reference

### HTTP Status Codes

| Status | Meaning | Common Causes |
|---|---|---|
| **200** | OK | Request successful |
| **201** | Created | Resource created successfully |
| **400** | Bad Request | Invalid input, missing fields, malformed JSON |
| **401** | Unauthorized | Missing JWT token or token expired |
| **403** | Forbidden | Admin access required or insufficient permissions |
| **404** | Not Found | Resource not found (user, order, company, etc.) |
| **409** | Conflict | Duplicate resource (email already registered) |
| **422** | Unprocessable Entity | Insufficient balance, insufficient holdings, business logic error |
| **500** | Internal Server Error | Database error or server crash |

---

### Common Error Messages

| Error | HTTP Status | Solution |
|---|---|---|
| `unauthorized` | 401 | JWT token is missing, invalid, or expired. Login again and use new token. |
| `invalid credentials` | 401 | Email or password is incorrect. |
| `user not found` | 404 | User does not exist. Check user ID. |
| `user already exists` | 409 | Email is already registered. Use different email. |
| `invalid email format` | 400 | Email format is invalid. Check spelling. |
| `password too short` | 400 | Password must be at least 6 characters. |
| `invalid amount` | 400 | Amount must be numeric and positive. |
| `insufficient balance` | 422 | Wallet balance is insufficient. Top up or transfer to wallet. |
| `insufficient available` | 422 | Available balance (excluding locked) is insufficient. |
| `insufficient holdings` | 422 | You don't have enough shares to sell. |
| `company not found` | 404 | Company symbol or ID doesn't exist. |
| `order not found` | 404 | Order ID doesn't exist or belongs to different user. |
| `IPO not found` | 404 | IPO ID doesn't exist. Check ID. |
| `admin access required` | 403 | Only admin users can perform this action. |
| `IPO already applied` | 422 | You have already applied for this IPO. |
| `IPO not open` | 422 | IPO is not currently open for applications. |
| `max shares exceeded` | 422 | Requested shares exceed per-applicant limit. |

---

## 10. Data Types & Format Reference

### Date/Time Format
All timestamps use **RFC3339** format with UTC timezone:
```
2026-03-10T15:30:00Z
2026-03-10T15:30:45.123456Z
```

### Decimal Numbers
All currency and price values are strings with decimals (to avoid precision loss):
```json
{
  "price": "105.500000",
  "balance": "15000.123456",
  "amount": "1000.00"
}
```

When sending decimals in requests, you can use either format:
- `"1000"` or `"1000.00"` or `"1000.123456"`

### Pagination
Current system returns recent records (up to 50 by default):
- `/orders/my` - Last 50 orders
- `/trades` - Last 50 trades
- `/wallet/transfers` - Last 50 transfers

---

## 11. Authentication Flow

### Step 1: Register
```bash
POST /auth/register
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "9841234567",
  "password": "password123",
  "role": "user"
}

Response: 201 Created (user created, not logged in)
```

### Step 2: Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...}
}
```

### Step 3: Use Token
```bash
GET /wallet
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Response: 200 OK
{
  "main_wallet": {...},
  "trading_wallet": {...}
}
```

---

## 12. Complete Trading Workflow Example

### Scenario: Buy and Sell Stocks

**1. Check Wallets:**
```bash
GET /wallet
Authorization: Bearer <TOKEN>
```

**2. Top Up Main Wallet (if needed):**
```bash
POST /wallet/topup
{ "amount": "5000.00" }
```

**3. Transfer to Trading Wallet:**
```bash
POST /wallet/transfer
{
  "amount": "2000.00",
  "direction": "main_to_trading"
}
```

**4. View Order Book:**
```bash
GET /orderbook/{company_id}
```

**5. Place Buy Order:**
```bash
POST /orders/buy
{
  "company_id": "comp-uuid",
  "quantity": 10,
  "price": "106.00",
  "order_type": "limit"
}
```

**6. Check Portfolio:**
```bash
GET /portfolio
Authorization: Bearer <TOKEN>
```

**7. Place Sell Order:**
```bash
POST /orders/sell
{
  "company_id": "comp-uuid",
  "quantity": 5,
  "price": "115.00"
}
```

**8. View Trade History:**
```bash
GET /trades
Authorization: Bearer <TOKEN>
```

---

## 13. CORS Policy

The API allows requests from:
- `http://localhost:3000` (development frontend)
- `*` (wildcard - all origins in development mode)

**Allowed Methods**:
- GET, POST, PUT, PATCH, DELETE, OPTIONS

**Allowed Headers**:
- Origin, Content-Type, Accept, Authorization

---

## 14. Rate Limiting

Currently, there is **no rate limiting** in the development version. For production, implement:
- 100 requests per minute per IP/user
- 10 requests per second for market data endpoints
- 1 request per second for order placement

---

## 15. Webhook Support (Future)

Planned webhook events:
- Order execution
- Price trigger fired
- IPO allocation complete
- Portfolio value alert

---

## Appendix: Testing with cURL

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9841234567",
    "password": "password123",
    "role": "user"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Wallets (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### List Companies
```bash
curl -X GET http://localhost:8080/api/v1/market/companies
```

### Place Buy Order
```bash
curl -X POST http://localhost:8080/api/v1/orders/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "company_id": "comp-uuid",
    "quantity": 10,
    "price": "106.00",
    "order_type": "limit"
  }'
```

---

## Support & Contacts

- **API Docs**: http://localhost:8080/swagger/index.html
- **Issue Tracker**: [GitHub Issues]
- **Email**: support@stockmarketsimu.dev

*Last Updated: March 9, 2026*
*Version: 2.0*

