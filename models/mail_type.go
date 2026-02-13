package models

type OneSignalEmailRequest struct {
	AppID                     string      `json:"app_id"`
	TemplateID                string      `json:"template_id,omitempty"`
	EmailFromName             string      `json:"email_from_name,omitempty"`
	EmailFromAddress          string      `json:"email_from_address,omitempty"`
	EmailSenderDomain         string      `json:"email_sender_domain,omitempty"`
	IncludeUnsubscribed       bool        `json:"include_unsubscribed,omitempty"`
	DisableEmailClickTracking bool        `json:"disable_email_click_tracking,omitempty"`
	Name                      string      `json:"name,omitempty"`
	EmailTo                   []string    `json:"email_to,omitempty"`
	CustomData                interface{} `json:"custom_data,omitempty"`
}

type SendPushNotification struct {
	AppID                   string                 `json:"app_id"`
	IncludeExternalUserIDs  []string               `json:"include_external_user_ids,omitempty"`
	Headings                map[string]string      `json:"headings,omitempty"`
	Contents                map[string]string      `json:"contents,omitempty"`
	Data                    map[string]interface{} `json:"data,omitempty"`
	ExistingAndroidChannel  string                 `json:"existing_android_channel_id,omitempty"`
	SmallIcon               string                 `json:"small_icon,omitempty"`
	IOSSound                string                 `json:"ios_sound,omitempty"`
	AndroidSound            string                 `json:"android_sound,omitempty"`
	TTL                     int                    `json:"ttl,omitempty"`
}