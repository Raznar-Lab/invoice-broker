package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"raznar.id/invoice-broker/config"
	"raznar.id/invoice-broker/internal/invoice"
)

func initXendit(app *fiber.App, conf config.PaymentConfig) {
	app.Post(fmt.Sprintf("/gateway/%s/invoice", conf.Label), func(c *fiber.Ctx) (err error) {
		xenHeader := invoice.XenditHeader{
			CallbackToken: c.Get("x-callback-token"),
			WebhookID:     c.Get("webhook-id"),
		}

		code := invoice.ProcessXendit(conf, xenHeader)
		if code == fiber.StatusOK {
			invoice.ForwardWebhookData(c.Body(), conf.CallbackURLs, "x-callback-token", xenHeader.CallbackToken)
		}

		return c.SendStatus(code)
	})
}
