# Stock Market Simulator — Complete API Documentation

**Base URL:** `http://localhost:8080/api/v1`

All protected endpoints require `Authorization: Bearer <JWT_TOKEN>` header.

---

## Table of Contents

1. [Health Check](#1-health-check)
2. [Authentication](#2-authentication)
3. [Stocks & Market Data](#3-stocks--market-data)
4. [Live Trading & Real-Time](#4-live-trading--real-time)
5. [Sectors](#5-sectors)
6. [Predictions](#6-predictions)
7. [Wallet (Fiat)](#7-wallet-fiat)
8. [Trading](#8-trading)
9. [Admin](#9-admin)

---

## 1. Health Check

### GET `/health`

Check if the server is running.

**Response** `200 OK`
```json
{
  "status": "ok"
}
```

---

## 2. Authentication

### POST `/api/v1/auth/register`

Register a new user account.

**Request Body**
```json
{
  "full_name": "Ram Sharma",
  "email": "ram@example.com",
  "phone": "9841234567",
  "role": "user",
  "password": "securepass123"
}
```

| Field     | Type   | Required | Validation           |
|-----------|--------|----------|----------------------|
| full_name | string | Yes      | min=2, max=100       |
| email     | string | Yes      | valid email          |
| phone     | string | Yes      | min=10, max=15       |
| role      | string | Yes      | min=3, max=10        |
| password  | string | Yes      | min=6                |

**Response** `201 Created`
```json
{
  "message": "user registered successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Ram Sharma",
    "email": "ram@example.com",
    "phone": "9841234567",
    "kyc_status": "pending",
    "role": "user",
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:00:00Z"
  }
}
```

**Error** `409 Conflict`
```json
{
  "error": "user already exists"
}
```

---

### POST `/api/v1/auth/login`

Login and receive a JWT token.

**Request Body**
```json
{
  "email": "ram@example.com",
  "password": "securepass123"
}
```

**Response** `200 OK`
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Ram Sharma",
    "email": "ram@example.com",
    "phone": "9841234567",
    "kyc_status": "pending",
    "role": "user",
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:00:00Z"
  }
}
```

**Error** `401 Unauthorized`
```json
{
  "error": "invalid credentials"
}
```

---

### GET `/api/v1/auth/profile` 🔒

Get the authenticated user's profile.

**Headers:** `Authorization: Bearer <token>`

**Response** `200 OK`
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Ram Sharma",
    "email": "ram@example.com",
    "phone": "9841234567",
    "kyc_status": "pending",
    "role": "user",
    "profile_image_id": "",
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:00:00Z"
  }
}
```

---

### PUT `/api/v1/auth/profile/update` 🔒

Update user profile fields.

**Headers:** `Authorization: Bearer <token>`

**Request Body**
```json
{
  "full_name": "Ram Kumar Sharma",
  "phone": "9841234568"
}
```

| Field     | Type   | Required | Validation           |
|-----------|--------|----------|----------------------|
| full_name | string | No       | min=2, max=100       |
| phone     | string | No       | min=10, max=15       |

**Response** `200 OK`
```json
{
  "message": "profile updated successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Ram Kumar Sharma",
    "email": "ram@example.com",
    "phone": "9841234568",
    "kyc_status": "pending",
    "role": "user",
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:05:00Z"
  }
}
```

---

### POST `/api/v1/auth/profile/image` 🔒

Upload a profile image.

**Headers:** `Authorization: Bearer <token>`

**Request:** `multipart/form-data`

| Field | Type | Required | Validation    |
|-------|------|----------|---------------|
| image | file | Yes      | max 5MB       |

**Response** `200 OK`
```json
{
  "message": "profile image uploaded successfully",
  "image_id": "img_abc123"
}
```

---

## 3. Stocks & Market Data

### GET `/api/v1/stocks`

List all active companies with pagination.

**Query Parameters**

| Param  | Type | Default | Max |
|--------|------|---------|-----|
| limit  | int  | 50      | 100 |
| offset | int  | 0       | -   |

**Response** `200 OK`
```json
{
  "companies": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "description": "Leading private sector bank in Nepal",
      "founded_year": 1984,
      "employees": 3500,
      "total_shares": 1800000000,
      "available_shares": 1799999000,
      "is_active": true,
      "created_at": "2026-02-23T10:00:00Z",
      "updated_at": "2026-02-23T10:00:00Z"
    }
  ],
  "limit": 50,
  "offset": 0,
  "count": 25,
  "total": 25
}
```

---

### GET `/api/v1/stocks/search`

Search companies by symbol or name.

**Query Parameters**

| Param | Type   | Required |
|-------|--------|----------|
| q     | string | Yes      |

**Response** `200 OK`
```json
{
  "companies": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "total_shares": 1800000000,
      "available_shares": 1799999000,
      "is_active": true
    }
  ]
}
```

---

### GET `/api/v1/stocks/:symbol`

Get a single company by symbol.

**URL Parameters:** `symbol` — Company ticker (e.g., `NABIL`)

**Response** `200 OK`
```json
{
  "company": {
    "id": "uuid-1",
    "symbol": "NABIL",
    "name": "Nabil Bank Limited",
    "sector": "Banking",
    "market_cap": 180000000000,
    "description": "Leading private sector bank in Nepal",
    "founded_year": 1984,
    "employees": 3500,
    "total_shares": 1800000000,
    "available_shares": 1799999000,
    "is_active": true,
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:00:00Z"
  }
}
```

---

### GET `/api/v1/stocks/:symbol/price`

Get the current (latest) price for a stock.

**URL Parameters:** `symbol` — Company ticker

**Response** `200 OK`
```json
{
  "price": {
    "id": "uuid-price",
    "company_id": "uuid-1",
    "open_price": "100.00",
    "high_price": "102.50",
    "low_price": "99.50",
    "close_price": "101.25",
    "volume": 15000,
    "timestamp": "2026-02-23T10:30:00Z",
    "timeframe": "1m"
  }
}
```

---

### GET `/api/v1/stocks/:symbol/history`

Get historical price data for charting.

**URL Parameters:** `symbol` — Company ticker

**Query Parameters**

| Param     | Type   | Default | Allowed Values       |
|-----------|--------|---------|----------------------|
| timeframe | string | 1D      | 1D, 1W, 1M, all     |
| days      | int    | 30      | 1–1825 (5 years)     |

**Response** `200 OK`
```json
{
  "symbol": "NABIL",
  "timeframe": "1D",
  "days": 30,
  "count": 22,
  "prices": [
    {
      "id": "uuid-price",
      "company_id": "uuid-1",
      "open_price": "100.00",
      "high_price": "102.50",
      "low_price": "99.50",
      "close_price": "101.25",
      "volume": 15000,
      "timestamp": "2026-02-23T00:00:00Z",
      "timeframe": "1D"
    }
  ]
}
```

---

### GET `/api/v1/stocks/:symbol/events`

Get upcoming market events for a company.

**URL Parameters:** `symbol` — Company ticker

**Response** `200 OK`
```json
{
  "events": [
    {
      "id": "uuid-event",
      "company_id": "uuid-1",
      "event_type": "earnings",
      "title": "Q1 Financial Results",
      "description": "Strong quarterly earnings with revenue growth",
      "impact_percentage": 3.5,
      "event_date": "2026-06-15T00:00:00Z"
    }
  ]
}
```

---

### GET `/api/v1/stocks/:symbol/candles`

Get candlestick (OHLCV) data for charting. Returns data sorted ascending by timestamp.

**URL Parameters:** `symbol` — Company ticker

**Query Parameters**

| Param     | Type   | Default | Max |
|-----------|--------|---------|-----|
| timeframe | string | 1D      | -   |
| days      | int    | 30      | 365 |

**Response** `200 OK`
```json
{
  "symbol": "NABIL",
  "timeframe": "1D",
  "days": 30,
  "count": 22,
  "candles": [
    {
      "timestamp": "2026-02-01T00:00:00Z",
      "open": "100.00",
      "high": "102.50",
      "low": "99.50",
      "close": "101.25",
      "volume": 15000
    },
    {
      "timestamp": "2026-02-02T00:00:00Z",
      "open": "101.25",
      "high": "103.00",
      "low": "100.80",
      "close": "102.10",
      "volume": 12000
    }
  ]
}
```

---

### GET `/api/v1/stocks/market-overview`

Get a basic market overview with top gainers, losers, and most active.

**Response** `200 OK`
```json
{
  "total_companies": 25,
  "top_gainers": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "current_price": "105.20",
      "previous_price": "100.00",
      "change": "5.20",
      "change_percent": "5.20"
    }
  ],
  "top_losers": [],
  "most_active": []
}
```

---

### GET `/api/v1/stocks/top-gainers`

Get stocks with the highest positive price change.

**Query Parameters**

| Param | Type | Default |
|-------|------|---------|
| limit | int  | 10      |

**Response** `200 OK`
```json
{
  "gainers": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "current_price": "105.20",
      "previous_price": "100.00",
      "change": "5.20",
      "change_percent": "5.20"
    }
  ]
}
```

---

### GET `/api/v1/stocks/top-losers`

Get stocks with the largest negative price change.

**Query Parameters**

| Param | Type | Default |
|-------|------|---------|
| limit | int  | 10      |

**Response** `200 OK`
```json
{
  "losers": [
    {
      "id": "uuid-2",
      "symbol": "IMS",
      "name": "IMS Software Solutions",
      "sector": "Information Technology",
      "market_cap": 30000000000,
      "current_price": "97.50",
      "previous_price": "100.00",
      "change": "-2.50",
      "change_percent": "-2.50"
    }
  ]
}
```

---

### GET `/api/v1/stocks/most-active`

Get stocks with the highest trading volume.

**Query Parameters**

| Param | Type | Default |
|-------|------|---------|
| limit | int  | 10      |

**Response** `200 OK`
```json
{
  "active": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "total_shares": 1800000000,
      "available_shares": 1799999000,
      "is_active": true
    }
  ]
}
```

---

## 4. Live Trading & Real-Time

### GET `/api/v1/stocks/live-trading`

Get live trading board data for all companies (like NEPSE live board).

**Response** `200 OK`
```json
{
  "live_trading": [
    {
      "symbol": "NABIL",
      "company_id": "uuid-1",
      "company_name": "Nabil Bank Limited",
      "sector": "Banking",
      "ltp": "101.25",
      "change_percent": "1.25",
      "open": "100.00",
      "high": "102.50",
      "low": "99.50",
      "volume": 15000,
      "previous_close": "100.00",
      "difference": "1.25",
      "turnover": "1518750.00",
      "last_updated": "2026-02-23T11:30:00Z"
    },
    {
      "symbol": "GBIME",
      "company_id": "uuid-2",
      "company_name": "Global IME Bank Limited",
      "sector": "Banking",
      "ltp": "98.50",
      "change_percent": "-1.50",
      "open": "100.00",
      "high": "100.50",
      "low": "98.00",
      "volume": 8000,
      "previous_close": "100.00",
      "difference": "-1.50",
      "turnover": "788000.00",
      "last_updated": "2026-02-23T11:25:00Z"
    }
  ],
  "count": 25
}
```

---

### GET `/api/v1/stocks/market-index`

Get the overall market index value (like NEPSE index).

**Response** `200 OK`
```json
{
  "market_index": {
    "index_value": "2250.50",
    "change": "15.30",
    "change_percent": "0.68",
    "total_turnover": "45000000.00",
    "total_volume": 350000,
    "total_market_cap": "2250500000000.00",
    "advances": 15,
    "declines": 8,
    "unchanged": 2,
    "total_companies": 25,
    "previous_close": "2235.20",
    "timestamp": "2026-02-23T11:30:00Z"
  }
}
```

---

### GET `/api/v1/stocks/market-summary`

Get a comprehensive market overview including index, top gainers, top losers, most active, and sector breakdown.

**Response** `200 OK`
```json
{
  "index": {
    "index_value": "2250.50",
    "change": "15.30",
    "change_percent": "0.68",
    "total_turnover": "45000000.00",
    "total_volume": 350000,
    "total_market_cap": "2250500000000.00",
    "advances": 15,
    "declines": 8,
    "unchanged": 2,
    "total_companies": 25,
    "previous_close": "2235.20",
    "timestamp": "2026-02-23T11:30:00Z"
  },
  "top_gainers": [
    {
      "symbol": "NABIL",
      "company_name": "Nabil Bank Limited",
      "ltp": "105.20",
      "change_percent": "5.20",
      "volume": 25000,
      "turnover": "2630000.00"
    }
  ],
  "top_losers": [
    {
      "symbol": "IMS",
      "company_name": "IMS Software Solutions",
      "ltp": "97.50",
      "change_percent": "-2.50",
      "volume": 5000,
      "turnover": "487500.00"
    }
  ],
  "most_active": [
    {
      "symbol": "NABIL",
      "company_name": "Nabil Bank Limited",
      "ltp": "105.20",
      "volume": 25000,
      "turnover": "2630000.00"
    }
  ],
  "sector_summary": [
    {
      "sector": "Banking",
      "change": "3.50",
      "change_percent": "0.70",
      "turnover": "15000000.00",
      "volume": 120000,
      "company_count": 5,
      "advances": 3,
      "declines": 2
    },
    {
      "sector": "Information Technology",
      "change": "-1.20",
      "change_percent": "-0.24",
      "turnover": "8000000.00",
      "volume": 60000,
      "company_count": 5,
      "advances": 2,
      "declines": 3
    }
  ],
  "as_of": "2026-02-23T11:30:00Z"
}
```

---

### GET `/api/v1/stocks/stream`

Server-Sent Events (SSE) stream for real-time price updates. Connect and receive live events as trades happen.

**Headers:** `Accept: text/event-stream`

**Response:** `200 OK` (streaming `text/event-stream`)

Each event is a JSON object:

**Trade Event:**
``` 
event: message
data: {"type":"trade","data":{"symbol":"NABIL","company_name":"Nabil Bank Limited","trade_type":"buy","quantity":100,"price":"101.25","total_amount":"10125.00","price_impact":"0.35","new_price":"101.60","timestamp":"2026-02-23T11:31:00Z"}}
```

**Price Update Event:**
```
event: message
data: {"type":"price_update","data":{"symbol":"NABIL","company_id":"uuid-1","company_name":"Nabil Bank Limited","ltp":"101.60","change_percent":"1.60","open":"100.00","high":"102.50","low":"99.50","volume":15100,"previous_close":"100.00","difference":"1.60","turnover":"1534160.00","last_updated":"2026-02-23T11:31:00Z"}}
```

---

## 5. Sectors

### GET `/api/v1/sectors`

Get all available sectors.

**Response** `200 OK`
```json
{
  "message": "sectors retrieved successfully",
  "sectors": [
    "Banking",
    "Hydropower",
    "Information Technology",
    "Insurance",
    "Manufacturing",
    "Pharma",
    "Real Estate"
  ],
  "count": 7
}
```

---

### GET `/api/v1/sectors/:sector/companies`

Get companies in a specific sector.

**URL Parameters:** `sector` — Sector name (e.g., `Banking`)

**Query Parameters**

| Param  | Type | Default | Max |
|--------|------|---------|-----|
| limit  | int  | 50      | 100 |
| offset | int  | 0       | -   |

**Response** `200 OK`
```json
{
  "sector": "Banking",
  "companies": [
    {
      "id": "uuid-1",
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "sector": "Banking",
      "market_cap": 180000000000,
      "total_shares": 1800000000,
      "available_shares": 1799999000,
      "is_active": true
    }
  ],
  "limit": 50,
  "offset": 0,
  "count": 5,
  "total": 5
}
```

---

### GET `/api/v1/sectors/:sector/stats`

Get statistics for a specific sector.

**URL Parameters:** `sector` — Sector name

**Response** `200 OK`
```json
{
  "message": "sector statistics",
  "sector": "Banking",
  "statistics": {
    "company_count": 5,
    "total_market_cap": "760000000000",
    "avg_market_cap": "152000000000",
    "total_employees": 17700,
    "avg_employees": 3540,
    "avg_founded_year": 1993
  },
  "top_5_companies": [
    {
      "symbol": "NABIL",
      "name": "Nabil Bank Limited",
      "market_cap": 180000000000
    }
  ]
}
```

---

## 6. Predictions (ML Algorithms)

This module implements **7 classical ML/statistical algorithms** for stock price prediction:

| # | Algorithm | Type | Description |
|---|-----------|------|-------------|
| 1 | **Linear Regression** | Statistical | Fits a straight line using least squares |
| 2 | **Exponential Moving Average (EMA)** | Technical Analysis | Exponentially weighted recent prices |
| 3 | **Holt's Linear Trend** | Time Series | Double Exponential Smoothing (level + trend) |
| 4 | **K-Nearest Neighbors (KNN)** | Machine Learning | Pattern-matching on similar price windows |
| 5 | **Auto-Regressive AR(3)** | Time Series | Foundation of ARIMA; uses past 3 values |
| 6 | **Weighted Moving Average (WMA)** | Technical Analysis | Linearly weighted moving average |
| 7 | **Ensemble** | Ensemble | Inverse-RMSE weighted combination of all models |

---

### GET `/api/v1/prediction/algorithms`

List all available prediction algorithms.

**Response** `200 OK`
```json
{
  "algorithms": [
    {
      "key": "linear_regression",
      "name": "Linear Regression",
      "description": "Simple linear regression fits a straight line to historical prices using least squares. Classic statistical method.",
      "type": "Statistical"
    },
    {
      "key": "ema",
      "name": "Exponential Moving Average (EMA)",
      "description": "Applies exponentially decreasing weights to older observations. Widely used in technical analysis.",
      "type": "Technical Analysis"
    },
    {
      "key": "holt",
      "name": "Holt's Linear Trend (Double Exponential Smoothing)",
      "description": "Extension of exponential smoothing that captures both level and trend in time series data.",
      "type": "Time Series"
    },
    {
      "key": "knn",
      "name": "K-Nearest Neighbors (KNN) Regression",
      "description": "Non-parametric method that predicts based on the K most similar historical price patterns.",
      "type": "Machine Learning"
    },
    {
      "key": "ar",
      "name": "Auto-Regressive AR(3) Model",
      "description": "Models the next value as a linear combination of previous values. Foundation of ARIMA.",
      "type": "Time Series"
    },
    {
      "key": "wma",
      "name": "Weighted Moving Average (WMA)",
      "description": "Moving average with linearly increasing weights, giving more importance to recent prices.",
      "type": "Technical Analysis"
    },
    {
      "key": "ensemble",
      "name": "Ensemble (Weighted Average)",
      "description": "Combines predictions from all models using inverse-RMSE weighted averaging.",
      "type": "Ensemble"
    }
  ],
  "total": 7
}
```

---

### GET `/api/v1/prediction/:symbol`

Get ML-based stock price prediction. Defaults to **Ensemble** (all models combined). Use `?algorithm=` to select a specific algorithm.

**URL Parameters:** `symbol` — Company ticker (e.g., `NABIL`)

**Query Parameters (optional):**

| Parameter | Type | Default | Options |
|-----------|------|---------|---------|
| `algorithm` | string | `ensemble` | `linear_regression`, `ema`, `holt`, `knn`, `ar`, `wma`, `ensemble` |

**Example:** `GET /api/v1/prediction/NABIL?algorithm=knn`

**Response** `200 OK`
```json
{
  "prediction": {
    "symbol": "NABIL",
    "current_price": "1250.00",
    "predicted_price": "1267.35",
    "predicted_change": "17.35",
    "predicted_change_pct": "1.39",
    "rmse": 12.4532,
    "mae": 9.8721,
    "algorithm": "K-Nearest Neighbors (KNN-5, window=5)",
    "datapoints": 45
  }
}
```

**Example with Ensemble:** `GET /api/v1/prediction/NABIL`

```json
{
  "prediction": {
    "symbol": "NABIL",
    "current_price": "1250.00",
    "predicted_price": "1262.18",
    "predicted_change": "12.18",
    "predicted_change_pct": "0.97",
    "rmse": 0,
    "mae": 0,
    "algorithm": "Ensemble (Weighted Average of All Models)",
    "datapoints": 45
  }
}
```

---

### GET `/api/v1/prediction/:symbol/compare`

Run **all 7 algorithms** on a stock and return a side-by-side comparison with the best model identified.

**URL Parameters:** `symbol` — Company ticker (e.g., `NABIL`)

**Response** `200 OK`
```json
{
  "comparison": {
    "symbol": "NABIL",
    "current_price": "1250.00",
    "datapoints": 45,
    "best_model": "Holt's Linear Trend (Double Exponential Smoothing)",
    "ensemble": {
      "symbol": "NABIL",
      "current_price": "1250.00",
      "predicted_price": "1262.18",
      "predicted_change": "12.18",
      "predicted_change_pct": "0.97",
      "rmse": 0,
      "mae": 0,
      "algorithm": "Ensemble (Weighted Average of All Models)",
      "datapoints": 45
    },
    "models": [
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1270.50",
        "predicted_change": "20.50",
        "predicted_change_pct": "1.64",
        "rmse": 15.2341,
        "mae": 12.1023,
        "algorithm": "Linear Regression",
        "datapoints": 45
      },
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1255.20",
        "predicted_change": "5.20",
        "predicted_change_pct": "0.42",
        "rmse": 8.5612,
        "mae": 6.4321,
        "algorithm": "Exponential Moving Average (EMA-5)",
        "datapoints": 45
      },
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1263.40",
        "predicted_change": "13.40",
        "predicted_change_pct": "1.07",
        "rmse": 7.2134,
        "mae": 5.8912,
        "algorithm": "Holt's Linear Trend (Double Exponential Smoothing)",
        "datapoints": 45
      },
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1267.35",
        "predicted_change": "17.35",
        "predicted_change_pct": "1.39",
        "rmse": 12.4532,
        "mae": 9.8721,
        "algorithm": "K-Nearest Neighbors (KNN-5, window=5)",
        "datapoints": 45
      },
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1258.90",
        "predicted_change": "8.90",
        "predicted_change_pct": "0.71",
        "rmse": 10.3456,
        "mae": 8.1234,
        "algorithm": "Auto-Regressive AR(3)",
        "datapoints": 45
      },
      {
        "symbol": "NABIL",
        "current_price": "1250.00",
        "predicted_price": "1254.10",
        "predicted_change": "4.10",
        "predicted_change_pct": "0.33",
        "rmse": 9.1234,
        "mae": 7.0123,
        "algorithm": "Weighted Moving Average (WMA-5)",
        "datapoints": 45
      }
    ],
    "model_details": [
      {
        "algorithm": "Linear Regression",
        "predicted": 1270.50,
        "rmse": 15.2341,
        "mae": 12.1023,
        "weight": 0.0982
      },
      {
        "algorithm": "Exponential Moving Average (EMA-5)",
        "predicted": 1255.20,
        "rmse": 8.5612,
        "mae": 6.4321,
        "weight": 0.1748
      },
      {
        "algorithm": "Holt's Linear Trend",
        "predicted": 1263.40,
        "rmse": 7.2134,
        "mae": 5.8912,
        "weight": 0.2074
      },
      {
        "algorithm": "KNN Regression",
        "predicted": 1267.35,
        "rmse": 12.4532,
        "mae": 9.8721,
        "weight": 0.1201
      },
      {
        "algorithm": "Auto-Regressive AR(3)",
        "predicted": 1258.90,
        "rmse": 10.3456,
        "mae": 8.1234,
        "weight": 0.1446
      },
      {
        "algorithm": "Weighted Moving Average",
        "predicted": 1254.10,
        "rmse": 9.1234,
        "mae": 7.0123,
        "weight": 0.1639
      }
    ]
  }
}
```

---

## 7. Wallet (Fiat)

### GET `/api/v1/wallet` 🔒

Get the user's fiat wallet (for top-ups and deposits).

**Headers:** `Authorization: Bearer <token>`

**Response** `200 OK`
```json
{
  "wallet": {
    "id": "uuid-wallet",
    "user_id": "uuid-user",
    "fiat_balance": "50000.00",
    "locked": false,
    "version": 1,
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:00:00Z"
  }
}
```

---

### POST `/api/v1/wallet/topup` 🔒

Add funds to the fiat wallet. The amount is added to both `fiat_balance` and `balance` (available for trading).

**Headers:** `Authorization: Bearer <token>`

**Request Body**
```json
{
  "amount": 50000
}
```

| Field  | Type  | Required | Validation |
|--------|-------|----------|------------|
| amount | float | Yes      | gt=0       |

**Response** `200 OK`
```json
{
  "message": "top-up successful",
  "wallet": {
    "id": "uuid-wallet",
    "user_id": "uuid-user",
    "fiat_balance": "50000.00",
    "locked": false,
    "version": 1
  },
  "transaction": {
    "id": "uuid-tx",
    "user_id": "uuid-user",
    "type": "topup",
    "amount": "50000.00",
    "status": "success",
    "reference_id": "topup_abc123",
    "created_at": "2026-02-23T10:05:00Z"
  }
}
```

**Error** `423 Locked`
```json
{
  "error": "wallet is locked"
}
```

---

### GET `/api/v1/transaction` 🔒

Get the user's fiat wallet transaction history (top-ups, refunds).

**Headers:** `Authorization: Bearer <token>`

**Response** `200 OK`
```json
{
  "message": "user transactions retrieved successfully",
  "data": [
    {
      "id": "uuid-tx",
      "user_id": "uuid-user",
      "type": "topup",
      "amount": "50000.00",
      "status": "success",
      "reference_id": "topup_abc123",
      "created_at": "2026-02-23T10:05:00Z"
    }
  ]
}
```

---

## 8. Trading

### GET `/api/v1/trading/wallet` 🔒

Get the user's trading wallet (balance, invested amount, profit/loss).

**Headers:** `Authorization: Bearer <token>`

**Response** `200 OK`
```json
{
  "wallet": {
    "id": "uuid-vwallet",
    "user_id": "uuid-user",
    "balance": "40000.00",
    "total_invested": "10000.00",
    "total_profit_loss": "250.00",
    "fiat_balance": "50000.00",
    "locked": false,
    "version": 1,
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T11:00:00Z"
  }
}
```

---

### GET `/api/v1/trading/portfolio` 🔒

Get the user's stock portfolio with current values and profit/loss.

**Headers:** `Authorization: Bearer <token>`

**Response** `200 OK`
```json
{
  "total_value": "10250.00",
  "total_invested": "10000.00",
  "total_profit_loss": "250.00",
  "profit_loss_pct": "2.50",
  "items": [
    {
      "id": "uuid-portfolio-item",
      "user_id": "uuid-user",
      "company_id": "uuid-1",
      "quantity": 100,
      "average_price": "100.00",
      "total_invested": "10000.00",
      "current_price": "102.50",
      "current_value": "10250.00",
      "profit_loss": "250.00",
      "profit_loss_pct": "2.50",
      "company_name": "Nabil Bank Limited",
      "company_symbol": "NABIL",
      "company_sector": "Banking"
    }
  ]
}
```

---

### POST `/api/v1/trading/buy` 🔒

Buy shares of a stock. Price impact is calculated via the PriceEngine — buying creates upward pressure on the stock price.

**Headers:** `Authorization: Bearer <token>`

**Request Body**
```json
{
  "symbol": "NABIL",
  "quantity": 100
}
```

| Field    | Type   | Required | Validation |
|----------|--------|----------|------------|
| symbol   | string | Yes      | -          |
| quantity | int    | Yes      | min=1      |

**Response** `200 OK`
```json
{
  "success": true,
  "message": "Stock purchased successfully",
  "quantity": 100,
  "price_per_share": "100.00",
  "total_amount": "10000.00",
  "new_balance": "40000.00",
  "new_market_price": "100.35",
  "price_impact": "0.35",
  "price_impact_pct": "0.35"
}
```

**Error** `400 Bad Request` — Insufficient balance
```json
{
  "success": false,
  "message": "insufficient balance",
  "quantity": 0,
  "price_per_share": "0",
  "total_amount": "0",
  "new_balance": "0"
}
```

**Error** `400 Bad Request` — Insufficient shares
```json
{
  "success": false,
  "message": "Insufficient shares available. Only 500 shares in market",
  "quantity": 0,
  "price_per_share": "0",
  "total_amount": "0",
  "new_balance": "0"
}
```

---

### POST `/api/v1/trading/sell` 🔒

Sell shares of a stock. Selling creates downward pressure on the stock price.

**Headers:** `Authorization: Bearer <token>`

**Request Body**
```json
{
  "symbol": "NABIL",
  "quantity": 50
}
```

| Field    | Type   | Required | Validation |
|----------|--------|----------|------------|
| symbol   | string | Yes      | -          |
| quantity | int    | Yes      | min=1      |

**Response** `200 OK`
```json
{
  "success": true,
  "message": "Stock sold successfully",
  "quantity": 50,
  "price_per_share": "102.50",
  "total_amount": "5125.00",
  "new_balance": "45125.00",
  "new_market_price": "102.15",
  "price_impact": "-0.35",
  "price_impact_pct": "-0.34"
}
```

**Error** `400 Bad Request` — Not enough shares
```json
{
  "success": false,
  "message": "insufficient shares",
  "quantity": 0,
  "price_per_share": "0",
  "total_amount": "0",
  "new_balance": "0"
}
```

---

### GET `/api/v1/trading/transactions` 🔒

Get the user's stock trading transaction history (buy/sell).

**Headers:** `Authorization: Bearer <token>`

**Query Parameters**

| Param  | Type | Default |
|--------|------|---------|
| limit  | int  | 50      |
| offset | int  | 0       |

**Response** `200 OK`
```json
{
  "transactions": [
    {
      "id": "uuid-stx",
      "user_id": "uuid-user",
      "company_id": "uuid-1",
      "type": "buy",
      "quantity": 100,
      "price_per_share": "100.00",
      "total_amount": "10000.00",
      "status": "completed",
      "reference_id": "",
      "created_at": "2026-02-23T11:00:00Z",
      "updated_at": "2026-02-23T11:00:00Z"
    },
    {
      "id": "uuid-stx-2",
      "user_id": "uuid-user",
      "company_id": "uuid-1",
      "type": "sell",
      "quantity": 50,
      "price_per_share": "102.50",
      "total_amount": "5125.00",
      "status": "completed",
      "reference_id": "",
      "created_at": "2026-02-23T11:30:00Z",
      "updated_at": "2026-02-23T11:30:00Z"
    }
  ],
  "limit": 50,
  "offset": 0
}
```

---

## 9. Admin

All admin endpoints require JWT token with `role: admin`.

### PUT `/api/v1/admin/users/:user_id/kyc` 🔒 👑

Update a user's KYC status and role.

**Headers:** `Authorization: Bearer <admin_token>`

**URL Parameters:** `user_id` — Target user UUID

**Request Body**
```json
{
  "kyc_status": "verified",
  "role": "user"
}
```

| Field      | Type   | Required | Allowed Values                         |
|------------|--------|----------|----------------------------------------|
| kyc_status | string | Yes      | pending, verified, rejected, under_review |
| role       | string | Yes      | user, admin                            |

**Response** `200 OK`
```json
{
  "message": "KYC status updated successfully",
  "user": {
    "id": "uuid-user",
    "full_name": "Ram Sharma",
    "email": "ram@example.com",
    "kyc_status": "verified",
    "role": "user"
  }
}
```

---

### POST `/api/v1/admin/seed-stocks` 🔒 👑

Seed the database with 25 Nepali companies, initial prices (Rs 100), and 25 market events per company. Only works on an empty database.

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:** None

**Response** `200 OK`
```json
{
  "success": true,
  "message": "Database seeded successfully",
  "companies_created": 25,
  "prices_created": 25,
  "events_created": 625,
  "notes": "All companies start at ₹100. No test users or transactions created. Users register via API, wallet starts at ₹0."
}
```

**Error** `400 Bad Request` — Already seeded
```json
{
  "error": "Database already seeded. Delete existing data via Supabase dashboard if re-seeding needed."
}
```

---

## Price Engine — How Order-Driven Pricing Works

Prices are **not random**. They change based on actual buy/sell trades:

| Action | Effect |
|--------|--------|
| **Buy** | Price goes **UP** (demand pressure) |
| **Sell** | Price goes **DOWN** (supply pressure) |

**Formula:**
```
impact = BaseSensitivity × √(tradeQuantity / availableShares)
newPrice = currentPrice × (1 + direction × impact + microNoise)
```

**Constants:**
- `BaseSensitivity = 0.5`
- `MicroNoiseRange = ±0.1%`
- `Circuit Breaker = ±10%` daily (like NEPSE)
- `MinPrice = Rs 1.00`

**Example:** Buying 1,000 shares of a stock with 1,000,000 available shares:
```
impact = 0.5 × √(1000 / 1000000) = 0.5 × 0.0316 = 0.0158 (1.58%)
```

---

## Seeded Companies (25)

| Symbol    | Name                           | Sector               |
|-----------|--------------------------------|-----------------------|
| NABIL     | Nabil Bank Limited             | Banking               |
| GBIME     | Global IME Bank Limited        | Banking               |
| NICA      | Nepal Investment Mega Bank     | Banking               |
| NMB       | NMB Bank Limited               | Banking               |
| HBL       | Himalayan Bank Limited         | Banking               |
| NTC       | Nepal Telecom                  | Information Technology|
| NITC      | Nepal IT Corporation           | Information Technology|
| F1SOFT    | F1Soft International           | Information Technology|
| ESewa     | eSewa Digital Services         | Information Technology|
| IMS       | IMS Software Solutions         | Information Technology|
| HYDRO1    | Upper Tamakoshi Hydropower     | Hydropower            |
| BPC       | Butwal Power Company           | Hydropower            |
| CHCL      | Chilime Hydropower             | Hydropower            |
| API       | Api Power Company              | Hydropower            |
| RSHP      | Rosuwa Shyamkhola Hydro Power  | Hydropower            |
| NLIC      | Nepal Life Insurance Company   | Insurance             |
| SICL      | Shikhar Insurance Company      | Insurance             |
| NMBHL     | NMB Health Insurance           | Insurance             |
| APOLLONP  | Apollo Nepal Hospitals         | Pharma                |
| MHPL      | Medical Health Products Limited| Pharma                |
| HDL       | Himalayan Distillery Limited   | Manufacturing         |
| UNL       | Unilever Nepal Limited         | Manufacturing         |
| BNL       | Bottlers Nepal Limited         | Manufacturing         |
| SHIVM     | Shivam Cements Limited         | Manufacturing         |
| DLFNP     | Nepal Housing Development Co.  | Real Estate           |

---

## Error Response Format

All errors follow this format:
```json
{
  "error": "error message here"
}
```

## Authentication

JWT token is returned on login. Include it in all 🔒 endpoints:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

Token is validated by middleware. If expired or invalid:
```json
{
  "error": "unauthorized"
}
```

## Setup Steps

1. **Run SQL migration** — Copy and paste `supabase/migrations/001_complete_schema.sql` into Supabase SQL Editor and execute
2. **Set environment variables** — `SUPABASE_URL`, `SUPABASE_KEY`, `JWT_SECRET`
3. **Seed data** — Either run `go run scripts/main.go` or call `POST /api/v1/admin/seed-stocks` with an admin token
4. **Start server** — `go run cmd/main.go`
