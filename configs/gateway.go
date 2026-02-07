package configs

type GatewayConfig struct {
	Xendit PaymentConfig `json:"xendit"`
	Paypal PaymentConfig `json:"paypal"`
}

type NotificationConfig struct {
	Webhooks []string `json:"webhooks"`
}
type PaymentConfig struct {
	ApiID         string             `json:"api_ID"`
	APIKey        string             `json:"api_key"`
	WebhookToken  string             `json:"webhook_token"`
	WebhookHeader string             `json:"webhook_header"`
	CallbackURLs  []string           `json:"callback_urls"`
	Sandbox       bool               `json:"sandbox"`
	Notification  NotificationConfig `json:"notification"`
}
