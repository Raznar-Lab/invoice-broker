package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/app/rest/routers"
	"raznar.id/invoice-broker/internal/pkg/database"
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

	routers.Init(app, conf, db)

	return app.Listen(fmt.Sprintf("%s:%s", conf.Web.Bind, conf.Web.Port))
}
