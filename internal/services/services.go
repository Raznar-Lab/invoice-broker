package services

import (
	"raznar.id/invoice-broker/configs"
	pakasir_service "raznar.id/invoice-broker/internal/services/pakasir"
	paypal_service "raznar.id/invoice-broker/internal/services/paypal"
	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)

type Services struct {
	Xendit  *xendit_service.XenditService
	Paypal  *paypal_service.PaypalService
	Pakasir *pakasir_service.PakasirService
}

func New(c *configs.Config) *Services {
	return &Services{
		Xendit:  xendit_service.New(c),
		Paypal:  paypal_service.New(c),
		Pakasir: pakasir_service.New(c),
	}
}
