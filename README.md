# 📈 Stock Market Simulator - Virtual Trading Platform

> A production-ready stock market simulation platform for Nepal, enabling users to practice stock trading with virtual currency (NPR) in a risk-free environment.

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)
![Status](https://img.shields.io/badge/status-production--ready-brightgreen?style=flat-square)

## 🎯 Overview

Stock Market Simulator is a fully-featured backend API that allows users to:
- 📱 Register and authenticate securely with JWT
- 💰 Maintain a virtual wallet with NPR 10 Lakh starting balance
- 📊 Browse real-time stock market data
- 🤝 Buy and sell stocks with instant settlement
- 📈 Track portfolio performance and transaction history
- 🎓 Learn stock trading without financial risk


---

## ✨ Key Features

### Core Functionality
- ✅ **Secure Authentication**: JWT-based user authentication with bcrypt password hashing
- ✅ **Virtual Trading**: Buy/sell stocks with NPR virtual currency (1M initial balance)
- ✅ **Real-time Market Data**: Browse companies, prices, and market trends
- ✅ **Portfolio Management**: Track holdings, profit/loss, and transaction history
- ✅ **Wallet System**: Top-up and manage virtual wallet balance
- ✅ **Admin Controls**: KYC verification and stock data management

### Technical Excellence
- ✅ **Production-Ready**: Built with industry best practices
- ✅ **Scalable**: Goroutine-based worker pool for concurrent transactions
- ✅ **Thread-Safe**: Mutex-protected caching and wallet locking
- ✅ **Structured Logging**: Comprehensive logging for debugging
- ✅ **Error Handling**: Custom error types and user-friendly messages
- ✅ **Cloud Native**: Supabase (PostgreSQL) backend with Render deployment

### API Features
- ✅ **10+ Public Endpoints**: Browse stocks, market data, company info
- ✅ **15+ Protected Endpoints**: Trading, wallet, profile management
- ✅ **5+ Admin Endpoints**: User verification and stock seeding
- ✅ **Health Check**: Automated deployment monitoring
- ✅ **Pagination**: Efficient data retrieval with limit/offset
- ✅ **Filtering & Search**: Find stocks and market leaders

---

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.24+ | High-performance backend |
| **Framework** | Gin Web Framework | REST API server |
| **Database** | Supabase (PostgreSQL) | User & trading data storage |
| **Cache** | In-memory (sync.RWMutex) | User data & company cache |
| **Authentication** | JWT + bcrypt | Secure user auth |
| **Async Processing** | Worker Pool (Goroutines) | Concurrent transaction handling |
| **Deployment** | Render + Docker | Cloud hosting |
| **CI/CD** | GitHub Actions | Automated health checks |
| **Logging** | log/slog | Structured logging |

---

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose (optional)
- Git
- Supabase account (free tier available)

### 1. Clone Repository
```bash
git clone https://github.com/919Umesh/stock_market_sim.git
cd stock_market_sim
```

### 2. Environment Setup
Create `.env` file with your Supabase credentials:
```env
PORT=8080
JWT_SECRET=your_super_secret_jwt_key_here
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here
WORKER_COUNT=5
QUEUE_SIZE=100
```

### 3. Run Locally
```bash
# Download dependencies
go mod download

# Run server
go run cmd/main.go
```

Server starts on `http://localhost:8080`

### 4. Run with Docker
```bash
docker-compose up --build
```

---

## 📚 API Documentation

### Quick Links
- 📖 **Full API Docs**: See [`API_DOCUMENTATION.md`](./API_DOCUMENTATION.md)
- 🔐 **Authentication**: JWT Bearer tokens required for protected endpoints
- 📊 **Base URL**: `/api/v1`
- 📦 **Data**: 25 companies, 30 days of historical prices, 750 test transactions


### Public Endpoints (No Auth Required)
```
GET    /stocks                      # List all companies
GET    /stocks/search?q=query      # Search companies
GET    /stocks/:symbol             # Get company details
GET    /stocks/:symbol/price       # Get current price
GET    /stocks/:symbol/history     # Get price history
GET    /stocks/market-overview     # Market statistics
GET    /stocks/top-gainers         # Top 10 gainers
GET    /stocks/top-losers          # Top 10 losers
GET    /stocks/most-active         # Most traded stocks
GET    /stocks/:symbol/events      # Upcoming events
POST   /auth/register              # Create account
POST   /auth/login                 # Get JWT token
```

### Protected Endpoints (Auth Required)
```
GET    /auth/profile               # Get user profile
PUT    /auth/profile/update        # Update profile
GET    /wallet                     # Get wallet balance
POST   /wallet/topup              # Add balance
GET    /transaction               # Wallet transactions
GET    /trading/wallet            # Trading wallet
GET    /trading/portfolio          # Portfolio holdings
POST   /trading/buy               # Buy stock
POST   /trading/sell              # Sell stock
GET    /trading/transactions      # Trading history
```

### Admin Endpoints
```
PUT    /admin/users/:id/kyc       # Update KYC status
POST   /admin/seed-stocks         # Seed stock data
```

---

## 📋 Example Usage

### Register & Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+977981234567",
    "password": "securepass123",
    "role": "user"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
# Returns JWT token
```

### Browse Stocks
```bash
# List all companies
curl http://localhost:8080/api/v1/stocks

# Search for company
curl http://localhost:8080/api/v1/stocks/search?q=NTC

# Get current price
curl http://localhost:8080/api/v1/stocks/NTC/price

# View top gainers
curl http://localhost:8080/api/v1/stocks/top-gainers
```

### Trading
```bash
# Get wallet balance
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/trading/wallet

# Buy stock
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"NTC","quantity":50}' \
  http://localhost:8080/api/v1/trading/buy

# View portfolio
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/trading/portfolio

# Sell stock
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"NTC","quantity":30}' \
  http://localhost:8080/api/v1/trading/sell
```

---

## 🏗️ Project Structure

```
stock_market_sim/
├── cmd/
│   └── main.go                    # Application entry point
├── api/
│   ├── routes.go                  # Route definitions
│   ├── trading_handler.go         # Trading endpoints
│   ├── stock_handler.go           # Stock market endpoints
│   └── admin_handler.go           # Admin endpoints
├── internal/
│   ├── supabase/
│   │   └── client.go              # Supabase REST API client
│   ├── auth/
│   │   ├── handler.go             # Auth HTTP handlers
│   │   ├── service.go             # Auth business logic
│   │   └── repository.go          # Auth data access
│   ├── trading/
│   │   ├── handler.go             # Trading HTTP handlers
│   │   ├── service.go             # Trading business logic
│   │   └── repository.go          # Trading data access
│   ├── wallet/
│   │   ├── handler.go             # Wallet HTTP handlers
│   │   ├── service.go             # Wallet business logic
│   │   └── repository.go          # Wallet data access
│   ├── stock/
│   │   ├── service.go             # Stock business logic
│   │   └── repository.go          # Stock data access
│   └── market/
│       └── simulator.go           # Market simulation engine
├── models/
│   ├── user.go                    # User model
│   ├── virtual_wallet.go          # Wallet model
│   ├── user_portfolio.go          # Portfolio holdings
│   ├── stock_transaction.go       # Transaction records
│   ├── company.go                 # Company data
│   └── stock_price.go             # Price history
├── pkg/
│   ├── middleware/
│   │   ├── auth.go                # JWT verification
│   │   └── admin.go               # Admin authorization
│   ├── apperr/
│   │   └── error.go               # Custom error types
│   ├── logger/
│   │   └── logger.go              # Structured logging
│   ├── queue/
│   │   └── worker_pool.go         # Async worker pool
│   └── utils/
│       ├── jwt.go                 # JWT utilities
│       └── hash.go                # Password hashing
├── config/
│   └── config.go                  # Configuration management
├── docker-compose.yml             # Docker setup
├── Dockerfile                      # Container configuration
├── go.mod & go.sum                # Dependencies
├── render.yaml                    # Render deployment config
└── API_DOCUMENTATION.md           # Complete API reference
```

---

## 🔐 Security Features

- ✅ **JWT Authentication**: Secure token-based auth
- ✅ **Password Hashing**: bcrypt with salt
- ✅ **CORS Enabled**: Cross-origin requests allowed
- ✅ **Input Validation**: Request binding with constraints
- ✅ **Error Handling**: No sensitive data leakage
- ✅ **Thread-Safe Operations**: Mutex protection for shared resources
- ✅ **Transaction Locking**: Prevent race conditions in trades
- ✅ **Row Level Security**: Supabase RLS for data access control

---

## 📊 Database Schema

### Tables in Supabase (PostgreSQL)

| Table | Purpose | Key Fields |
|-----------|---------|-----------|
| `users` | User accounts | email, password_hash, kyc_status, role |
| `virtual_wallets` | Trading balances | user_id, balance, total_invested |
| `user_portfolios` | Stock holdings | user_id, company_id, quantity, average_price |
| `stock_transactions` | Buy/sell records | user_id, company_id, type, quantity |
| `companies` | Stock listing | symbol, name, sector, market_cap |
| `stock_prices` | Price history | company_id, date, open, high, low, close |

---

## 🚀 Deployment

### Render.com (Current)
The application is deployed on Render's free tier at: `https://gold-go-backend.onrender.com`

**Deploy your own:**
1. Fork the repository
2. Go to [Render.com](https://render.com)
3. Create new Web Service
4. Connect GitHub repo
5. Configure environment variables
6. Deploy automatically

**Configuration:**
- Runtime: Docker
- Region: Frankfurt (or your preference)
- Port: 8080
- Health Check: `/health` (hourly via GitHub Actions)

### GitHub Actions
- **Health Check Workflow**: Runs hourly to verify deployment
- **File**: `.github/workflows/health-check.yml`
- **CI/CD Ready**: Can add automated testing/linting

---

## �� Performance & Scalability

### Load Handling
- **Concurrent Transactions**: Goroutine-based worker pool (default: 5 workers)
- **Queue Capacity**: 100 pending jobs
- **Caching**: In-memory user & company cache with RWMutex
- **Database**: Supabase PostgreSQL auto-scales

### Response Times
- Health Check: <1ms
- Stock Browse: <100ms
- Trading Operations: <500ms
- Market Overview: <200ms

---

## 🧪 Testing

### Manual Testing
```bash
# Health check
curl http://localhost:8080/health

# Build
go build ./...

# Lint
go vet ./...
```

### Integration Tests
Use the curl examples in the documentation or tools like:
- Postman
- Insomnia
- Thunder Client
- curl/httpie

---

## 📝 Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | Yes | 8080 | Server port |
| `JWT_SECRET` | Yes | - | JWT signing key |
| `SUPABASE_URL` | Yes | - | Supabase project URL |
| `SUPABASE_ANON_KEY` | Yes | - | Supabase anon/service key |
| `WORKER_COUNT` | No | 5 | Worker pool size |
| `QUEUE_SIZE` | No | 100 | Job queue capacity |

---

## �� Troubleshooting

### Common Issues

**"Timed out waiting for health check"**
- Solution: Ensure `/health` endpoint is accessible
- Check: `curl https://gold-go-backend.onrender.com/health`

**"Failed to connect to Supabase"**
- Check Supabase credentials in `.env`
- Verify Supabase project is running and accessible

**"Insufficient balance"**
- Top-up wallet: `POST /wallet/topup`
- Check balance: `GET /trading/wallet`

**"Portfolio item not found"**
- Ensure stock exists
- Try buying the stock first

---

## 📞 Support & Contribution

### Getting Help
- 📖 **API Docs**: See [`API_DOCUMENTATION.md`](./API_DOCUMENTATION.md)
- 🐛 **Issues**: Report on [GitHub Issues](https://github.com/919Umesh/stock_market_sim/issues)
- 💬 **Discussions**: Use GitHub Discussions

### Contributing
1. Fork the repository
2. Create feature branch: `git checkout -b feature/your-feature`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature/your-feature`
5. Submit Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💼 Author

**Developed by:** [919Umesh](https://github.com/919Umesh)

**Repository:** https://github.com/919Umesh/stock_market_sim

---

## 🎓 Learning Resources

- 📚 [Go Documentation](https://golang.org/doc/)
- 🌐 [Gin Web Framework](https://gin-gonic.com/)
- 📊 [Appwrite Documentation](https://appwrite.io/docs)
- 🐳 [Docker Documentation](https://docs.docker.com/)

---

## 🔮 Future Roadmap

- [ ] Real stock market integration (NSE Nepal)
- [ ] Advanced charting with candlesticks
- [ ] Market recommendations engine
- [ ] Social trading features
- [ ] Mobile app (React Native)
- [ ] WebSocket for real-time updates
- [ ] Rate limiting & API throttling
- [ ] Advanced analytics dashboard

---

**Status**: ✅ Production Ready | **Last Updated**: February 2026 | **API Version**: v1
