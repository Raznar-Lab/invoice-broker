package configs

type ServerConfig struct {
	ApiToken     string   `env:"API_TOKEN"`
	EnableProxy  bool     `env:"ENABLE_PROXY" envDefault:"false"`
	ProxyHeader  string   `env:"PROXY_HEADER" envDefault:"X-Forwarded-For"`
	TrustedProxy []string `env:"TRUSTED_PROXY" envSeparator:","`
	Bind         string   `env:"BIND" validate:"required"`
	Port         string   `env:"PORT" validate:"required"`
}

