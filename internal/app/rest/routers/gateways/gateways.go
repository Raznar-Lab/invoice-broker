package router_gateways

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	base_router "raznar.id/invoice-broker/internal/app/rest/routers/base"
	router_xendit_gateway "raznar.id/invoice-broker/internal/app/rest/routers/gateways/router"

	"raznar.id/invoice-broker/internal/pkg/database"
)

type GatewayRouter struct {
	base_router.BaseRouter
}

func (r GatewayRouter) Init(g fiber.Router) {
	xendit := router_xendit_gateway.New(r.App, r.Config, r.DB)
	xendit.Init(g)
}

func New(app *fiber.App, config *configs.Config, database *database.Database) base_router.IBaseRouter {
	router := &GatewayRouter{}
	router.Set(app, config, database)

	return router
}
