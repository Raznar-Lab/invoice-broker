package config

type WebConfig struct {
	EnableProxy  string   `yaml:"enable_proxy"`
	ProxyHeader  string   `yaml:"proxy_header"`
	TrustedProxy []string `yaml:"trusted_proxy"`
	Bind         string   `yaml:"bind" validate:"required"`
	Port         string   `yaml:"port" validate:"required"`
}
