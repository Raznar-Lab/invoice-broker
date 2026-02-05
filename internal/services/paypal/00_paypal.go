package paypal_service

import (
	"raznar.id/invoice-broker/configs"
	base_service "raznar.id/invoice-broker/internal/services/base"
)

type PaypalService struct {
	base_service.Service
}

func paypalEndpoint(cfg *configs.GatewayConfig) string {
	if cfg.Paypal.Sandbox {
		return "https://api-m.sandbox.paypal.com"
	}
	return "https://api-m.paypal.com"
}

func New(c *configs.Config) *PaypalService {
	v := &PaypalService{}
	v.Set(c)

	return v
}
