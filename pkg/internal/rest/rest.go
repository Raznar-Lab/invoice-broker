package rest

import (
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/pkg/internal/database"
)

func Start(conf *configs.Config, db *database.Database) error {
	fiberConf := fiber.Config{
		TrustedProxies: conf.Web.TrustedProxy,
	}

	if len(fiberConf.TrustedProxies) > 0 {
		fiberConf.EnableTrustedProxyCheck = true
	}

	fiberConf.ProxyHeader = conf.Web.ProxyHeader
	app := fiber.New(fiberConf)
	// Middleware
	app.Use(logger.New())

	initRoutes(app, conf, db)

	return app.Listen(fmt.Sprintf("%s:%s", conf.Web.Bind, conf.Web.Port))
}

func initRoutes(app *fiber.App, conf *configs.Config, db *database.Database) {
	initShortener(app, db)
	initXendit(app, conf.Gateway.Xendit, conf.Web, db)
}

func getApiData(apiToken string, ip string, apiConfigs []configs.APIConfig) *configs.APIConfig {
	for _, apiConfig := range apiConfigs {
		if apiConfig.Token != apiToken {
			continue
		}

		if len(apiConfig.AllowedIPs) == 0 || slices.Contains(apiConfig.AllowedIPs, ip) {
			return &apiConfig
		}
	}
	return nil
}