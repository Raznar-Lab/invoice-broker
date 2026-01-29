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
			cfg,
		)

		label := strings.ToLower(labelKey)

		// Authenticated API
		api := r.RG.Group("/paypal/" + label)
		{
			api.Use(r.Middlewares.Auth)
			api.POST("/invoice", ctrl.CreateInvoice)
		}

		// IPN endpoint (NO auth)
		ipn := r.RG.Group("/paypal/ipn/" + label)
		{
			ipn.POST("", ctrl.ValidateIPN)
		}
	}
}


func New(c *configs.Config, s *services.Services, m *middlewares.Middlewares, rg *gin.RouterGroup) *PaypalRoute {
	x := &PaypalRoute{}
	x.Set(c, s, m, rg)

	return x
}
