package services

import (
	"raznar.id/invoice-broker/configs"
	paypal_service "raznar.id/invoice-broker/internal/services/paypal"
	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)

type Services struct {
	Xendit *xendit_service.XenditService
	Paypal *paypal_service.PaypalService
}

func New(c *configs.Config) *Services {
	return &Services{
		Xendit: xendit_service.New(c),
		Paypal: paypal_service.New(c),
	}
}
