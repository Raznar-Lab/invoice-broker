package paypal_route

import (
	"strings"

	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
	paypal_controller "raznar.id/invoice-broker/internal/http/controllers/paypal"
	"raznar.id/invoice-broker/internal/http/middlewares"
	base_routes "raznar.id/invoice-broker/internal/http/routes/base"
	"raznar.id/invoice-broker/internal/services"
)

type PaypalRoute struct {
	base_routes.Router
}

func (r *PaypalRoute) Register() {
	for labelKey, cfg := range r.Config.Gateway {
		if cfg == nil {
			return
		}

		ctrl := paypal_controller.New(
			r.Config,
			r.Services,
			&cfg.Paypal,
		)

		label := strings.ToLower(labelKey)

		// Authenticated API (your client creates orders)
		api := r.RG.Group("/paypal/" + label)
		{
			api.Use(r.Middlewares.Auth)
			api.POST("/order", ctrl.CreateOrder) // CREATE ORDER
		}

		// Unauthenticated API (PayPal callbacks)
		r.RG.POST("/paypal/webhook/"+label, ctrl.Webhook) // webhook from PayPal
		r.RG.GET("/paypal/return/"+label, ctrl.Return)    // user redirected after approval

	}
}

func New(c *configs.Config, s *services.Services, m *middlewares.Middlewares, rg *gin.RouterGroup) *PaypalRoute {
	x := &PaypalRoute{}
	x.Set(c, s, m, rg)

	return x
}
