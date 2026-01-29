package services

import (
	"raznar.id/invoice-broker/configs"
	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)

type Services struct {
	Xendit *xendit_service.XenditService
}

func New(c *configs.Config) *Services {
	return &Services{
		Xendit: xendit_service.New(c),
	}
}
