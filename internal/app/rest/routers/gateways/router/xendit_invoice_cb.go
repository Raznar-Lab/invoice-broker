package router_xendit_gateway

import (
	"fmt"
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
		CallbackToken: c.Get("x-callback-token"),
		WebhookId:     c.Get("webhook-id"),
	}

	if !slices.Contains(paymentConfig.CallbackTokens, xenditHeader.CallbackToken) {
		return c.SendStatus(fiber.ErrUnauthorized.Code)
	}

	var invoiceData xenInvoice.Invoice
	if err := c.BodyParser(&invoiceData); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	transaction := r.DB.GetTransactionByTrID(invoiceData.GetId())
	if transaction == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	transaction.Status = invoiceData.GetStatus().String()
	r.DB.SilentSave()

	body := c.Body()
	go func() {
		apiConf := r.Config.GetAPIConfigByOrganization(transaction.Organization)
		if apiConf == nil {
			fmt.Printf("An error occured when forwarding webhook data, api conf with organization %s is null\n", transaction.Organization)
			return
		}

		urls := append(apiConf.CallbackURLS, transaction.CallbackURLS...)
		services.Invoice().ForwardWebhookData(body, urls, []string{"x-callback-token", "X-Callback-Token"}, apiConf.Token)
	}()

	return c.SendStatus(fiber.StatusOK)
}
