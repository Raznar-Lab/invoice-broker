package router_xendit_gateway

import (
	"slices"
	"github.com/gofiber/fiber/v2"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/internal/app/services"
)

type xenditInvoiceCallbackHeader struct {
	CallbackToken string
	WebhookId     string
}

func (r XenditGatewayRouter) InvoiceCallbackHandler(c *fiber.Ctx) (err error) {
	paymentConfig := r.Config.Gateway.Xendit

	// webhook id is unused, idk what to implement atm
	xenditHeader := xenditInvoiceCallbackHeader{
		WebhookId:     c.Get("webhook-id"),
	}

	if c.Get("x-callback-token") != "" {
		xenditHeader.CallbackToken = c.Get("x-callback-token")
	} else if c.Get("X-Callback-Token") != "" {
		xenditHeader.CallbackToken = c.Get("X-Callback-Token")
	} else if c.Get("X-CALLBACK-TOKEN") != "" {
		xenditHeader.CallbackToken = c.Get("X-CALLBACK-TOKEN")
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !slices.Contains(paymentConfig.CallbackTokens, xenditHeader.CallbackToken) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var invoiceData xenInvoice.Invoice
	if err := c.BodyParser(&invoiceData); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	transaction := r.DB.GetTransactionByTrID(invoiceData.GetId())
	if transaction != nil {
		transaction.Status = invoiceData.GetStatus().String()
		r.DB.SilentSave()
	}

	body := c.Body()
	go func() {
		urls := paymentConfig.CallbackURLS

		apiConf := r.Config.GetAPIConfigByOrganization(transaction.Organization)
		if apiConf != nil {
			urls = append(urls, apiConf.CallbackURLS...)
		} 


		if transaction != nil {
			urls = append(urls, transaction.CallbackURLS...)
		}
		services.Invoice().ForwardWebhookData(body, urls, []string{"x-callback-token", "X-Callback-Token"}, apiConf.Token)
	}()

	return c.SendStatus(fiber.StatusOK)
}
