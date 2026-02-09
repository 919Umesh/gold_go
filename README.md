# Gold Investment Backend

# Inspiration 
Inspired by the success of Jio Gold in India, we're bringing digital gold investment to Nepal. While Nepal may be smaller in scale, the potential is immense. Gold isn't just a material asset - it carries deep emotional and cultural significance for our people.

Built with the vision of making gold investment accessible, secure, and convenient for every Nepali household.

Contact:
Umesh Shahi
📞 9868732774
📧 thakuriumesh919@gmail.com

## 🚀 Features

- **User Authentication** - JWT-based secure authentication
- **Wallet Management** - Fiat balance and gold grams management
- **Gold Trading** - Buy and sell gold with real-time pricing
- **Price Updates** - Automated background gold price updates
- **Transaction History** - Complete audit trail of all transactions
- **Concurrent Safety** - Thread-safe operations with mutex protection
- **RESTful API** - Clean and well-structured API endpoints

## 🛠️ Technology Stack

- **Backend**: Go 1.20+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Security**: bcrypt password hashing

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- Git

## ⚙️ Installation

### 1. Clone the Repository
```bash
git clone https://github.com/919Umesh/gold_go.git
cd gold_go
```

### 2. Environment Configuration
Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gold_investment
DB_PORT=5432

# Server Configuration
PORT=8080
JWT_SECRET=your_super_secure_jwt_secret_key_here_min_32_chars

# Application Settings
WORKER_COUNT=5
QUEUE_SIZE=100
ENVIRONMENT=development
LOG_LEVEL=info

# Gold Provider (Mock for development)
GOLD_PROVIDER_URL=http://localhost:9000
```

### 3. Database Setup
```sql
-- Connect to PostgreSQL and create database
CREATE DATABASE gold_investment;

-- Or use existing database, the tables will be created automatically
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Run the Application
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080` and automatically create all necessary database tables.

## 🗄️ Database Schema

The system automatically creates the following tables:

### Users Table
- User accounts with authentication details
- KYC status tracking
- Contact information

### Wallets Table
- Fiat balance (NPR)
- Gold grams holdings
- Optimistic locking for concurrent safety

### Transactions Table
- Complete transaction history
- Buy/sell/topup operations
- Status tracking

### Gold Prices Table
- Historical gold prices
- Real-time price updates
- Source tracking

## 🔑 API Endpoints

### Authentication Endpoints

#### Register User
- **POST** `/api/v1/auth/register`
- **Body**:
```json
{
  "full_name": "Umesh Shahi",
  "email": "thakuriumesh919@gmail.com",
  "phone": "9868732774",
  "password": "Thakuri@8848"
}
```

#### Login
- **POST** `/api/v1/auth/login`
- **Body**:
```json
{
  "email": "thakuriumesh919@gmail.com",
  "password": "Thakuri@8848"
}
```

#### Get Profile
- **GET** `/api/v1/auth/profile`
- **Headers**: `Authorization: Bearer <token>`

### Wallet Endpoints (Protected)

#### Get Wallet Balance
- **GET** `/api/v1/wallet`
- **Headers**: `Authorization: Bearer <token>`

#### Top Up Balance
- **POST** `/api/v1/wallet/topup`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "amount": 5000.00
}
```

#### Buy Gold
- **POST** `/api/v1/wallet/buy`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "grams": 2.5,
  "price_per_gram": 6500.00
}
```

#### Sell Gold
- **POST** `/api/v1/wallet/sell`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
  "grams": 1.0,
  "price_per_gram": 6600.00
}
```

### Gold Price Endpoints (Public)

#### Get Current Price
- **GET** `/api/v1/gold/price`

#### Get Price History
- **GET** `/api/v1/gold/history?days=7`

### Health Check
- **GET** `/health`

## 🎯 Usage Examples

### Complete Workflow

1. **Register a new user**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Umesh Shahi","email":"thakuriumesh919@gmail.com","phone":"9868732774","password":"Thakuri@8848"}'
```

2. **Login to get JWT token**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"thakuriumesh919@gmail.com","password":"Thakuri@8848"}'
```

3. **Use token for protected routes**
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get wallet balance
curl -X GET http://localhost:8080/api/v1/wallet \
  -H "Authorization: Bearer $TOKEN"

# Top up balance
curl -X POST http://localhost:8080/api/v1/wallet/topup \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount": 10000}'

# Buy gold
curl -X POST http://localhost:8080/api/v1/wallet/buy \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"grams": 1.5, "price_per_gram": 6500}'
```

## 🔧 Advanced Go Features Implemented

### Goroutines
- Background gold price updates every 30 seconds
- Worker pools for transaction processing
- Concurrent request handling

### Channels
- Job queue for asynchronous transaction processing
- Communication between price updater and consumers
- Buffered channels for efficient resource usage

### Pointers
- Efficient struct passing in method receivers
- Memory optimization for large data structures
- Pointer-based database operations

### Interfaces
- Dependency injection for testability
- Repository pattern for data access
- Service abstraction layers

### Context
- Request cancellation and timeouts
- Graceful shutdown handling
- Background job management

### Sync Package
- Mutexes for concurrent wallet access
- RWMutex for price cache
- WaitGroups for worker coordination

### Error Handling
- Proper error wrapping with context
- Structured error responses
- Panic recovery middleware

### Struct Tags
- JSON serialization/deserialization
- GORM database mappings
- Validation rules

## 🏗️ Project Structure

```
gold_investment_backend/
├── cmd/
│   └── main.go                 
├── config/
│   ├── config.go             
│   └── database.go            
├── internal/
│   ├── auth/                  
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── wallet/                
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── gold/                  
│   │   ├── handler.go
│   │   └── service.go
│   └── transaction/           
├── pkg/
│   ├── middleware/            
│   │   └── auth.go
│   ├── utils/                 
│   │   ├── hash.go
│   │   └── jwt.go
│   └── queue/                 
├── models/                    
│   ├── user.go
│   ├── wallet.go
│   ├── transaction.go
│   └── gold_price.go
├── api/
│   └── routes.go             
└── go.mod                     
```

## 🔒 Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Secure token expiration (72 hours)
- Input validation and sanitization
- SQL injection prevention with GORM
- Concurrent access protection

## 📈 Performance Features

- Database connection pooling
- In-memory price caching
- Background job processing
- Optimistic locking for wallets
- Efficient goroutine management

## 🚀 Deployment

### Production Deployment on Render (Recommended)

This project is configured for automatic deployment to Render with Redis support.

#### Quick Deploy

1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Deploy to Render"
   git push origin main
   ```

2. **Connect to Render**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click "New +" → "Blueprint"
   - Select your repository
   - Render will auto-detect `render.yaml`

3. **Set Environment Variables**
   - In Render dashboard, set `APPWRITE_API_KEY` securely
   - All other variables are pre-configured

4. **Deploy**
   - Click "Apply" to create services
   - Your API will be live at: `https://your-service.onrender.com`

#### Architecture

- **Backend**: Go service on Render (auto-scaling)
- **Database**: Appwrite Cloud (managed)
- **Cache**: Redis on Render (managed)
- **CI/CD**: GitHub Actions + Render auto-deploy

#### Deployment Files

- `Dockerfile` - Multi-stage Docker build
- `render.yaml` - Infrastructure as Code
- `docker-compose.yml` - Local development
- `.github/workflows/ci-cd.yml` - CI/CD pipeline
- `Makefile` - Development commands

For detailed deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md)

### Local Development with Docker

```bash
# Build and run with Docker Compose
make docker-run

# Or manually
docker-compose up --build

# Stop services
make docker-stop
```

### Manual Deployment

```bash
# Build binary
make build

# Run binary
./bin/server
```

### Development Commands

```bash
make help           # Show all available commands
make test           # Run tests
make lint           # Run linter
make deploy-check   # Validate deployment config
make pre-push       # Run all checks before pushing
```

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```


## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request


## 🔮 Future Enhancements

- Real gold provider integration
- Payment gateway integration (eSewa, Khalti)
- SMS/Email notifications
- Admin dashboard
- Advanced reporting and analytics
- Mobile app support

---

**Built with  using Go, Gin, and PostgreSQL**