package router_xendit_gateway

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/configs"
	base_router "raznar.id/invoice-broker/internal/app/rest/routers/base"
	"raznar.id/invoice-broker/internal/pkg/database"
)

type XenditGatewayRouter struct {
	base_router.BaseRouter
}

func New(app *fiber.App, config *configs.Config, database *database.Database) base_router.IBaseRouter {
	router := &XenditGatewayRouter{}
	router.Set(app, config, database)

	return router
}

func (r XenditGatewayRouter) Init() {
	xGroup := r.App.Group(fmt.Sprintf("/gateway/%s", r.Config.Gateway.Xendit.Label))

	invGroup := xGroup.Group("/invoice")
	invGroup.Post("/cb", r.InvoiceCallbackHandler)
	invGroup.Post("/", r.InvoiceCreateHandler)
	invGroup.Get("/:id", r.InvoiceGetHandler)
}
