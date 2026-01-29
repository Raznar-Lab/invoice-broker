package paypal_controller

import (
	"raznar.id/invoice-broker/configs"
	base_controller "raznar.id/invoice-broker/internal/http/controllers/base"
	"raznar.id/invoice-broker/internal/services"
)

type PaypalController struct {
	paymentConfig *configs.GatewayConfig
	base_controller.BaseController
}

func New(c *configs.Config, s *services.Services, paymentConfig *configs.GatewayConfig) *PaypalController {

	x := &PaypalController{
		paymentConfig: paymentConfig,
	}

	// 2. Set the base dependencies
	x.Set(c, s)

	return x
}
