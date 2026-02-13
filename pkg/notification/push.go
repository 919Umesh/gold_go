package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/919Umesh/stock_market_sim/models"
)

const (
	ONE_SIGNAL_REST_API_KEY = "YOUR_ONESIGNAL_REST_API_KEY"
	ONE_SIGNAL_APP_ID       = "YOUR_ONESIGNAL_APP_ID"
)

var CHANNEL_MAP = map[string]string{
	"ORDERS":     "order_channel",
	"TICKETS":    "ticket_channel",
	"UPDATES":    "update_channel",
	"INVENTORY":  "inventory_channel",
	"PROMOTIONS": "promotion_channel",
	"PAYMENTS":   "payment_channel",
	"REVIEWS":    "review_channel",
	"QUERIES":    "query_channel",
	"ISSUES":     "issue_channel",
	"GENERAL":    "general_channel",
}



func SendPushNotification(
	customerID string,
	title string,
	message string,
	notificationType string,
	linkedID interface{},
) error {

	if customerID == "" || title == "" || message == "" {
		return errors.New("missing required fields")
	}

	if notificationType == "" {
		notificationType = "GENERAL"
	}

	channelID := CHANNEL_MAP[notificationType]

	reqBody := &models.SendPushNotification{
		AppID:                  ONE_SIGNAL_APP_ID,
		IncludeExternalUserIDs: []string{customerID},
		Headings:               map[string]string{"en": title},
		Contents:               map[string]string{"en": message},
		ExistingAndroidChannel: channelID,
		SmallIcon:              "ic_stat_onesignal",
		IOSSound:               "default",
		AndroidSound:           "default",
		TTL:                    3600,
		Data: map[string]interface{}{
			"linked_id": linkedID,
			"type":      notificationType,
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"https://onesignal.com/api/v1/notifications",
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+ONE_SIGNAL_REST_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("onesignal error %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Println("OneSignal response:", string(respBody))
	return nil
}
