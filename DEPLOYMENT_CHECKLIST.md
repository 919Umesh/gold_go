# Quick Deployment Checklist

Use this checklist to ensure a smooth deployment to Render.

## Pre-Deployment ✅

- [ ] Run deployment validation: `./scripts/deploy-check.sh`
- [ ] All tests passing: `make test`
- [ ] Build successful: `make build`
- [ ] Code formatted: `make fmt`
- [ ] Dependencies tidied: `make tidy`
- [ ] No sensitive data in code (API keys, passwords, etc.)
- [ ] `.env` file is in `.gitignore`
- [ ] All changes committed to Git
- [ ] On `main` branch (or your deployment branch)

## Render Setup 🚀

- [ ] GitHub account connected to Render
- [ ] Repository accessible to Render (public or authorized)
- [ ] Appwrite API key ready to set in Render dashboard

## Deployment Steps 📋

### 1. Push to GitHub
```bash
git add .
git commit -m "Deploy to Render with Redis support"
git push origin main
```

### 2. Create Blueprint on Render
- [ ] Go to https://dashboard.render.com
- [ ] Click "New +" → "Blueprint"
- [ ] Select repository: `919Umesh/gold_go`
- [ ] Confirm `render.yaml` is detected

### 3. Configure Environment Variables
- [ ] Navigate to `gold-go-backend` service
- [ ] Go to "Environment" tab
- [ ] Set `APPWRITE_API_KEY` to your actual API key
- [ ] Verify other variables are auto-configured

### 4. Deploy Services
- [ ] Click "Apply" to create services
- [ ] Wait for deployment to complete (5-10 minutes)
- [ ] Note your service URL (e.g., `https://gold-go-backend.onrender.com`)

## Post-Deployment Verification ✅

### 1. Health Check
```bash
curl https://your-service.onrender.com/health
```
Expected: `{"status":"healthy"}`

- [ ] Health endpoint returns 200 OK

### 2. Test Public Endpoints
```bash
# Get stocks list
curl https://your-service.onrender.com/api/v1/stocks

# Get market overview
curl https://your-service.onrender.com/api/v1/stocks/market-overview
```

- [ ] Public endpoints accessible
- [ ] Responses are valid JSON

### 3. Test Authentication
```bash
# Register a test user
curl -X POST https://your-service.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Test User","email":"test@example.com","phone":"1234567890","password":"Test@1234"}'

# Login
curl -X POST https://your-service.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test@1234"}'
```

- [ ] User registration works
- [ ] Login returns JWT token
- [ ] Token is valid

### 4. Test Protected Endpoints
```bash
# Use token from login
export TOKEN="your_jwt_token_here"

# Get wallet
curl https://your-service.onrender.com/api/v1/wallet \
  -H "Authorization: Bearer $TOKEN"

# Get portfolio
curl https://your-service.onrender.com/api/v1/trading/portfolio \
  -H "Authorization: Bearer $TOKEN"
```

- [ ] Protected endpoints require authentication
- [ ] Valid token grants access
- [ ] Invalid token returns 401

### 5. Check Logs
- [ ] Go to Render dashboard
- [ ] Select `gold-go-backend` service
- [ ] Click "Logs" tab
- [ ] Verify no errors in logs
- [ ] Check Redis connection successful

### 6. Monitor Metrics
- [ ] CPU usage is reasonable
- [ ] Memory usage is stable
- [ ] No error spikes
- [ ] Response times are acceptable

## Update Flutter App 📱

### 1. Update API Configuration
```dart
// lib/config/api_config.dart
class ApiConfig {
  static const String productionUrl = 'https://your-service.onrender.com/api/v1';
  static const String developmentUrl = 'http://localhost:8080/api/v1';
  
  static String get baseUrl {
    return kReleaseMode ? productionUrl : developmentUrl;
  }
}
```

- [ ] Production URL updated in Flutter app
- [ ] App tested with production backend
- [ ] All features working correctly

### 2. Test Flutter App
- [ ] Registration flow works
- [ ] Login flow works
- [ ] Stock listing loads
- [ ] Trading features work
- [ ] Wallet operations work
- [ ] Error handling works

## Production Readiness 🎯

### Security
- [ ] HTTPS enabled (automatic on Render)
- [ ] API keys stored as environment variables
- [ ] JWT secret is secure and auto-generated
- [ ] No sensitive data in logs
- [ ] CORS configured correctly

### Performance
- [ ] Redis caching working
- [ ] Response times acceptable
- [ ] No memory leaks
- [ ] Worker pool functioning

### Monitoring
- [ ] Health checks configured
- [ ] Logs accessible
- [ ] Metrics visible
- [ ] Alerts set up (optional)

### Scaling (Optional)
- [ ] Consider upgrading to paid plan for production
- [ ] No spin-down (always-on service)
- [ ] Custom domain configured (optional)
- [ ] CDN for static assets (if applicable)

## Troubleshooting 🔧

If something goes wrong:

### Build Fails
- [ ] Check build logs in Render dashboard
- [ ] Verify `go.mod` and `go.sum` are committed
- [ ] Test build locally: `make build`
- [ ] Check Dockerfile syntax

### Service Won't Start
- [ ] Verify all environment variables are set
- [ ] Check Appwrite credentials are correct
- [ ] Review logs for error messages
- [ ] Ensure health endpoint is accessible

### Redis Connection Issues
- [ ] Verify Redis service is running
- [ ] Check environment variables are linked
- [ ] Ensure both services in same region
- [ ] Review Redis logs

### Appwrite Connection Issues
- [ ] Verify API key is correct
- [ ] Check project ID and database ID
- [ ] Ensure endpoint URL is correct
- [ ] Test Appwrite connection locally

## Rollback Plan 🔄

If deployment fails:

1. **Render Auto-Rollback**
   - Render automatically rolls back if health check fails
   - Previous version remains active

2. **Manual Rollback**
   - Go to Render dashboard
   - Select service
   - Click "Rollback" to previous deployment

3. **Git Rollback**
   ```bash
   git revert HEAD
   git push origin main
   ```
   - Render will auto-deploy the reverted version

## Success Criteria ✨

Your deployment is successful when:

- ✅ All health checks passing
- ✅ All API endpoints accessible
- ✅ Authentication working
- ✅ Database operations successful
- ✅ Redis caching functional
- ✅ No errors in logs
- ✅ Flutter app connects successfully
- ✅ All features working as expected

## Next Steps After Successful Deployment 🎉

1. **Monitor for 24 hours**
   - Watch logs for any issues
   - Monitor performance metrics
   - Check error rates

2. **Consider Upgrades**
   - Upgrade to paid plan for production
   - Add custom domain
   - Set up monitoring alerts
   - Configure backup strategy

3. **Documentation**
   - Document production URL
   - Update team on deployment
   - Create runbook for common issues

4. **Continuous Improvement**
   - Set up automated testing
   - Implement monitoring dashboards
   - Plan for scaling
   - Regular security updates

---

**Congratulations! 🎊 Your backend is now live on Render!**

For detailed troubleshooting, see [DEPLOYMENT.md](./DEPLOYMENT.md)
