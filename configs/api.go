package configs

import "slices"

type APIConfig struct {
	Organization string   `yaml:"organization" validate:"required"`
	Token        string   `yaml:"token" validate:"required"`
	AllowedIPs   []string `yaml:"allowed_ips" validate:"required"`
	CallbackURLS []string `yaml:"callback_urls" validate:"required"`
}

func (c Config) GetAPIConfigByOrganization(organization string) (apiConf *APIConfig) {
	for _, apiConfig := range c.APIConfigs {
		if apiConfig.Organization == organization {
			apiConf = &apiConfig
			return
		}
	}

	return
}

func (c Config) GetAPIConfig(apiToken string, ip string) (apiConf *APIConfig) {
	for _, apiConfig := range c.APIConfigs {
		if apiConfig.Token != apiToken || (len(apiConfig.AllowedIPs) != 0 && slices.Contains(apiConfig.AllowedIPs, ip)) {
			continue
		}

		apiConf = &apiConfig
		return
	}

	return
}

func (c Config) GetAPIConfigByToken(token string) *APIConfig {
	for _, apiConfig := range c.APIConfigs {
		if apiConfig.Token == token {
			return &apiConfig
		}
	}

	return nil
}