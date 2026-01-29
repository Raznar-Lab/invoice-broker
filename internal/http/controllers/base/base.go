package base_controller

import (
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/services"
)

type BaseController struct {
	Config      *configs.Config
	Services    *services.Services
}

func (b *BaseController) Set(c *configs.Config, s *services.Services) {
	b.Config = c
	b.Services = s
}
