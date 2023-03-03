package twilio

type WebhookRequestParameters struct {
	MessageSid          string `json:"MessageSid"`
	SmsSid              string `json:"SmsSid"`
	AccountSid          string `json:"AccountSid"`
	MessagingServiceSid string `json:"MessagingServiceSid"`
	From                string `json:"From"`
	To                  string `json:"To"`
	Body                string `json:"Body"`
	NumMedia            int    `json:"NumMedia"`
	MediaContentType0   string `json:"MediaContentType0"`
	MediaUrl0           string `json:"MediaUrl0"`
	MediaContentType1   string `json:"MediaContentType1"`
	MediaUrl2           string `json:"MediaUrl2"`
	MediaContentType3   string `json:"MediaContentType3"`
	MediaUrl3           string `json:"MediaUrl3"`
}
