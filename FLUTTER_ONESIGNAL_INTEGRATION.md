# Flutter OneSignal Push Notification Integration Guide

## Overview
This guide shows how to integrate OneSignal push notifications in your Flutter app to receive wallet top-up notifications and other events from your backend.

## Prerequisites
- Flutter project setup
- OneSignal account (free tier available)
- Your backend API base URL
- OneSignal App ID: `bb337783-7787-4f2f-8c9f-0dd994f2ccf5`

---

## Step 1: Install OneSignal Flutter Package

### Update `pubspec.yaml`
```yaml
dependencies:
  flutter:
    sdk: flutter
  onesignal_flutter: ^4.0.0
  intl: ^0.19.0
```

### Install Dependencies
```bash
flutter pub get
```

---

## Step 2: Platform-Specific Setup

### iOS Setup

#### Update `ios/Podfile`
```ruby
# Uncomment the platform at the top
platform :ios, '12.0'

post_install do |installer|
  installer.pods_project.targets.each do |target|
    flutter_additional_ios_build_settings(target)
    target.build_configurations.each do |config|
      config.build_settings['GCC_PREPROCESSOR_DEFINITIONS'] ||= [
        '$(inherited)',
        'PERMISSION_NOTIFICATIONS=1',
      ]
    end
  end
end
```

#### Update `ios/Runner/Info.plist`
Add OneSignal required keys:
```xml
<key>NSUserNotificationAlertType</key>
<string>alert</string>
<key>OneSignalAppId</key>
<string>bb337783-7787-4f2f-8c9f-0dd994f2ccf5</string>
```

#### Register for Remote Notifications (`ios/Runner/GeneratedPluginRegistrant.m`)
Ensure OneSignal is initialized before other plugins.

### Android Setup

#### Update `android/build.gradle`
```gradle
buildscript {
    repositories {
        google()
        mavenCentral()
    }

    dependencies {
        classpath 'com.android.tools.build:gradle:7.3.1'
        classpath 'com.google.gms:google-services:4.3.15'
    }
}
```

#### Update `android/app/build.gradle`
```gradle
apply plugin: 'com.android.application'
apply plugin: 'com.google.gms.google-services'

dependencies {
    // OneSignal SDK
    implementation 'com.onesignal:OneSignal:[5.0.0, 5.99.99]'
    
    // Google Play Services
    implementation 'com.google.android.gms:play-services-location:21.0.1'
}

android {
    compileSdkVersion 34
    
    defaultConfig {
        applicationId "com.your_company.app_name"
        minSdkVersion 21
        targetSdkVersion 34
    }
}
```

#### Update `android/app/src/main/AndroidManifest.xml`
```xml
<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    
    <uses-permission android:name="android.permission.INTERNET" />
    <uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
    
    <application>
        <!-- Your app configuration -->
    </application>
</manifest>
```

#### Update `android/app/src/main/java/MainActivity.kt`
```kotlin
package com.example.stock_market_app

import android.os.Bundle
import io.flutter.embedding.android.FlutterActivity

class MainActivity: FlutterActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }
}
```

---

## Step 3: Create OneSignal Service Class

Create `lib/services/onesignal_service.dart`:

```dart
import 'package:onesignal_flutter/onesignal_flutter.dart';
import 'package:flutter/foundation.dart';
import 'dart:developer' as developer;

class OneSignalService {
  static const String appId = 'bb337783-7787-4f2f-8c9f-0dd994f2ccf5';

  static Future<void> initialize() async {
    try {
      // Initialize OneSignal
      OneSignal.initialize(appId);

      // Request user permission for notifications (iOS)
      OneSignal.Notifications.requestPermission(true);

      // Set logging level for debugging
      OneSignal.Debug.setLogLevel(OSLogLevel.verbose);

      // Set up notification handlers
      _setupNotificationHandlers();

      developer.log('OneSignal initialized successfully');
    } catch (e) {
      developer.log('OneSignal initialization error: $e');
    }
  }

  /// Setup notification handlers for foreground and click events
  static void _setupNotificationHandlers() {
    // Handle notification while app is in foreground
    OneSignal.Notifications.addForegroundWillDisplayListener((event) {
      developer.log('Notification received in foreground: ${event.notification.title}');

      // You can optionally display a custom dialog or handle it
      // by default OneSignal will display the notification
      event.notification.display();
    });

    // Handle notification clicks
    OneSignal.Notifications.addClickListener((OSNotificationClickEvent event) {
      developer.log('Notification clicked: ${event.notification.title}');
      _handleNotificationClick(event.notification);
    });
  }

  /// Handle notification clicks and route accordingly
  static void _handleNotificationClick(OSNotification notification) {
    final additionalData = notification.additionalData;
    
    if (additionalData == null) return;

    final String notificationType = additionalData['type'] ?? '';
    final String linkedId = additionalData['linked_id']?.toString() ?? '';

    developer.log('Notification type: $notificationType, LinkedID: $linkedId');

    // Route based on notification type
    switch (notificationType) {
      case 'TOP-UP':
        _handleTopUpNotification(additionalData, linkedId);
        break;
      case 'TRADE_EXECUTED':
        _handleTradeNotification(additionalData, linkedId);
        break;
      case 'ANNOUNCEMENT':
        _handleAnnouncementNotification(additionalData);
        break;
      default:
        developer.log('Unknown notification type: $notificationType');
    }
  }

  static void _handleTopUpNotification(Map<String, dynamic> data, String transactionId) {
    // Navigate to wallet screen with transaction details
    // Example: navigatorKey.currentState?.pushNamed('/wallet', arguments: {'transactionId': transactionId});
    
    double? amount = double.tryParse(data['amount']?.toString() ?? '0');
    double? balance = double.tryParse(data['new_balance']?.toString() ?? '0');
    
    developer.log('Top-Up: ₹$amount, New Balance: ₹$balance');
    
    // You can show a snackbar or navigate here
    // ScaffoldMessenger.of(context).showSnackBar(
    //   SnackBar(content: Text('Wallet credited: ₹$amount')),
    // );
  }

  static void _handleTradeNotification(Map<String, dynamic> data, String tradeId) {
    // Navigate to portfolio or trade details
    String? stock = data['stock']?.toString();
    int? quantity = int.tryParse(data['quantity']?.toString() ?? '0');
    double? price = double.tryParse(data['price']?.toString() ?? '0');
    
    developer.log('Trade executed: $stock x$quantity @ ₹$price');
  }

  static void _handleAnnouncementNotification(Map<String, dynamic> data) {
    // Show announcement dialog or navigate to announcements page
    developer.log('Announcement received');
  }

  /// Set user ID (call after login)
  static Future<void> setUserID(String userId) async {
    try {
      OneSignal.login(userId);
      developer.log('User ID set for OneSignal: $userId');
    } catch (e) {
      developer.log('Error setting user ID: $e');
    }
  }

  /// Get device state
  static Future<String?> getDeviceState() async {
    try {
      final state = await OneSignal.User.getOnesignalId();
      developer.log('Device OneSignal ID: $state');
      return state;
    } catch (e) {
      developer.log('Error getting device state: $e');
      return null;
    }
  }

  /// Logout user (call on logout)
  static Future<void> logoutUser() async {
    try {
      OneSignal.logout();
      developer.log('User logged out from OneSignal');
    } catch (e) {
      developer.log('Error logging out: $e');
    }
  }

  /// Add custom tag (optional)
  static Future<void> addTag(String key, String value) async {
    try {
      OneSignal.User.addTag(key, value);
      developer.log('Added tag: $key = $value');
    } catch (e) {
      developer.log('Error adding tag: $e');
    }
  }

  /// Remove tag (optional)
  static Future<void> removeTag(String key) async {
    try {
      OneSignal.User.removeTag(key);
      developer.log('Removed tag: $key');
    } catch (e) {
      developer.log('Error removing tag: $e');
    }
  }
}
```

---

## Step 4: Initialize OneSignal in Main App

Update `lib/main.dart`:

```dart
import 'package:flutter/material.dart';
import 'services/onesignal_service.dart';
import 'screens/home_screen.dart';
import 'screens/wallet_screen.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize OneSignal
  await OneSignalService.initialize();
  
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Stock Market Simulator',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const HomeScreen(),
      routes: {
        '/wallet': (context) => const WalletScreen(),
        '/portfolio': (context) => const PortfolioScreen(),
      },
      navigatorObservers: [
        // Optional: Add observer to track navigation
      ],
    );
  }
}
```

---

## Step 5: Update Login Screen

Update `lib/screens/login_screen.dart`:

```dart
import 'package:flutter/material.dart';
import 'package:your_app/services/onesignal_service.dart';
import 'package:your_app/services/auth_service.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({Key? key}) : super(key: key);

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  bool isLoading = false;

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (emailController.text.isEmpty || passwordController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter email and password')),
      );
      return;
    }

    setState(() => isLoading = true);

    try {
      // Call your backend login API
      final loginResponse = await AuthService.login(
        email: emailController.text,
        password: passwordController.text,
      );

      // Set user ID in OneSignal
      await OneSignalService.setUserID(loginResponse.userId);

      // Add user properties as tags (optional)
      await OneSignalService.addTag('email', emailController.text);
      await OneSignalService.addTag('kyc_status', loginResponse.kycStatus);

      // Navigate to home screen
      if (mounted) {
        Navigator.of(context).pushReplacementNamed('/home');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Login failed: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Login')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            TextField(
              controller: emailController,
              decoration: const InputDecoration(
                labelText: 'Email',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: passwordController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Password',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 24),
            isLoading
                ? const CircularProgressIndicator()
                : ElevatedButton(
                    onPressed: _handleLogin,
                    child: const Text('Login'),
                  ),
          ],
        ),
      ),
    );
  }
}
```

---

## Step 6: Create Wallet Service for Top-Up

Create `lib/services/wallet_service.dart`:

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;

class WalletService {
  static const String baseUrl = 'https://your-api.com/api';
  static String? authToken;

  /// Set auth token after login
  static void setAuthToken(String token) {
    authToken = token;
  }

  /// Perform wallet top-up
  static Future<TopUpResponse> topUp({
    required double amount,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/wallet/topup'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $authToken',
        },
        body: jsonEncode({
          'amount': amount,
        }),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        developer.log('Top-up successful: ${jsonEncode(data)}');
        return TopUpResponse.fromJson(data);
      } else if (response.statusCode == 401) {
        throw 'Unauthorized. Please login again.';
      } else {
        throw 'Top-up failed: ${response.body}';
      }
    } catch (e) {
      developer.log('Top-up error: $e');
      rethrow;
    }
  }

  /// Get wallet balance
  static Future<WalletResponse> getWallet() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/wallet'),
        headers: {
          'Authorization': 'Bearer $authToken',
        },
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return WalletResponse.fromJson(data);
      } else {
        throw 'Failed to fetch wallet';
      }
    } catch (e) {
      developer.log('Get wallet error: $e');
      rethrow;
    }
  }
}

// Models
class TopUpResponse {
  final String message;
  final WalletData wallet;
  final TransactionData transaction;

  TopUpResponse({
    required this.message,
    required this.wallet,
    required this.transaction,
  });

  factory TopUpResponse.fromJson(Map<String, dynamic> json) {
    return TopUpResponse(
      message: json['message'] ?? '',
      wallet: WalletData.fromJson(json['wallet'] ?? {}),
      transaction: TransactionData.fromJson(json['transaction'] ?? {}),
    );
  }
}

class WalletResponse {
  final WalletData wallet;

  WalletResponse({required this.wallet});

  factory WalletResponse.fromJson(Map<String, dynamic> json) {
    return WalletResponse(
      wallet: WalletData.fromJson(json['wallet'] ?? {}),
    );
  }
}

class WalletData {
  final String id;
  final double fiatBalance;
  final double balance;
  final double totalInvested;
  final double totalProfitLoss;

  WalletData({
    required this.id,
    required this.fiatBalance,
    required this.balance,
    required this.totalInvested,
    required this.totalProfitLoss,
  });

  factory WalletData.fromJson(Map<String, dynamic> json) {
    return WalletData(
      id: json['\$id'] ?? '',
      fiatBalance: (json['fiat_balance'] as num?)?.toDouble() ?? 0.0,
      balance: (json['balance'] as num?)?.toDouble() ?? 0.0,
      totalInvested: (json['total_invested'] as num?)?.toDouble() ?? 0.0,
      totalProfitLoss: (json['total_profit_loss'] as num?)?.toDouble() ?? 0.0,
    );
  }
}

class TransactionData {
  final String id;
  final double amount;
  final String status;
  final String referenceId;
  final DateTime createdAt;

  TransactionData({
    required this.id,
    required this.amount,
    required this.status,
    required this.referenceId,
    required this.createdAt,
  });

  factory TransactionData.fromJson(Map<String, dynamic> json) {
    return TransactionData(
      id: json['\$id'] ?? '',
      amount: (json['amount'] as num?)?.toDouble() ?? 0.0,
      status: json['status'] ?? '',
      referenceId: json['reference_id'] ?? '',
      createdAt: DateTime.tryParse(json['\$createdAt'] ?? '') ?? DateTime.now(),
    );
  }
}
```

---

## Step 7: Create Wallet Screen

Create `lib/screens/wallet_screen.dart`:

```dart
import 'package:flutter/material.dart';
import 'package:your_app/services/wallet_service.dart';
import 'package:intl/intl.dart';

class WalletScreen extends StatefulWidget {
  const WalletScreen({Key? key}) : super(key: key);

  @override
  State<WalletScreen> createState() => _WalletScreenState();
}

class _WalletScreenState extends State<WalletScreen> {
  late Future<WalletResponse> walletFuture;
  final topUpAmountController = TextEditingController();
  bool isTopUpLoading = false;

  @override
  void initState() {
    super.initState();
    walletFuture = WalletService.getWallet();
  }

  @override
  void dispose() {
    topUpAmountController.dispose();
    super.dispose();
  }

  Future<void> _handleTopUp() async {
    final amount = double.tryParse(topUpAmountController.text);
    if (amount == null || amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter a valid amount')),
      );
      return;
    }

    setState(() => isTopUpLoading = true);

    try {
      final response = await WalletService.topUp(amount: amount);
      
      if (mounted) {
        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${response.message}\nNew Balance: ₹${response.wallet.fiatBalance}'),
            backgroundColor: Colors.green,
          ),
        );

        // Refresh wallet data
        setState(() {
          walletFuture = WalletService.getWallet();
          topUpAmountController.clear();
        });

        // Close dialog if open
        if (Navigator.of(context).canPop()) {
          Navigator.pop(context);
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Top-up failed: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => isTopUpLoading = false);
      }
    }
  }

  void _showTopUpDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add Money to Wallet'),
        content: TextField(
          controller: topUpAmountController,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(
            hintText: 'Enter amount',
            prefixText: '₹ ',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: isTopUpLoading ? null : _handleTopUp,
            child: isTopUpLoading
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Top-up'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Wallet'),
        centerTitle: true,
      ),
      body: FutureBuilder<WalletResponse>(
        future: walletFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }

          if (snapshot.hasError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Error: ${snapshot.error}'),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => setState(() => walletFuture = WalletService.getWallet()),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          final wallet = snapshot.data?.wallet;
          if (wallet == null) {
            return const Center(child: Text('No wallet data'));
          }

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                // Wallet Card
                Card(
                  elevation: 4,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(24.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Wallet Balance',
                          style: TextStyle(
                            fontSize: 14,
                            color: Colors.grey,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '₹${wallet.fiatBalance.toStringAsFixed(2)}',
                          style: const TextStyle(
                            fontSize: 32,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 16),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Total Invested',
                                  style: TextStyle(fontSize: 12, color: Colors.grey),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  '₹${wallet.totalInvested.toStringAsFixed(2)}',
                                  style: const TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ],
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Profit/Loss',
                                  style: TextStyle(fontSize: 12, color: Colors.grey),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  '₹${wallet.totalProfitLoss.toStringAsFixed(2)}',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: wallet.totalProfitLoss >= 0
                                        ? Colors.green
                                        : Colors.red,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                // Top-up Button
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: _showTopUpDialog,
                    child: const Text(
                      'Add Money',
                      style: TextStyle(fontSize: 16),
                    ),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
```

---

## Step 8: Update Logout Functionality

Update your logout handler:

```dart
import 'package:your_app/services/onesignal_service.dart';

Future<void> handleLogout() async {
  // Logout from OneSignal
  await OneSignalService.logoutUser();
  
  // Clear any stored tokens
  // await AuthService.clearToken();
  
  // Navigate to login screen
  // Navigator.of(context).pushReplacementNamed('/login');
}
```

---

## Step 9: Testing

### Test Top-Up Flow
1. **Login** to your Flutter app
2. **Go to Wallet** screen
3. **Click "Add Money"** button
4. **Enter amount** (e.g., 500)
5. **Submit** the top-up request

### Expected Results
- ✅ Backend processes the top-up
- ✅ Email notification sent (if email available)
- ✅ **Push notification appears on device**
  - **Foreground**: Notification banner appears
  - **Background**: Notification in notification center
  - **Locked screen**: Notification shows on lock screen
- ✅ Wallet balance updates
- ✅ Transaction appears in history

### Verify in OneSignal Dashboard
1. Go to [OneSignal Dashboard](https://onesignal.com/apps)
2. Select your app
3. Navigate to **Notifications → All Sent**
4. Look for your top-up notification
5. Check delivery status

---

## Step 10: Advanced Features (Optional)

### A. Send Deep Links in Notifications

```dart
// In wallet handler - send deep link
push := models.PushInput{
    CustomerID:       userID,
    Title:            "Wallet Top-Up",
    Message:          "Successful",
    LaunchURL:        "https://your-app.com/wallet/topup",
    AppURL:           "lagani://wallet/topup",
    Data: map[string]interface{}{
        // ...
    },
}
```

### B. Add Custom Sounds

```go
// In backend
push := models.PushInput{
    // ...
    IOSSound:     "notification_sound.wav",
    AndroidSound: "notification_sound",
}
```

### C. Send to Specific User Segments

```go
// Tag users by KYC status
await OneSignalService.addTag('kyc_status', 'verified');

// Then target in backend - send to "verified" segment only
```

### D. Track Notification Analytics

```dart
// In Flutter - track if user opened notification from specific action
OneSignal.Notifications.addClickListener((OSNotificationClickEvent event) {
  final action = event.notification.actionId;
  developer.log('User clicked action: $action');
});
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| No notifications received | 1. Check user ID is set correctly after login<br>2. Verify App ID in code<br>3. Check device is connected to internet<br>4. Restart app and trigger top-up |
| Notifications only in foreground | Check app permissions on device Settings |
| Wrong notification format | Verify `CustomerID` matches `external_id` in OneSignal |
| App crashes on notification | Check notification handler doesn't have null reference errors |
| Notifications delayed | OneSignal queues notifications - check delivery in dashboard |

---

## Quick Reference

### Key Methods
```dart
// Initialize
OneSignalService.initialize();

// After login
OneSignalService.setUserID(userId);

// Add user properties
OneSignalService.addTag('email', userEmail);

// On logout
OneSignalService.logoutUser();

// Get device ID
String? deviceId = await OneSignalService.getDeviceState();
```

### Notification Types Handled
- `TOP-UP`: Navigate to wallet
- `TRADE_EXECUTED`: Navigate to portfolio
- `ANNOUNCEMENT`: Show dialog

### Data Passed from Backend
```json
{
  "type": "TOP-UP",
  "linked_id": "transaction-id",
  "amount": 500.00,
  "new_balance": 5500.00,
  "reference_id": "topup_xxxx"
}
```

---

## Next Steps
1. ✅ Implement notification handlers for other events (trading, announcements)
2. ✅ Add notification history screen
3. ✅ Implement notification preferences
4. ✅ Add unsubscribe option
5. ✅ Track notification analytics

