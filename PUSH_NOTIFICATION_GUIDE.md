# Push Notification Service Integration Guide

## Overview
Your backend has a OneSignal-based push notification service ready to send notifications to your Flutter app.

## Service Files
- **`pkg/notification/push.go`** - Push notification sender
- **`pkg/notification/mail.go`** - Email notification sender (also uses OneSignal)
- **`models/mail_type.go`** - Request/response models

## How It Works

### 1. OneSignal Setup (Already Done)
```
App ID: bb337783-7787-4f2f-8c9f-0dd994f2ccf5
API Key: os_v2_app_xmzxpa3xq5hs7de7bxmzj4wm6vbpekcdunwucx5mavjqsjn3rk7q7y7nijanhv7mhfocb5be4v5h7pnq6izbuibhn7fyoadocgwrdwa
```

### 2. Push Notification Format
The `SendOneSignalPush()` function accepts a `models.PushInput`:

```go
type PushInput struct {
    CustomerID        string                 // OneSignal external_id (set on device)
    Title             string                 // Notification title
    Message           string                 // Notification body
    NotificationType  string                 // e.g., "TOP-UP", "PAYMENT_SUCCESS"
    LinkedID          interface{}            // Context ID (transaction ID, etc.)
    ImageURL          string                 // Notification image
    LaunchURL         string                 // URL/deep link
    AppURL            string                 // App-specific deep link
    WebURL            string                 // Web-specific URL
    AndroidChannel    string                 // Android channel ID
    SmallIcon         string                 // Android icon
    IOSSound          string                 // iOS sound
    AndroidSound      string                 // Android sound
    TTL               int                    // TTL in seconds (default 3600)
    IsAndroid         *bool                  // Filter Android
    IsIos             *bool                  // Filter iOS
    IsAnyWeb          *bool                  // Filter web
    Data              map[string]interface{} // Custom payload
}
```

### 3. Backend Usage Example (in your handlers)

#### Example 1: Send Top-Up Notification
```go
import "github.com/919Umesh/stock_market_sim/pkg/notification"


respBody, err := notification.SendOneSignalPush(push)
if err != nil {
    slog.Error("failed to send push notification", "error", err)
    // Continue without failing the API call
}
```

## Flutter App Integration

### Step 1: Add OneSignal SDK to Flutter
```yaml
dependencies:
  onesignal_flutter: ^3.14.0
```

### Step 2: Initialize OneSignal in Flutter
```dart
import 'package:onesignal_flutter/onesignal_flutter.dart';

void main() {
  OneSignal.initialize("bb337783-7787-4f2f-8c9f-0dd994f2ccf5");
  OneSignal.Notifications.addForegroundWillDisplayListener((event) {
    // Handle notification in foreground
  });
  OneSignal.Notifications.addClickListener((OSNotificationClickEvent event) {
    // Handle notification tap
  });
  runApp(MyApp());
}
```

### Step 3: Set External ID for Device
Once user logs in, set their user ID as external_id:

```dart
String userId = authToken.userId; 
OneSignal.login(userId); 

```

### Step 4: Handle Incoming Notifications
```dart
OneSignal.Notifications.addForegroundWillDisplayListener((event) {
  OSNotification notification = event.notification;
  
  Map<String, dynamic> data = notification.additionalData ?? {};
  String notificationType = data['type'] ?? '';
  String linkedId = data['linked_id'] ?? '';
  
  // Handle different notification types
  switch(notificationType) {
    case 'TOP-UP':
      navigateTo('/wallet', {'transactionId': linkedId});
      break;
    case 'TRADE_EXECUTED':
      navigateTo('/portfolio', {'tradeId': linkedId});
      break;
    case 'ANNOUNCEMENT':
      showAnnouncementDialog(notification.title, notification.body);
      break;
  }
});


OneSignal.Notifications.addClickListener((OSNotificationClickEvent event) {
  OSNotification notification = event.notification;
  Map<String, dynamic> data = notification.additionalData ?? {};
  
  String launchUrl = notification.launchURL ?? '';
  if (launchUrl.isNotEmpty) {
    launchUrlInBrowser(launchUrl);
  }
});
```

## Security & Best Practices

1. **Never commit API keys** - Use environment variables:
   ```go
   oneSignalRestAPIKey = os.Getenv("ONE_SIGNAL_REST_API_KEY")
   oneSignalAppID = os.Getenv("ONE_SIGNAL_APP_ID")
   ```

2. **Validate user ownership** - Always verify the notification is being sent to the authenticated user:
   ```go
   // In handler
   userID := c.GetString("user_id")
   if userID != requestedUserID {
       c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
       return
   }
   ```

3. **Fail gracefully** - Don't fail the entire API if push notification fails:
   ```go
   _, err := notification.SendOneSignalPush(push)
   if err != nil {
       slog.Warn("push notification failed", "error", err)
       // Continue processing
   }
   ```

4. **Rate limit notifications** - Don't spam users with too many notifications

## Integrating into Your Handlers

### Wallet Handler (Top-Up)
Add after successful transaction in `internal/wallet/handler.go`:
```go
// After TopUp succeeds
wallet, transaction, userEmail, err := h.service.TopUp(userID, req.Amount, "", referenceID)
// ... error handling ...

// Send email
if userEmail != "" {
    mail.SendEmail(userEmail, "Top-Up", topupData)
}

// Send push notification
push := models.PushInput{
    CustomerID:       userID,
    Title:            "Wallet Credited",
    Message:          fmt.Sprintf("₹%.2f has been added to your wallet", transaction.Amount),
    NotificationType: "TOP-UP",
    LinkedID:         transaction.ID,
    Data: map[string]interface{}{
        "amount":       transaction.Amount,
        "balance":      wallet.FiatBalance,
        "reference_id": transaction.ReferenceID,
    },
}
notification.SendOneSignalPush(push) // Non-blocking
```

## Testing

1. **Test from backend:**
   ```bash
   curl -X POST https://your-api.com/wallet/topup \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"amount": 500}'
   ```

2. **Verify in OneSignal Dashboard:**
   - Go to OneSignal Dashboard
   - Check "Notifications" → "All Sent"
   - Look for your notification

3. **Check Flutter app:**
   - Foreground notification should appear if app is open
   - Background notification should trigger when app is closed

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Notifications not received | Verify `external_id` is set in Flutter app matching user ID |
| Wrong notification format | Check `CustomerID` matches OneSignal `external_id` |
| API errors | Verify API key and App ID in environment variables |
| Notifications not clickable | Ensure `LaunchURL`, `AppURL`, or deep links are set |

## Next Steps

1. Add push notifications to other handlers (trading, stock predictions, announcements)
2. Create notification templates in OneSignal Dashboard for better formatting
3. Implement notification history/preferences in your Flutter app
4. Add unsubscribe/preference management endpoints
