# Deployment Guide: Seeding Stock Market Data on Render

## Step 1: Deploy Updated Code to Render

First, commit and push all the new stock market code to GitHub:

```bash
git add .
git commit -m "Add virtual stock market trading system with seed endpoint"
git push origin main
```

Render will automatically detect the changes and redeploy your application.

## Step 2: Wait for Deployment to Complete

Monitor your Render dashboard until the deployment shows "Live" status.

## Step 3: Seed the Remote Database

Once deployed, you need to seed the stock market data **ONE TIME ONLY**.

### Option A: Using the Admin API Endpoint (Recommended)

1. **Login as Admin** to get your JWT token:

```bash
curl -X POST https://your-app.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-admin-email@example.com",
    "password": "your-password"
  }'
```

2. **Call the Seed Endpoint** (requires admin role):

```bash
curl -X POST https://your-app.onrender.com/api/v1/admin/seed-stocks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Database seeded successfully",
  "companies_created": 21,
  "prices_created": 5355,
  "events_created": 20
}
```

> **⚠️ Important:** This endpoint can only be called once. If you try to call it again, you'll get an error saying the database is already seeded.

### Option B: Using Render Shell

Alternatively, you can use Render's shell feature:

1. Go to your Render dashboard
2. Click on your `gold-go-app` service
3. Click on the "Shell" tab
4. Run the seed script:

```bash
go run scripts/seed_stocks.go
```

## Step 4: Verify the Data

Test that the stock data is available:

```bash
# Get list of companies
curl https://your-app.onrender.com/api/v1/stocks

# Get market overview
curl https://your-app.onrender.com/api/v1/stocks/market-overview

# Get specific company
curl https://your-app.onrender.com/api/v1/stocks/NABIL
```

## Step 5: Test Trading

1. **Register a new user:**

```bash
curl -X POST https://your-app.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "9841234567",
    "password": "Test123!",
    "role": "user"
  }'
```

2. **Login:**

```bash
curl -X POST https://your-app.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }'
```

3. **Check virtual wallet** (should have NPR 1,000,000):

```bash
curl https://your-app.onrender.com/api/v1/trading/wallet \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

4. **Buy some stocks:**

```bash
curl -X POST https://your-app.onrender.com/api/v1/trading/buy \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "NABIL",
    "quantity": 10
  }'
```

5. **Check portfolio:**

```bash
curl https://your-app.onrender.com/api/v1/trading/portfolio \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## What Gets Seeded

- **21 Nepali Companies** across 6 sectors:
  - Banking (5 companies)
  - Hydropower (4 companies)
  - Insurance (4 companies)
  - Manufacturing (3 companies)
  - Hotels (3 companies)
  - Finance (2 companies)

- **1 Year of Historical Price Data** (~5,355 price records)
  - Daily OHLCV data
  - Realistic price movements
  - Trading volume simulation

- **20 Market Events**
  - Earnings announcements
  - News events
  - Dividend declarations

## Market Simulator

The market simulator will automatically start when your app deploys and will:
- Update stock prices every 5 minutes during market hours
- Respect Nepal Stock Exchange timing (Sun-Thu, 11 AM - 3 PM)
- Use Geometric Brownian Motion for realistic price movements

## Troubleshooting

### "Database already seeded" Error

If you get this error, it means the data is already in the database. You don't need to seed again unless you want to reset everything.

To reset (⚠️ **WARNING: This deletes all stock data**):

```sql
-- Connect to your Render PostgreSQL database and run:
DELETE FROM stock_predictions;
DELETE FROM market_events;
DELETE FROM stock_transactions;
DELETE FROM user_portfolios;
DELETE FROM virtual_wallets;
DELETE FROM stock_prices;
DELETE FROM companies;
```

Then you can call the seed endpoint again.

### Seed Endpoint Returns 401 Unauthorized

Make sure:
1. You're logged in as a user with `admin` role
2. Your JWT token is valid and included in the Authorization header
3. The admin middleware is working correctly

## Next Steps

After seeding:
1. ✅ Test all stock market endpoints
2. ✅ Verify market simulator is running
3. ✅ Test buy/sell functionality
4. ✅ Monitor Render logs for any errors
5. 📱 Start building your frontend!
