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

type PushInput struct {
	CustomerID       string                 `json:"customer_id"`                 // OneSignal external_id you set on the device
	Title            string                 `json:"title"`                       // Notification title
	Message          string                 `json:"message"`                     // Notification body
	NotificationType string                 `json:"notification_type,omitempty"` // e.g., "TOP-UP"
	LinkedID         interface{}            `json:"linked_id,omitempty"`         // any ID to relate context
	ImageURL         string                 `json:"image_url,omitempty"`         // dashboard "Image"
	LaunchURL        string                 `json:"launch_url,omitempty"`        // dashboard "Launch URL" (use url/app_url/web_url)
	AppURL           string                 `json:"app_url,omitempty"`           // optional deep link for native
	WebURL           string                 `json:"web_url,omitempty"`           // optional web specific URL
	AndroidChannel   string                 `json:"android_channel,omitempty"`   // OneSignal existing Android channel ID
	SmallIcon        string                 `json:"small_icon,omitempty"`        // Android small icon
	IOSSound         string                 `json:"ios_sound,omitempty"`         // iOS sound
	AndroidSound     string                 `json:"android_sound,omitempty"`     // Android sound
	TTL              int                    `json:"ttl,omitempty"`               // seconds
	IsAndroid        *bool                  `json:"is_android,omitempty"`        // platform filters
	IsIos            *bool                  `json:"is_ios,omitempty"`
	IsAnyWeb         *bool                  `json:"is_any_web,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"` // custom payload ("data"/"custom_data")
}

// OneSignal request for Push API (?c=push)
type OneSignalPush struct {
	AppID                    string                 `json:"app_id"`
	IncludeAliases           map[string][]string    `json:"include_aliases,omitempty"` // { "external_id": ["..."] }
	TargetChannel            string                 `json:"target_channel,omitempty"`  // "push" when using include_aliases
	Headings                 map[string]string      `json:"headings,omitempty"`        // {"en": "..."}
	Contents                 map[string]string      `json:"contents,omitempty"`        // {"en": "..."}
	Data                     map[string]interface{} `json:"data,omitempty"`            // delivered to app
	CustomData               map[string]interface{} `json:"custom_data,omitempty"`     // also supported by API
	URL                      string                 `json:"url,omitempty"`             // generic Launch URL
	AppURL                   string                 `json:"app_url,omitempty"`
	WebURL                   string                 `json:"web_url,omitempty"`
	BigPicture               string                 `json:"big_picture,omitempty"`      // Android big picture
	ChromeWebImage           string                 `json:"chrome_web_image,omitempty"` // Web image
	IOSAttachments           map[string]string      `json:"ios_attachments,omitempty"`  // {"id": "<image_url>"}
	ExistingAndroidChannelID string                 `json:"existing_android_channel_id,omitempty"`
	SmallIcon                string                 `json:"small_icon,omitempty"`
	IOSSound                 string                 `json:"ios_sound,omitempty"`
	AndroidSound             string                 `json:"android_sound,omitempty"`
	TTL                      int                    `json:"ttl,omitempty"`
	IsAndroid                *bool                  `json:"isAndroid,omitempty"`
	IsIos                    *bool                  `json:"isIos,omitempty"`
	IsAnyWeb                 *bool                  `json:"isAnyWeb,omitempty"`
}
