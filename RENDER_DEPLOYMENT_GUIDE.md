# 🚀 Complete Guide: Deploy to Render with Supabase

## Overview
This guide walks you through deploying your Go backend to Render and connecting it to Supabase.

---

## Step 1: Prepare Supabase Project

### 1.1 Create a New Supabase Project (if needed)
1. Go to [supabase.com](https://supabase.com)
2. Sign in or create account
3. Click "New Project"
4. Fill in:
   - **Project Name**: `gold-go-backend` (or your choice)
   - **Database Password**: Generate strong password (save it!)
   - **Region**: Choose closest to your users
5. Click "Create new project" (takes 1-2 minutes)

### 1.2 Run Database Migration
1. Once project is created, go to **SQL Editor** (left sidebar)
2. Click "New Query"
3. Copy the entire contents of `supabase/migrations/001_complete_schema.sql`
4. Paste into SQL Editor
5. Click "Run" button (bottom right)
6. Wait for success message ✅

### 1.3 Get Your Supabase Credentials

In Supabase Dashboard, go to **Settings** → **API**:

#### Find These 3 Values:

**A. Project URL** (looks like: `https://xxxxxxxxx.supabase.co`)
```
Settings → API → Project URL
Copy this URL
```

**B. Anon Key** (public, safe to expose in frontend)
```
Settings → API → Project API keys → anon key
Copy this key (looks like: eyJh...)
```

**C. Service Role Key** (SECRET! for backend only)
```
Settings → API → Project API keys → service_role (use this one for backend!)
⚠️ NEVER expose this in frontend or public repos
```

> **Important**: Use `service_role` key for backend, not `anon` key!

---

## Step 2: Prepare GitHub Repository

### 2.1 Push Code to GitHub
```bash
# From your project root
git add .
git commit -m "Ready for Render deployment"
git push origin main
```

### 2.2 Ensure `.gitignore` has these entries
```
.env
.env.local
bin/
main
/dist/
```

Your `.gitignore` should already have these. Verify it does.

---

## Step 3: Connect Render to GitHub

### 3.1 Create Render Account
1. Go to [render.com](https://render.com)
2. Sign up with GitHub (recommended for auto-deploys)
3. Authorize Render to access your GitHub account

### 3.2 Create New Web Service
1. In Render dashboard, click **"New +"** → **"Web Service"**
2. Select your GitHub repository
3. Configure:
   - **Name**: `gold-go-backend` (or your choice)
   - **Environment**: `Docker`
   - **Region**: `Frankfurt` (or closest to you)
   - **Plan**: `Free` (can upgrade later)
   - **Auto-Deploy**: Toggle ON (auto-deploy on git push)

---

## Step 4: Set Environment Variables in Render

### 4.1 Add Environment Variables
After creating the service, go to **Settings** → **Environment**:

Add these 4 variables:

| Key | Value | Source |
|-----|-------|--------|
| `SUPABASE_URL` | `https://xxxxxxxxx.supabase.co` | Supabase Settings → API → Project URL |
| `SUPABASE_ANON_KEY` | `eyJh...` | Supabase Settings → API → anon key |
| `JWT_SECRET` | `your-super-secret-key-min-32-chars` | Generate random strong password |
| `PORT` | `8080` | Default (Render provides this automatically, but set anyway) |
| `WORKER_COUNT` | `5` | Optional, default is 5 |
| `QUEUE_SIZE` | `100` | Optional, default is 100 |

**Example for SUPABASE_URL:**
```
https://abcdef123456.supabase.co
```

**Example for JWT_SECRET** (generate online at [randomkeygen.com](https://randomkeygen.com)):
```
7$9#mK@2vL8*pQ%4xN&6yR+5zW^3cD!1bF0eG8hJ5kMoP2sT7uV4wX6yZ1aB3cD5eF
```

### 4.2 Important Notes on Keys
- ✅ **SUPABASE_ANON_KEY**: Safe to use in backend (not frontend)
- ✅ **SERVICE_ROLE_KEY**: NOT NEEDED (you have anon key, which is enough)
- ✅ **JWT_SECRET**: Keep unique and strong (min 32 chars recommended)

> **Why Service Role Key is NOT included**: Your backend code uses the anon key for Supabase access. The service_role key is only needed if you're doing advanced admin operations, which this code doesn't do.

---

## Step 5: Deploy

### 5.1 Manual Deploy (First Time)
1. In Render dashboard, click the service name
2. Click **"Manual Deploy"** → **"Deploy latest commit"**
3. Watch logs (click **"Logs"** tab)
4. Wait for deployment to complete (2-3 minutes)
5. Once complete, you'll see a green checkmark ✅

### 5.2 Auto-Deploy (On Every Git Push)
Since you enabled "Auto-Deploy", any future push to `main` branch will auto-deploy:

```bash
git add .
git commit -m "Update API endpoint"
git push origin main
# Render automatically deploys!
```

---

## Step 6: Verify Deployment

### 6.1 Get Your Live URL
In Render dashboard:
1. Click on your service name
2. Look at the top - you'll see your live URL
   ```
   https://gold-go-backend.onrender.com
   ```

### 6.2 Test Health Endpoint
Open in browser or use curl:
```bash
curl https://gold-go-backend.onrender.com/health
```

Expected response:
```json
{"status":"ok"}
```

### 6.3 Test Registration Endpoint
```bash
curl -X POST https://gold-go-backend.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "test@example.com",
    "phone": "9841234567",
    "password": "password123",
    "role": "user"
  }'
```

Expected response:
```json
{
  "message": "user registered successfully",
  "user": {...}
}
```

---

## Step 7: Common Issues & Troubleshooting

### Issue: "SUPABASE_URL is required"
**Solution**: You forgot to add `SUPABASE_URL` in Render environment variables
- Go to Settings → Environment
- Add `SUPABASE_URL` = your Supabase project URL

### Issue: "Failed to connect to Supabase"
**Solution**: Check these:
1. Is `SUPABASE_URL` correct? (should start with `https://` and end with `.supabase.co`)
2. Is `SUPABASE_ANON_KEY` correct? (should be long string starting with `eyJ`)
3. Did you run the migration in Supabase? Check SQL Editor history
4. Check Render logs: Settings → Logs, scroll for errors

### Issue: "Deployment failed"
**Solution**:
1. Check Render logs for error message
2. Common causes:
   - Go code compilation error - fix in local repo, push again
   - Missing environment variable - add to Render Settings → Environment
   - Port conflicts - unlikely on Render

### Issue: "Build failing - Go version mismatch"
**Solution**: Render uses Go 1.21 by default. Our code uses Go 1.24.
- Add `runtime.txt` to root with: `go1.24`
- Or ensure `go.mod` has correct version

---

## Step 8: Setup Continuous Integration (Optional)

### 8.1 Create GitHub Actions Workflow
Create file: `.github/workflows/deploy.yml`

```yaml
name: Deploy to Render

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Test Go build
        run: go build ./cmd/main.go
```

This tests your code before Render deploys it.

---

## Step 9: Monitor Your Deployment

### 9.1 View Logs
In Render dashboard → Logs tab:
- See real-time server activity
- Spot errors immediately
- Debug issues

### 9.2 Set up Health Checks
Render automatically pings `/health` endpoint every 5 minutes
- If it fails 3 times, deployment restarts automatically
- Your code already has this: `api/routes.go` line ~47

### 9.3 View Metrics
Click **Metrics** tab to see:
- CPU usage
- Memory usage
- Request count
- Restart history

---

## Step 10: Production Best Practices

### 10.1 Upgrade from Free Tier (Optional)
Free tier limitations:
- Spins down after 15 mins of inactivity
- Can't scale to multiple instances

Upgrade to **Starter** ($7/month):
- No spin-down (always on)
- Better performance
- Better support

In Render dashboard → Settings → Plan → Upgrade

### 10.2 Setup Custom Domain (Optional)
1. Buy domain (GoDaddy, Namecheap, etc.)
2. In Render Settings → Custom Domain
3. Point domain DNS to Render
4. Get free SSL certificate automatically

### 10.3 Setup Database Backups
In Supabase:
1. Settings → Backups
2. Enable automatic backups (daily)
3. Can restore to any point in time

### 10.4 Monitor Performance
Use tools:
- [Uptime Robot](https://uptimerobot.com) - monitor if service is up
- New Relic - detailed performance monitoring
- DataDog - comprehensive observability

---

## Step 11: Seed Sample Data

Once deployed, seed your database:

### Option A: Using Your Local Script
```bash
# From local machine
export SUPABASE_URL="https://your-project.supabase.co"
export SUPABASE_ANON_KEY="your-anon-key"
go run scripts/main.go
```

### Option B: Using curl (recommended)
```bash
# Call the admin seed endpoint
curl -X POST https://gold-go-backend.onrender.com/api/v1/admin/seed-stocks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

But first, you need:
1. Admin user account
2. KYC verified
3. Valid JWT token

**Easier way**: Run local script, which is already configured.

---

## Summary Checklist

- [ ] Created Supabase project
- [ ] Ran migration SQL in Supabase
- [ ] Got 3 Supabase credentials (URL, anon key, JWT secret)
- [ ] Pushed code to GitHub
- [ ] Created Render account
- [ ] Created Web Service on Render
- [ ] Added 4 environment variables
- [ ] Deployed service
- [ ] Tested `/health` endpoint
- [ ] Tested `/api/v1/auth/register` endpoint
- [ ] Seeded sample data with `go run scripts/main.go`

---

## Your Live API is Ready! 🎉

**Endpoint**: `https://gold-go-backend.onrender.com`

**All Endpoints Available**:
- `GET /health` - Health check
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/stocks` - List all stocks
- `GET /api/v1/stocks/:symbol` - Get stock details
- `POST /api/v1/trading/buy` - Buy stock (protected)
- `POST /api/v1/trading/sell` - Sell stock (protected)
- ... and 20+ more endpoints

---

## Need Help?

- **Render Support**: https://render.com/docs
- **Supabase Docs**: https://supabase.com/docs
- **Go Best Practices**: https://golang.org/doc/effective_go

Good luck with your deployment! 🚀
