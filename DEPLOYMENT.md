# Deployment Instructions

Follow these steps to deploy the application to a production Linux server (e.g., Ubuntu).

## Prerequisites
-   A Linux server (Ubuntu 22.04 recommended).
-   Domain name pointing to your server's IP.
-   Docker and Docker Compose installed.

## 1. Install Docker & Docker Compose
```bash
sudo apt update
sudo apt install -y docker.io docker-compose
sudo systemctl enable --now docker
```

## 2. Deploy Code
Clone your repository to the server:
```bash
git clone https://github.com/your/repo.git /opt/gold_go
cd /opt/gold_go
```

## 3. Configure Environment
Create a production `.env` file (do NOT commit this to Git).
```bash
cp .env.example .env
nano .env
```
**Example .env:**
```ini
DB_USER=gold_user
DB_PASSWORD=secure_prod_password_123
DB_NAME=gold_prod
JWT_SECRET=very_long_secure_random_string
REDIS_PASSWORD=secure_redis_password
SERVER_PORT=8080
```

## 4. SSL Certificates (Let's Encrypt)
Before starting Nginx, request certificates using Certbot.
(Requires your domain to point to this server IP)
```bash
sudo docker run -it --rm --name certbot \
    -v "/opt/gold_go/nginx/certbot/conf:/etc/letsencrypt" \
    -v "/opt/gold_go/nginx/certbot/www:/var/www/certbot" \
    certbot/certbot certonly --standalone -d your-domain.com
```

## 5. Build and Start Services
```bash
sudo docker-compose up -d --build
```

## 6. Run Database Migrations (Prisma)
With the containers running, run the migration inside the app container or locally if connected.
Recommended: Run inside container (since we didn't include Node/Prisma in the Go image to save space, you might need to run this from your CI/CD or local machine pointing to the remote DB, OR use a separate 'migrator' service).

**Option A (Recommended): Remote Migration (CI/CD or Local)**
```bash
# On your local machine:
export DATABASE_URL="postgresql://gold_user:secure_prod_password_123@your-server-ip:5432/gold_prod"
npx prisma migrate deploy
```
*(Make sure port 5432 is protected by firewall but accessible to your IP for this operation)*

**Option B: One-off Node Container**
```bash
docker run --rm -it \
    -v $(pwd):/app \
    --network gold_network \
    -e DATABASE_URL="postgresql://gold_user:secure_prod_password_123@postgres:5432/gold_prod" \
    node:18-alpine sh -c "cd /app && npm install prisma && npx prisma migrate deploy"
```

## 7. Verify Deployment
Check logs:
```bash
docker-compose logs -f app
```

Visit `https://your-domain.com/health` (if endpoint exists) or test API.

---

# Alternative: Deploy to Render.com (Free Tier)

**Render** is a cloud provider that offers a free tier for Web Services, PostgreSQL, and Redis, making it ideal for testing this project.

## 1. Push to GitHub
Ensure your code (including `render.yaml` and `Dockerfile`) is pushed to a public or private GitHub repository.

## 2. Sign up for Render
Go to [dashboard.render.com](https://dashboard.render.com/) and sign up/login with GitHub.

## 3. Create Blueprint
1.  Click **New +** -> **Blueprint**.
2.  Connect your GitHub account and select your `gold_go` repository.
3.  Render will automatically detect the `render.yaml` file.
4.  Click **Apply**.

## 4. Post-Deployment Steps
Render will provision the Go app, Postgres DB, and Redis instance automatically.
-   **Database**: It will be linked automatically via environment variables defined in `render.yaml`.
-   **Migrations**: 
    -   Render does not support running `npx` commands easily inside the Go Docker container (since it's Alpine). 
    -   **Solution**: Connect to the remote Render database from your local machine to run migrations.
    1.  Get the **External Connection String** from the Render Dashboard (Postgres service info).
    2.  Run locally:
        ```bash
        export DATABASE_URL="postgres://..."
        npx prisma migrate deploy
        ```

## 5. Notes on Free Tier
-   **Spin Down**: The free web service will "sleep" after 15 minutes of inactivity. The first request will take ~30s to wake it up.
-   **Database**: The free Postgres database expires after 90 days.
-   **Redis**: The free Redis instance is not persistent (data clears on restart).
