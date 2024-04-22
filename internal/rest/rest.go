package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"raznar.id/invoice-broker/config"
	"raznar.id/invoice-broker/constants"
	"raznar.id/invoice-broker/internal/invoice"
)

type GatewayInvoiceData struct {
	ID          string
	Status      string
	Gateway     constants.Gateway
	Currency    constants.Currency
	GuildID     string
	Amount      float64
	Fee         float64
	Description string
}

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

	return app.Listen(fmt.Sprintf("%s:%s", conf.Web.Bind, conf.Web.Port))
}

func initRoutes(app *fiber.App, conf *config.Config) {
	app.Post(fmt.Sprintf("/gateway/%s/invoice", constants.GATEWAY_XENDIT_ID.String()), func(c fiber.Ctx) (err error) {
		xenHeader := invoice.XenditHeader{
			CallbackToken: c.Get("x-callback-token"),
			WebhookID:     c.Get("webhook-id"),
		}

		code := invoice.ProcessXendit(conf.Gateway.Xendit, xenHeader)
		if code == fiber.StatusOK {
			invoice.ForwardWebhookData(c.Body(), conf.Gateway.Xendit.CallbackURLs)
		}

		return c.SendStatus(code)
	})
}
