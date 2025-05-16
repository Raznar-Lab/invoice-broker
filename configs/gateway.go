package configs

type GatewayConfig struct {
	Xendit PaymentConfig `yaml:"xendit"`
	Paypal PaymentConfig `yaml:"paypal"`
}

type PaymentConfig struct {
	APIKey         string                          `yaml:"api_key"`
	Label          string                          `yaml:"label"`
	Webhooks       map[string]PaymentWebhookConfig `yaml:"webhooks"`
	CallbackTokens []string                        `yaml:"callback_tokens"`
	CallbackURLS   []string                        `yaml:"callback_urls"`
}

type PaymentWebhookConfig struct {
	URL         string            `yaml:"url"`
	ContentType string            `yaml:"content_type"`
	Content     string            `yaml:"content"`
	Variables   map[string]string `yaml:"variables"`
}
