package paypal_controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	paypal_service "raznar.id/invoice-broker/internal/services/paypal"
)

func (x *PaypalController) Webhook(c *gin.Context) {
	log.Debug().
		Str("ip", c.ClientIP()).
		Str("path", c.FullPath()).
		Msg("paypal webhook received")

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to read paypal webhook body")

		c.Status(http.StatusBadRequest)
		return
	}

	payload := paypal_service.WebhookValidationPayload{
		PaymentConfig: x.paymentConfig,
		Headers:       c.Request.Header,
		RawBody:       rawBody,
	}

	if !x.Services.Paypal.ValidateWebhook(payload) {
		log.Warn().
			Msg("paypal webhook validation failed")

		c.Status(http.StatusBadRequest)
		return
	}

	log.Info().
		Msg("paypal webhook verified successfully")

	c.Status(http.StatusOK)
}
