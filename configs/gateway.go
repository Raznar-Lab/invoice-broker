package configs

type GatewayConfig struct {
	Xendit PaymentConfig `json:"xendit"`
	Paypal PaymentConfig `json:"paypal"`
}

type PaymentConfig struct {
	APIKey        string   `json:"api_key"`
	WebhookTokens []string `json:"webhook_tokens"`
	CallbackURLs  []string `json:"callback_urls"`
}