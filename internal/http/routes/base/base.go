package base_routes

import (
	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/http/middlewares"
	"raznar.id/invoice-broker/internal/services"
)

type Router struct {
	Config   *configs.Config
	RG       *gin.RouterGroup
	Services *services.Services
	Middlewares *middlewares.Middlewares
}

func (r *Router) Set(c *configs.Config, services *services.Services, middlewares *middlewares.Middlewares, rg *gin.RouterGroup) {
	r.Config = c
	r.RG = rg
	r.Services = services
	r.Middlewares = middlewares
}
