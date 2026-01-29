package xendit_route

import (
	"fmt"
	"strings"

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
	for labelKey, cfg := range r.Config.Gateway {
		if cfg == nil {
			return
		}

		ctrl := xendit_controller.New(r.Config, r.Services, cfg)

		label := fmt.Sprintf("/xendit/%s", strings.ToLower(labelKey))

		// API routes (authenticated)
		api := r.RG.Group(label)
		{
			api.Use(r.Middlewares.Auth)
			api.POST("/invoice", ctrl.CreateInvoice)
		}

		// Webhook routes (NO auth)
		webhook := r.RG.Group("/webhook" + label)
		{
			webhook.POST("", ctrl.ValidateWebhook)
		}
	}
}

func New(c *configs.Config, s *services.Services, m *middlewares.Middlewares, rg *gin.RouterGroup) *XenditRoute {
	x := &XenditRoute{}
	x.Set(c, s, m, rg)

	return x
}
