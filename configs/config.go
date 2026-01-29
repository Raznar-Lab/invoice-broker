package configs

import (
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server  ServerConfig              `envPrefix:"SERVER_" validate:"required"`
	Gateway map[string]*GatewayConfig `env:"-"`
}

var validate = validator.New()

func New() (*Config, error) {
	config := &Config{
		Gateway: make(map[string]*GatewayConfig),
	}

	// 1. Parse Server Config
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	// 2. Parse Dynamic Gateways from flat ENV vars
	// Pattern: GATEWAY_STOREALPHA_XENDIT_API_KEY
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "GATEWAY_") {
			parts := strings.Split(key, "_")
			if len(parts) < 4 {
				continue
			}

			label := strings.ToLower(parts[1])    // e.g., "storealpha"
			provider := strings.ToLower(parts[2]) // e.g., "xendit"
			field := strings.Join(parts[3:], "_") // e.g., "API_KEY"

			// Get or Create the GatewayConfig for this label
			gc, ok := config.Gateway[label]
			if !ok || gc == nil {
				gc = &GatewayConfig{}
				config.Gateway[label] = gc
			}
			// Assign values based on provider and field
			switch provider {
			case "xendit":
				fillPaymentConfig(&gc.Xendit, field, value)
			case "paypal":
				fillPaymentConfig(&gc.Paypal, field, value)
			}

			config.Gateway[label] = gc
		}
	}

	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Helper to fill the sub-struct fields
func fillPaymentConfig(p *PaymentConfig, field, value string) {
	log.Debug().
		Str("field", field).
		Msg("Filling payment config field")

	switch field {
	case "API_KEY":
		// don't log the value, it's sensitive
		p.APIKey = value

	case "WEBHOOK_TOKENS":
		p.WebhookTokens = strings.Split(value, ",")
		log.Debug().
			Int("count", len(p.WebhookTokens)).
			Msg("Parsed webhook tokens")

	case "SANDBOX":
		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Warn().
				Str("value", value).
				Msg("Invalid SANDBOX value, defaulting to false")
			return
		}
		p.Sandbox = v

	case "CALLBACK_URLS":
		p.CallbackURLs = strings.Split(value, ",")
		log.Debug().
			Int("count", len(p.CallbackURLs)).
			Msg("Parsed callback URLs")
	}

}
