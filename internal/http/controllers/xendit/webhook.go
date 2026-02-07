package xendit_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	xendit_service "raznar.id/invoice-broker/internal/services/xendit"
)

func (x *XenditController) ValidateWebhook(c *gin.Context) {
	log.Debug().
		Str("ip", c.ClientIP()).
		Str("path", c.FullPath()).
		Msg("Xendit webhook received")

	var invoice xenInvoice.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		log.Debug().
			Err(err).
			Msg("Failed to bind Xendit webhook payload")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	cbToken := c.GetHeader(x.paymentConfig.WebhookHeader)

	log.Debug().
		Str("callback_token", cbToken).
		Str("invoice_id", invoice.GetId()).
		Msg("Xendit webhook payload parsed")

	payload := xendit_service.ValidationPayload{
		PaymentConfig: x.paymentConfig,
		Invoice:       &invoice,
		CallbackToken: cbToken,
	}

	if !x.Services.Xendit.ValidateWebhook(payload) {
		log.Debug().
			Str("invoice_id", invoice.GetId()).
			Msg("Xendit webhook validation failed")

		c.Status(http.StatusUnauthorized)
		return
	}

	log.Debug().
		Str("invoice_id", invoice.GetId()).
		Msg("Xendit webhook validated successfully")

	c.Status(http.StatusOK)
}
