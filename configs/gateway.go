package configs

type GatewayConfig struct {
	Xendit PaymentConfig `json:"xendit"`
	Paypal PaymentConfig `json:"paypal"`
}

type PaymentConfig struct {
	ApiID         string   `json:"api_ID"`
	APIKey        string   `json:"api_key"`
	WebhookTokens []string `json:"webhook_tokens"`
	CallbackURLs  []string `json:"callback_urls"`
	Sandbox       bool     `json:"sandbox"`
}

func (p PaymentConfig) WebhookToken() string {
	if len(p.WebhookTokens) == 0 {
		return ""
	}
	return p.WebhookTokens[0]
}
