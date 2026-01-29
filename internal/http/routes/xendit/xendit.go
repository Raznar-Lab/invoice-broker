package xendit_route

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
	xendit_controller "raznar.id/invoice-broker/internal/http/controllers/xendit"
	"raznar.id/invoice-broker/internal/http/middlewares"
	base_routes "raznar.id/invoice-broker/internal/http/routes/base"
	"raznar.id/invoice-broker/internal/services"
)

type XenditRoute struct {
	base_routes.Router
}

func (r *XenditRoute) Register() {
	// Loop through all gateways defined in the config
	for label, cfg := range r.Config.Gateway {
		// Initialize the controller for this specific gateway
		ctrl := xendit_controller.New(r.Config, r.Services, r.Middlewares, cfg)

		// Create a group for this label: /api/xendit/{label}
		g := r.RG.Group(fmt.Sprintf("/xendit/%s", label))
		{
			g.POST("/invoice", ctrl.CreateInvoice)
		}

		gg := r.RG.Group(fmt.Sprintf("/webhook/xendit/%s", label))
		{
			gg.POST("/", ctrl.ValidateWebhook)
		}
	}
}

func New(c *configs.Config, s *services.Services, m *middlewares.Middlewares, rg *gin.RouterGroup) *XenditRoute {
	x := &XenditRoute{}
	x.Set(c, s, m, rg)

	return x
}
