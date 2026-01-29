package base_controller

import (
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/http/middlewares"
	"raznar.id/invoice-broker/internal/services"
)

type BaseController struct {
	Config      *configs.Config
	Services    *services.Services
	Middlewares *middlewares.Middlewares
}

func (b *BaseController) Set(c *configs.Config, s *services.Services, m *middlewares.Middlewares) {
	b.Config = c
	b.Services = s
	b.Middlewares = m
}
