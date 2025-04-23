package routers

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	base_router "raznar.id/invoice-broker/internal/app/rest/routers/base"
	router_gateways "raznar.id/invoice-broker/internal/app/rest/routers/gateways"
	router_health "raznar.id/invoice-broker/internal/app/rest/routers/health"

	router_shortener "raznar.id/invoice-broker/internal/app/rest/routers/shortener"
	"raznar.id/invoice-broker/internal/pkg/database"
)

func Init(app *fiber.App, config *configs.Config, db *database.Database) {
	routers := []base_router.IBaseRouter{
		router_gateways.New(app, config, db),
		router_shortener.New(app, config, db),
		router_health.New(app, config, db),
	}

	for _, r := range routers {
		r.Init(app.Group("/api"))
	}
}
