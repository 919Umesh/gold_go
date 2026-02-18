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
	CustomerID       string                 `json:"customer_id"`                 
	Title            string                 `json:"title"`                      
	Message          string                 `json:"message"`                     
	NotificationType string                 `json:"notification_type,omitempty"` 
	LinkedID         interface{}            `json:"linked_id,omitempty"`         
	ImageURL         string                 `json:"image_url,omitempty"`        
	LaunchURL        string                 `json:"launch_url,omitempty"`        
	AppURL           string                 `json:"app_url,omitempty"`           
	WebURL           string                 `json:"web_url,omitempty"`           
	AndroidChannel   string                 `json:"android_channel,omitempty"`   
	SmallIcon        string                 `json:"small_icon,omitempty"`        
	IOSSound         string                 `json:"ios_sound,omitempty"`        
	AndroidSound     string                 `json:"android_sound,omitempty"`     
	TTL              int                    `json:"ttl,omitempty"`               
	IsAndroid        *bool                  `json:"is_android,omitempty"`        
	IsIos            *bool                  `json:"is_ios,omitempty"`
	IsAnyWeb         *bool                  `json:"is_any_web,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"` 
}


type OneSignalPush struct {
	AppID                    string                 `json:"app_id"`
	IncludeAliases           map[string][]string    `json:"include_aliases,omitempty"` 
	TargetChannel            string                 `json:"target_channel,omitempty"`  
	Headings                 map[string]string      `json:"headings,omitempty"`        
	Contents                 map[string]string      `json:"contents,omitempty"`        
	Data                     map[string]interface{} `json:"data,omitempty"`           
	CustomData               map[string]interface{} `json:"custom_data,omitempty"`     
	URL                      string                 `json:"url,omitempty"`            
	AppURL                   string                 `json:"app_url,omitempty"`
	WebURL                   string                 `json:"web_url,omitempty"`
	BigPicture               string                 `json:"big_picture,omitempty"`    
	ChromeWebImage           string                 `json:"chrome_web_image,omitempty"` 
	IOSAttachments           map[string]string      `json:"ios_attachments,omitempty"` 
	ExistingAndroidChannelID string                 `json:"existing_android_channel_id,omitempty"`
	SmallIcon                string                 `json:"small_icon,omitempty"`
	IOSSound                 string                 `json:"ios_sound,omitempty"`
	AndroidSound             string                 `json:"android_sound,omitempty"`
	TTL                      int                    `json:"ttl,omitempty"`
	IsAndroid                *bool                  `json:"isAndroid,omitempty"`
	IsIos                    *bool                  `json:"isIos,omitempty"`
	IsAnyWeb                 *bool                  `json:"isAnyWeb,omitempty"`
}
