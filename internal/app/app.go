package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/http/routes"
	"raznar.id/invoice-broker/internal/services"
)

func Start(conf *configs.Config, s *services.Services) error {
	if !gin.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Initialize Routes from the separate register function
	routes.Init(r, conf, s)

	log.Info().Str("bind", conf.Server.Bind).Str("port", conf.Server.Port).Msg("REST API is live")
	return r.Run(":" + conf.Server.Port)
}
