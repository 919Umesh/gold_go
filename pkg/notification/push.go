package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	oneSignalRestAPIKey = "os_v2_app_xmzxpa3xq5hs7de7bxmzj4wm6vbpekcdunwucx5mavjqsjn3rk7q7y7nijanhv7mhfocb5be4v5h7pnq6izbuibhn7fyoadocgwrdwa"
	oneSignalAppID      = "bb337783-7787-4f2f-8c9f-0dd994f2ccf5"
	httpClient          = &http.Client{Timeout: 10 * time.Second}
)

const oneSignalPushEndpoint = "https://api.onesignal.com/notifications?c=push"

// PushInput contains the input for sending a push notification
type PushInput struct {
	CustomerID       string
	Title            string
	Message          string
	LaunchURL        string
	AppURL           string
	WebURL           string
	AndroidChannel   string
	SmallIcon        string
	IOSSound         string
	AndroidSound     string
	TTL              int
	IsAndroid        bool
	IsIos            bool
	IsAnyWeb         bool
	ImageURL         string
	NotificationType string
	LinkedID         string
	Data             map[string]interface{}
}

// OneSignalPush is the request body for OneSignal push API
type OneSignalPush struct {
	AppID                    string                 `json:"app_id"`
	IncludeAliases           map[string][]string    `json:"include_aliases,omitempty"`
	TargetChannel            string                 `json:"target_channel"`
	Headings                 map[string]string      `json:"headings,omitempty"`
	Contents                 map[string]string      `json:"contents,omitempty"`
	URL                      string                 `json:"url,omitempty"`
	AppURL                   string                 `json:"app_url,omitempty"`
	WebURL                   string                 `json:"web_url,omitempty"`
	ExistingAndroidChannelID string                 `json:"existing_android_channel_id,omitempty"`
	SmallIcon                string                 `json:"small_icon,omitempty"`
	IOSSound                 string                 `json:"ios_sound,omitempty"`
	AndroidSound             string                 `json:"android_sound,omitempty"`
	TTL                      int                    `json:"ttl,omitempty"`
	IsAndroid                bool                   `json:"isAndroid,omitempty"`
	IsIos                    bool                   `json:"isIos,omitempty"`
	IsAnyWeb                 bool                   `json:"isAnyWeb,omitempty"`
	BigPicture               string                 `json:"big_picture,omitempty"`
	ChromeWebImage           string                 `json:"chrome_web_image,omitempty"`
	IOSAttachments           map[string]string      `json:"ios_attachments,omitempty"`
	Data                     map[string]interface{} `json:"data,omitempty"`
}

func SendOneSignalPush(in PushInput) ([]byte, error) {
	if oneSignalAppID == "" || oneSignalRestAPIKey == "" {
		return nil, errors.New("missing ONE_SIGNAL_APP_ID or ONE_SIGNAL_REST_API_KEY")
	}

	if in.CustomerID == "" {
		return nil, errors.New("customer_id is required")
	}

	if in.Title == "" && in.Message == "" {
		return nil, errors.New("at least one of title or message required")
	}

	reqBody := &OneSignalPush{
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
