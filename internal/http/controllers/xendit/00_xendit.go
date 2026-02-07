package xendit_controller

import (
	"raznar.id/invoice-broker/configs"
	base_controller "raznar.id/invoice-broker/internal/http/controllers/base"
	"raznar.id/invoice-broker/internal/services"
)

type XenditController struct {
	paymentConfig *configs.PaymentConfig
	base_controller.BaseController
}

func New(c *configs.Config, s *services.Services, paymentConfig *configs.PaymentConfig) *XenditController {

	x := &XenditController{
		paymentConfig: paymentConfig,
	}

	// 2. Set the base dependencies
	x.Set(c, s)

	return x
}
