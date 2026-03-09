package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	ONESIGNAL_API_KEY = "os_v2_app_xmzxpa3xq5hs7de7bxmzj4wm6vbpekcdunwucx5mavjqsjn3rk7q7y7nijanhv7mhfocb5be4v5h7pnq6izbuibhn7fyoadocgwrdwa"
	OneSignalAppID    = "bb337783-7787-4f2f-8c9f-0dd994f2ccf5"
)

var templateIDs = map[string]string{
	"Top-Up": "5189dd66-5a29-4d31-a16d-03cc27ff5881",
}

// OneSignalEmailRequest is the request body for OneSignal email API
type OneSignalEmailRequest struct {
	AppID                     string      `json:"app_id"`
	TemplateID                string      `json:"template_id"`
	EmailFromName             string      `json:"email_from_name,omitempty"`
	EmailFromAddress          string      `json:"email_from_address,omitempty"`
	EmailSenderDomain         string      `json:"email_sender_domain,omitempty"`
	IncludeUnsubscribed       bool        `json:"include_unsubscribed,omitempty"`
	DisableEmailClickTracking bool        `json:"disable_email_click_tracking,omitempty"`
	Name                      string      `json:"name,omitempty"`
	EmailTo                   []string    `json:"include_email_tokens,omitempty"`
	CustomData                interface{} `json:"custom_data,omitempty"`
}

func SendEmail(email, emailType string, topUpData interface{}) error {
	if email == "" || emailType == "" || topUpData == nil {
		return errors.New("invalid input: email, emailType and topUpData are required")
	}

	templateID, ok := templateIDs[emailType]
	if !ok || templateID == "" {
		return fmt.Errorf("unknown template id for emailType: %s", emailType)
	}

	reqBody := &OneSignalEmailRequest{
		AppID:                     OneSignalAppID,
		TemplateID:                templateID,
		EmailFromName:             "LaganiPlus",
		EmailFromAddress:          "dale@umesh-shahi.com.np",
		EmailSenderDomain:         "mail.umesh-shahi.com.np",
		IncludeUnsubscribed:       true,
		DisableEmailClickTracking: false,
		Name:                      "Booking Confirmation Email",
		EmailTo:                   []string{email},
		CustomData:                topUpData,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://api.onesignal.com/notifications?c=email", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+ONESIGNAL_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("onesignal returned status %d: %s", resp.StatusCode, string(respBody))
	}

	slog.Info("OneSignal email response", "response", string(respBody))
	return nil
}
