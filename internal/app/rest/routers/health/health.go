package router_health

import (
	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	base_router "raznar.id/invoice-broker/internal/app/rest/routers/base"
	"raznar.id/invoice-broker/internal/pkg/database"
)

type HealthRouter struct {
	base_router.BaseRouter
}

func New(app *fiber.App, config *configs.Config, db *database.Database) base_router.IBaseRouter {
	router := &HealthRouter{}
	router.Set(app, config, db)

	return router
}

func (r HealthRouter) Handler(c *fiber.Ctx) (err error) {

	return c.SendStatus(fiber.StatusNoContent)
}

func (r HealthRouter) Init(g fiber.Router) {
	group := g.Group("/health")
	group.Get("/", r.Handler)
}
