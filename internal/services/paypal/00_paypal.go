package paypal_service

import (
	"raznar.id/invoice-broker/configs"
	base_service "raznar.id/invoice-broker/internal/services/base"
)

type PaypalService struct {
	base_service.Service
}

func New(c *configs.Config) *PaypalService {
	v := &PaypalService{}
	v.Set(c)

	return v
}
