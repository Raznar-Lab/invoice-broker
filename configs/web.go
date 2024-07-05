package configs

type WebConfig struct {
	EnableProxy  string      `yaml:"enable_proxy"`
	ProxyHeader  string      `yaml:"proxy_header"`
	TrustedProxy []string    `yaml:"trusted_proxy"`
	Bind         string      `yaml:"bind" validate:"required"`
	Port         string      `yaml:"port" validate:"required"`
	APIConfigs   []APIConfig `yaml:"api_configs" validate:"required"`
}

func (w WebConfig) GetAPIConfig(organization string) *APIConfig {
	for _, apiConfig := range w.APIConfigs {
		if apiConfig.Organization == organization {
			return &apiConfig
		}
	}

	return nil
}

func (w WebConfig) GetAPIConfigByToken(token string) *APIConfig {
	for _, apiConfig := range w.APIConfigs {
		if apiConfig.Token == token {
			return &apiConfig
		}
	}

	return nil
}

type APIConfig struct {
	Organization string   `yaml:"organization" validate:"required"`
	Token        string   `yaml:"token" validate:"required"`
	AllowedIPs   []string `yaml:"allowed_ips" validate:"required"`
	CallbackURLS []string `yaml:"callback_urls" validate:"required"`
}
