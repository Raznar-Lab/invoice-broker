package router_shortener

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	base_router "raznar.id/invoice-broker/internal/app/rest/routers/base"
	"raznar.id/invoice-broker/internal/pkg/database"
)

type ShortenerRoute struct {
	base_router.BaseRouter
}

func New(app *fiber.App, config *configs.Config, db *database.Database) base_router.IBaseRouter {
	router := &ShortenerRoute{}
	router.Set(app, config, db)

	return router
}

func (r ShortenerRoute) Init() {
	group := r.App.Group("/shortener")
	group.Get("/", r.ShortenerGetHandler)
}
