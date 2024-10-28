package configs

type GatewayConfig struct {
	Xendit PaymentConfig `yaml:"xendit"`
	Paypal PaymentConfig `yaml:"paypal"`
}

type PaymentConfig struct {
	APIKey        string   `yaml:"api_key"`
	Label         string   `yaml:"label"`
	CallbackTokens []string `yaml:"callback_tokens"`
}
