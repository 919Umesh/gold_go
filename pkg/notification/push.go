package notification

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

var (
	oneSignalRestAPIKey = "os_v2_app_xmzxpa3xq5hs7de7bxmzj4wm6vbpekcdunwucx5mavjqsjn3rk7q7y7nijanhv7mhfocb5be4v5h7pnq6izbuibhn7fyoadocgwrdwa"
	oneSignalAppID      = "bb337783-7787-4f2f-8c9f-0dd994f2ccf5"
	httpClient          = &http.Client{Timeout: 10 * time.Second}
)

const oneSignalPushEndpoint = "https://api.onesignal.com/notifications?c=push"

func SendOneSignalPush(in models.PushInput) ([]byte, error) {
	if oneSignalAppID == "" || oneSignalRestAPIKey == "" {
		return nil, errors.New("missing ONE_SIGNAL_APP_ID or ONE_SIGNAL_REST_API_KEY")
	}

	if in.CustomerID == "" {
		return nil, errors.New("customer_id is required")
	}

	if in.Title == "" && in.Message == "" {
		return nil, errors.New("at least one of title or message required")
	}

	reqBody := &models.OneSignalPush{
		AppID: oneSignalAppID,
		IncludeAliases: map[string][]string{
			"external_id": {in.CustomerID},
		},
		TargetChannel:            "push",
		Headings:                 map[string]string{"en": in.Title},
		Contents:                 map[string]string{"en": in.Message},
		URL:                      in.LaunchURL,
		AppURL:                   in.AppURL,
		WebURL:                   in.WebURL,
		ExistingAndroidChannelID: in.AndroidChannel,
		SmallIcon:                in.SmallIcon,
		IOSSound:                 nonEmpty(in.IOSSound, "default"),
		AndroidSound:             nonEmpty(in.AndroidSound, "default"),
		TTL:                      firstNonZero(in.TTL, 3600),
		IsAndroid:                in.IsAndroid,
		IsIos:                    in.IsIos,
		IsAnyWeb:                 in.IsAnyWeb,
	}

	if in.ImageURL != "" {
		reqBody.BigPicture = in.ImageURL
		reqBody.ChromeWebImage = in.ImageURL
		reqBody.IOSAttachments = map[string]string{
			"id": in.ImageURL,
		}
	}

	reqBody.Data = map[string]interface{}{
		"type":      in.NotificationType,
		"linked_id": in.LinkedID,
	}
	for k, v := range in.Data {
		reqBody.Data[k] = v
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		oneSignalPushEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+oneSignalRestAPIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("onesignal request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("onesignal error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
func nonEmpty(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func firstNonZero(value int, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
