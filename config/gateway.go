package config

type GatewayConfig struct {
	Xendit PaymentConfig `yaml:"xendit"`
	Paypal PaymentConfig `yaml:"paypal"`
}

type PaymentConfig struct {
	APIKey        string   `yaml:"api_key"`
	AdminToken    []string `yaml:"admin_token"`
	AdminIP       []string `yaml:"admin_ip"`
	Label         string   `yaml:"label"`
	WebhookTokens []string `yaml:"webhook_tokens"`

	// not recommended though.. only for apps like WHMCS, in a custom application by my own i wouldn't use this. i made this for commercial purposes.
	CallbackURLs []string `yaml:"callback_urls"`
}
