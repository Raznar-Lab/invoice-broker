package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"raznar.id/invoice-broker/config"
)

func Start(conf *config.Config) (err error) {
	fiberConf := fiber.Config{
		TrustedProxies: conf.Web.TrustedProxy,
	}

	if len(fiberConf.TrustedProxies) > 0 {
		fiberConf.EnableTrustedProxyCheck = true
	}

	fiberConf.ProxyHeader = conf.Web.ProxyHeader
	app := fiber.New(fiberConf)
	initRoutes(app, conf)

	// middleware
	app.Use(logger.New())

	return app.Listen(fmt.Sprintf("%s:%s", conf.Web.Bind, conf.Web.Port))
}

func initRoutes(app *fiber.App, conf *config.Config) {
	initXendit(app, conf.Gateway.Xendit)
}
