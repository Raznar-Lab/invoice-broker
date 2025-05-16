package router_xendit_gateway

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/internal/app/services"
	"raznar.id/invoice-broker/internal/dtos"
	"raznar.id/invoice-broker/internal/pkg/constants"
)

type xenditInvoiceCallbackHeader struct {
	CallbackToken string
	WebhookId     string
}

func (r XenditGatewayRouter) InvoiceCallbackHandler(c *fiber.Ctx) (err error) {
	paymentConfig := r.Config.Gateway.Xendit

	// webhook id is unused, idk what to implement atm
	xenditHeader := xenditInvoiceCallbackHeader{
		WebhookId: c.Get("webhook-id"),
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
		tokenHeaders := []string{"x-callback-token", "X-Callback-Token", "X-CALLBACK-TOKEN"}
		if transaction != nil {
			apiConf := r.Config.GetAPIConfigByOrganization(transaction.Organization)
			if apiConf != nil {
				services.Invoice().ForwardCallbackData(body, apiConf.CallbackURLS, tokenHeaders, apiConf.Token)
				services.Invoice().ForwardCallbackData(body, transaction.CallbackURLS, tokenHeaders, apiConf.Token)
			} else {
				services.Invoice().ForwardCallbackData(body, transaction.CallbackURLS, tokenHeaders, xenditHeader.CallbackToken)
			}

		}

		services.Invoice().ForwardCallbackData(body, paymentConfig.CallbackURLS, tokenHeaders, xenditHeader.CallbackToken)
		invoiceFees := 0.0
		for _, fee := range invoiceData.GetFees() {
			invoiceFees += float64(fee.GetValue())
		}

		services.Invoice().ForwardWebhookData(dtos.WebhookInvoiceData{
			InvoiceId:  invoiceData.GetId(),
			ExternalId: invoiceData.GetExternalId(),
			Amount:     invoiceData.GetAmount(),
			Fee:        invoiceFees,
			Gateway:    constants.GATEWAY_XENDIT,
			Currency:   invoiceData.GetCurrency().String(),
			Status:     invoiceData.GetStatus().String(),
			URL:        invoiceData.GetInvoiceUrl(),
			Description: invoiceData.GetDescription(),
		}, paymentConfig.Webhooks)
	}()

	return c.SendStatus(fiber.StatusOK)
}
