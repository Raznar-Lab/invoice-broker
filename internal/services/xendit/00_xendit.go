package xendit_service

import (
	"raznar.id/invoice-broker/configs"
	base_service "raznar.id/invoice-broker/internal/services/base"
)

type XenditService struct {
	base_service.Service
}

func New(c *configs.Config) *XenditService {
	v := &XenditService{}
	v.Set(c)

	return v
}