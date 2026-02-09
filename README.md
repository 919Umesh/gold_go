# Gold Investment & Stock Market Simulator (Nepal)

Inspired by digital gold investment successes, this platform brings accessible gold and stock market simulation to Nepal. Built for security, scale, and ease of use.

## 🚀 Key Features

*   **Dual Investment System**: Manage both Digital Gold and Virtual Stock portfolios.
*   **User Authentication**: Secure JWT-based auth with Appwrite Cloud integration.
*   **Virtual Wallet**: NPR-based virtual balance for risk-free trading practice.
*   **Real-time Simulation**: Automated gold price updates and stock market events.
*   **Performance**: Redis-based rate limiting and caching for high-speed API responses.
*   **Worker Pool**: Asynchronous background processing for concurrent safety.

## 🛠️ Technology Stack

*   **Backend**: Go 1.24+ (Gin Framework)
*   **Database**: Appwrite Cloud (NoSQL)
*   **Cache/Speed**: Redis (Managed)
*   **Deployment**: Render (Dockerized)
*   **DevOps**: GitHub Actions (CI/CD)

## ⚙️ Quick Start (Development)

1.  **Environment Setup**:
    Copy `.env.example` to `.env` and fill in your Appwrite credentials.
    ```env
    PORT=8080
    JWT_SECRET=your_secret_key
    APPWRITE_ENDPOINT=https://fra.cloud.appwrite.io/v1
    APPWRITE_PROJECT_ID=your_project_id
    APPWRITE_API_KEY=your_api_key
    APPWRITE_DATABASE_ID=your_db_id
    REDIS_URL=redis://localhost:6379
    ```

2.  **Run with Docker**:
    ```bash
    docker-compose up --build
    ```

3.  **Run Locally**:
    ```bash
    go mod tidy
    go run cmd/main.go
    ```

## 🔑 Primary API Endpoints

*   `POST /api/v1/auth/register` - New user signup
*   `POST /api/v1/auth/login` - Secure login (returns JWT)
*   `GET /api/v1/stocks` - List simulated companies
*   `POST /api/v1/trading/buy` - Execute a virtual trade
*   `GET /api/v1/trading/wallet` - Check balance and profit/loss

For a full list of endpoints, see [API_ENDPOINTS.md](./API_ENDPOINTS.md).

## 🚀 Deployment to Render

This project is configured for **Render Blueprint**. 

1.  Push your changes to GitHub.
2.  In [Render Dashboard](https://dashboard.render.com), create a **New Blueprint**.
3.  Connect this repository. Render will auto-detect `render.yaml` and set up the Backend + Redis.
4.  Manually add `APPWRITE_API_KEY` in the Render environment settings.

---
**Maintained by Umesh Shahi**
📞 9868732774 | 📧 thakuriumesh919@gmail.com