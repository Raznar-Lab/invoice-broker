package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"raznar.id/invoice-broker/config"
)

func Start(conf *config.Config, rdb *redis.Client) (err error) {
	fiberConf := fiber.Config{
		TrustedProxies: conf.Web.TrustedProxy,
	}

	if len(fiberConf.TrustedProxies) > 0 {
		fiberConf.EnableTrustedProxyCheck = true
	}

	fiberConf.ProxyHeader = conf.Web.ProxyHeader
	app := fiber.New(fiberConf)
	// middleware
	app.Use(logger.New())

	initRoutes(app, conf, rdb)

	return app.Listen(fmt.Sprintf("%s:%s", conf.Web.Bind, conf.Web.Port))
}

func initRoutes(app *fiber.App, conf *config.Config, rdb *redis.Client) {
	initXendit(app, conf.Gateway.Xendit, rdb)
}
