# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐         ┌──────────────────┐             │
│  │  Flutter Mobile  │         │   Flutter Web    │             │
│  │      App         │         │      App         │             │
│  └────────┬─────────┘         └────────┬─────────┘             │
│           │                            │                        │
│           └────────────┬───────────────┘                        │
│                        │                                        │
│                        │ HTTPS/REST API                         │
│                        │                                        │
└────────────────────────┼────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    RENDER PLATFORM (Cloud)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐    │
│  │              Go Backend Service                        │    │
│  │  ┌──────────────────────────────────────────────────┐  │    │
│  │  │  API Layer (Gin Framework)                       │  │    │
│  │  │  - Authentication (JWT)                          │  │    │
│  │  │  - Stock Trading                                 │  │    │
│  │  │  - Wallet Management                             │  │    │
│  │  │  - Rate Limiting                                 │  │    │
│  │  │  - CORS                                          │  │    │
│  │  └──────────────────────────────────────────────────┘  │    │
│  │                         │                               │    │
│  │  ┌──────────────────────┼──────────────────────────┐   │    │
│  │  │  Service Layer       │                          │   │    │
│  │  │  - Business Logic    │                          │   │    │
│  │  │  - Validation        │                          │   │    │
│  │  │  - Transaction Mgmt  │                          │   │    │
│  │  └──────────────────────┼──────────────────────────┘   │    │
│  │                         │                               │    │
│  │  ┌──────────────────────┼──────────────────────────┐   │    │
│  │  │  Repository Layer    │                          │   │    │
│  │  │  - Data Access       │                          │   │    │
│  │  │  - Appwrite Client   │                          │   │    │
│  │  └──────────────────────┼──────────────────────────┘   │    │
│  │                         │                               │    │
│  │  ┌──────────────────────┼──────────────────────────┐   │    │
│  │  │  Worker Pool         │                          │   │    │
│  │  │  - Background Jobs   │                          │   │    │
│  │  │  - Async Processing  │                          │   │    │
│  │  └──────────────────────┴──────────────────────────┘   │    │
│  │                                                          │    │
│  │  Port: 8080 (Internal)                                  │    │
│  │  URL: https://gold-go-backend.onrender.com              │    │
│  └────────────┬────────────────────────┬────────────────────┘    │
│               │                        │                         │
│               │                        │                         │
│  ┌────────────▼──────────┐  ┌─────────▼──────────────────┐     │
│  │   Redis Cache         │  │   Appwrite Cloud           │     │
│  │   (Managed)           │  │   (External)               │     │
│  │                       │  │                            │     │
│  │  - Session Cache      │  │  - Database (NoSQL)        │     │
│  │  - Rate Limiting      │  │  - Authentication          │     │
│  │  - Price Cache        │  │  - User Management         │     │
│  │                       │  │  - Collections:            │     │
│  │  Internal Network     │  │    * Users                 │     │
│  │  Port: 6379           │  │    * Wallets               │     │
│  └───────────────────────┘  │    * Transactions          │     │
│                             │    * Companies             │     │
│                             │    * Stock Prices          │     │
│                             │    * Portfolios            │     │
│                             │                            │     │
│                             │  Endpoint:                 │     │
│                             │  fra.cloud.appwrite.io     │     │
│                             └────────────────────────────┘     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. User Registration/Login
```
Flutter App → Backend API → Appwrite Auth → JWT Token → Flutter App
```

### 2. Stock Trading
```
Flutter App → Backend API → Redis (Check Cache) → Appwrite DB
                          ↓
                    Worker Pool (Async Processing)
                          ↓
                    Transaction Processing
                          ↓
                    Wallet Update
```

### 3. Price Updates
```
Background Worker → External API → Redis Cache → Appwrite DB
                                              ↓
                                        Real-time Updates
```

## Technology Stack

### Backend
- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Authentication**: JWT (JSON Web Tokens)
- **Concurrency**: Goroutines, Channels, Worker Pools

### Database & Storage
- **Primary DB**: Appwrite Cloud (NoSQL)
- **Cache**: Redis (Managed by Render)
- **Session Store**: Redis

### Infrastructure
- **Hosting**: Render (PaaS)
- **Container**: Docker (Multi-stage build)
- **CI/CD**: GitHub Actions + Render Auto-deploy

### Development Tools
- **Version Control**: Git/GitHub
- **Dependency Management**: Go Modules
- **Build Tool**: Make
- **Testing**: Go testing package

## Security Layers

```
┌─────────────────────────────────────────┐
│  1. HTTPS/TLS (Render Automatic)        │
├─────────────────────────────────────────┤
│  2. CORS (Configured in Backend)        │
├─────────────────────────────────────────┤
│  3. Rate Limiting (Redis-based)         │
├─────────────────────────────────────────┤
│  4. JWT Authentication                  │
├─────────────────────────────────────────┤
│  5. Input Validation                    │
├─────────────────────────────────────────┤
│  6. Environment Variables (Secrets)     │
├─────────────────────────────────────────┤
│  7. Appwrite Security Rules             │
└─────────────────────────────────────────┘
```

## Deployment Pipeline

```
Developer
    │
    │ git push
    ▼
GitHub Repository
    │
    ├─────────────────────┬──────────────────────┐
    │                     │                      │
    ▼                     ▼                      ▼
GitHub Actions      Render Platform      Code Review
    │                     │
    │ Run Tests          │ Detect Push
    │ Lint Code          │
    │ Build Docker       │
    │ Security Scan      │
    │                     │
    ▼                     ▼
  Pass/Fail          Build Docker Image
    │                     │
    │                     ▼
    │              Deploy to Render
    │                     │
    │                     ├─────────────┬──────────────┐
    │                     │             │              │
    │                     ▼             ▼              ▼
    │              Backend Service  Redis Service  Health Check
    │                     │             │              │
    │                     └─────────────┴──────────────┘
    │                                   │
    │                                   ▼
    │                            Live Production
    │                                   │
    └───────────────────────────────────┘
                 Feedback Loop
```

## Scaling Strategy

### Current (Free Tier)
- **Backend**: 1 instance (spins down after 15 min)
- **Redis**: 25MB storage
- **Cost**: $0/month

### Production (Recommended)
- **Backend**: Always-on, auto-scaling
- **Redis**: 256MB storage, persistence
- **Cost**: ~$17/month

### Future Scaling
- **Horizontal Scaling**: Multiple backend instances
- **Load Balancing**: Automatic (Render)
- **Database Sharding**: Appwrite scaling
- **CDN**: For static assets
- **Monitoring**: Prometheus + Grafana

## Monitoring & Observability

```
┌─────────────────────────────────────────┐
│           Render Dashboard              │
├─────────────────────────────────────────┤
│  • Real-time Logs                       │
│  • CPU/Memory Metrics                   │
│  • Request Count                        │
│  • Response Times                       │
│  • Error Rates                          │
│  • Health Check Status                  │
└─────────────────────────────────────────┘
```

## High Availability

- **Health Checks**: Every 30 seconds
- **Auto-Restart**: On failure
- **Zero-Downtime Deploy**: Blue-green deployment
- **Automatic Rollback**: On health check failure
- **99.9% Uptime SLA**: Render guarantee (paid plans)

## Disaster Recovery

1. **Database**: Appwrite automatic backups
2. **Code**: Git version control
3. **Configuration**: Infrastructure as Code (render.yaml)
4. **Secrets**: Render encrypted storage
5. **Rollback**: Instant via Render dashboard

## Performance Optimization

- **Caching**: Redis for frequently accessed data
- **Connection Pooling**: Database connections
- **Goroutines**: Concurrent request handling
- **Worker Pool**: Background job processing
- **CDN**: Static asset delivery (future)
- **Compression**: Gzip responses

## Cost Optimization

### Current Setup (Free)
- Backend: $0 (with limitations)
- Redis: $0 (25MB)
- Appwrite: $0 (free tier)
- **Total**: $0/month

### Production Setup
- Backend Starter: $7/month
- Redis Starter: $10/month
- Appwrite Pro: $15/month (optional)
- **Total**: $17-32/month

### Enterprise Setup
- Backend Standard: $25/month
- Redis Standard: $25/month
- Appwrite Pro: $15/month
- Custom Domain: Included
- **Total**: $65/month
