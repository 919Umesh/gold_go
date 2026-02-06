# 🚀 Gold Go API Documentation (Appwrite Backend)

This document provides a comprehensive guide to the API endpoints available in the Gold Go project, refactored for use with Appwrite.

**Base URL**: `http://localhost:8080/api/v1`

---

## 🔐 Authentication

All protected routes require an `Authorization` header with a valid JWT Bearer token.

**Header Format**: `Authorization: Bearer <YOUR_JWT_TOKEN>`

### 1. Register User
- **Method**: `POST`
- **Endpoint**: `/auth/register`
- **Request Body**:
  ```json
  {
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9800000000",
    "password": "securepassword",
    "role": "user"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": "6985...",
    "full_name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "kyc_status": "pending"
  }
  ```

### 2. Login
- **Method**: `POST`
- **Endpoint**: `/auth/login`
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "user": {
      "id": "6985...",
      "full_name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    },
    "token": "eyJhbGciOiJIUzI1..."
  }
  ```

### 3. Get User Profile
- **Method**: `GET`
- **Endpoint**: `/auth/profile`
- **Protection**: 🔓 Protected
- **Response (200 OK)**:
  ```json
  {
    "id": "6985...",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "9800000000",
    "role": "user",
    "kyc_status": "verified"
  }
  ```

---

## 📈 Stock Market (Public)

These endpoints do not require authentication.

### 4. List All Companies
- **Method**: `GET`
- **Endpoint**: `/stocks`
- **Query Params**: `limit=10&offset=0`
- **Response**: List of company objects (Symbol, Name, Sector, MarketCap, CurrentPrice).

### 5. Get Company Detail
- **Method**: `GET`
- **Endpoint**: `/stocks/:symbol` (e.g., `/stocks/NABIL`)

### 6. Get Price History
- **Method**: `GET`
- **Endpoint**: `/stocks/:symbol/history`
- **Query Params**: `timeframe=1d` (Available: `1h`, `1d`, `1w`)

### 7. Market Overview
- **Method**: `GET`
- **Endpoint**: `/stocks/market-overview`
- **Description**: Returns top gainers, top losers, and most active stocks.

---

## 💰 Trading & Wallet (Protected)

Authorization required for all endpoints below.

### 8. Get Virtual Wallet
- **Method**: `GET`
- **Endpoint**: `/trading/wallet`
- **Response**:
  ```json
  {
    "balance": 50000.00,
    "total_invested": 12000.00,
    "total_profit_loss": 450.50
  }
  ```

### 9. Top Up Wallet
- **Method**: `POST`
- **Endpoint**: `/wallet/topup`
- **Body**: `{"amount": 10000.0}`

### 10. Buy Stock
- **Method**: `POST`
- **Endpoint**: `/trading/buy`
- **Request Body**:
  ```json
  {
    "symbol": "NICA",
    "quantity": 5
  }
  ```
- **Response**: `TradeResult` object with success status and transaction details.

### 11. Sell Stock
- **Method**: `POST`
- **Endpoint**: `/trading/sell`
- **Request Body**:
  ```json
  {
    "symbol": "NICA",
    "quantity": 2
  }
  ```

### 12. Get Portfolio
- **Method**: `GET`
- **Endpoint**: `/trading/portfolio`
- **Description**: Returns all stocks owned by the user with current valuation.

---

## 🛠 Admin & Setup

### 13. Seed Stock Data (Run ONCE)
- **Method**: `POST`
- **Endpoint**: `/admin/seed-stocks`
- **Protection**: 🔐 Admin Only
- **Description**: This endpoint populates your Appwrite database with initial companies, price history, and market events for testing.

---

## 📱 Mobile App Integration Guide

1.  **Base URL**: Use your local IP or a tunnel (like Ngrok) if testing on a physical device.
2.  **Token Storage**: Store the JWT token securely (e.g., `flutter_secure_storage`).
3.  **Appwrite Connection**: The Go API handles all Appwrite communication on the backend. You do **not** need to use the Appwrite SDK in your mobile app unless you want to use Appwrite's storage or Push notifications directly.
4.  **Error Handling**:
    - `401 Unauthorized`: Token missing or expired.
    - `403 Forbidden`: User needs Admin role or Verified KYC.
    - `429 Too Many Requests`: Rate limit reached.
