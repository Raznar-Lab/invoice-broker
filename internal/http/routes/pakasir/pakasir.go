package pakasir_route

import (
	"strings"

	"github.com/gin-gonic/gin"
	"raznar.id/invoice-broker/configs"
	pakasir_controller "raznar.id/invoice-broker/internal/http/controllers/pakasir"
	"raznar.id/invoice-broker/internal/http/middlewares"
	base_routes "raznar.id/invoice-broker/internal/http/routes/base"
	"raznar.id/invoice-broker/internal/services"
)

type PakasirRoute struct {
	base_routes.Router
}

func (r *PakasirRoute) Register() {
	for labelKey, cfg := range r.Config.Gateway {
		if cfg == nil {
			return
		}

		ctrl := pakasir_controller.New(r.Config, r.Services, &cfg.Pakasir)

		label := strings.ToLower(labelKey)

		// Webhook routes (NO auth)
		webhook := r.RG.Group("/pakasir/webhook/" + label)
		{
			webhook.POST("", ctrl.ValidateWebhook)
		}
	}
}

func New(c *configs.Config, s *services.Services, m *middlewares.Middlewares, rg *gin.RouterGroup) *PakasirRoute {
	x := &PakasirRoute{}
	x.Set(c, s, m, rg)

	return x
}
