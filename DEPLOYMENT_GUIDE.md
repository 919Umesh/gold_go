# Deployment Guide for Render

## 1. Prerequisites
- A GitHub repository containing this code.
- A Cloud Render account.

## 2. Setup on Render

1.  **Log in to Render** and go to your **Dashboard**.
2.  Click **New +** and select **Blueprint**.
3.  Connect your GitHub repository.
4.  Render will automatically detect the `render.yaml` file in the root directory.
5.  **Review the resources**:
    - `gold-go-app`: The Go web service.
    - `gold-db`: Managed PostgreSQL database.
    - `gold-redis`: Managed Redis instance.
6.  **Apply** the blueprint. Render will start deploying your services.

## 3. Environment Variables
The `render.yaml` automatically sets up connection strings for the database and Redis.
- `JWT_SECRET` is auto-generated.
- `PORT` defaults to `8080`.

## 4. Verification
Once deployed:
1.  Check the logs of the `gold-go-app` service. You should see "Server starting address :8080".
2.  Visit the service URL + `/health`. You should see `{"status":"healthy"}`.
