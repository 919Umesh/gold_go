# Gold Go API Documentation

Base URL: `https://<your-render-url>/api/v1`

## Authentication

### Register
- **URL**: `/auth/register`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword",
    "full_name": "John Doe"
  }
  ```

### Login
- **URL**: `/auth/login`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
- **Response**: Returns a JWT token.

## Gold

### Get Current Price
- **URL**: `/gold/price`
- **Method**: `GET`
- **Response**: Current gold price per gram in INR.

### Get Price History
- **URL**: `/gold/history?days=30`
- **Method**: `GET`
- **Response**: List of historical prices.

## Wallet (Protected)
**Headers**: `Authorization: Bearer <token>`

### Get Wallet
- **URL**: `/wallet`
- **Method**: `GET`

### Top Up
- **URL**: `/wallet/topup`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "amount": 1000.00,
    "reference_id": "unique-ref-123"
  }
  ```

### Buy Gold
- **URL**: `/wallet/buy`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "grams": 2.5,
    "price_per_gram": 6500.00,
    "reference_id": "buy-ref-456"
  }
  ```

### Sell Gold
- **URL**: `/wallet/sell`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "grams": 1.0,
    "price_per_gram": 6600.00,
    "reference_id": "sell-ref-789"
  }
  ```

### Get Transactions
- **URL**: `/transaction`
- **Method**: `GET`

## User (Protected)
**Headers**: `Authorization: Bearer <token>`

### Get Profile
- **URL**: `/auth/profile`
- **Method**: `GET`

### Update Profile
- **URL**: `/auth/profile/update`
- **Method**: `PUT`
