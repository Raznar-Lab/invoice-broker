package paypal_controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	paypal_service "raznar.id/invoice-broker/internal/services/paypal"
)

func (x *PaypalController) ValidateIPN(c *gin.Context) {
	log.Debug().
		Str("ip", c.ClientIP()).
		Str("path", c.FullPath()).
		Msg("PayPal IPN received")

	// Read raw body (MANDATORY for IPN)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read PayPal IPN body")

		c.Status(http.StatusBadRequest)
		return
	}

	payload := paypal_service.ValidationPayload{
		PaymentConfig: x.paymentConfig,
		RawBody:       rawBody,
	}

	if !x.Services.Paypal.ValidateIPN(payload) {
		log.Warn().
			Msg("PayPal IPN validation failed")

		// PayPal expects 200 even on logical failure
		c.Status(http.StatusOK)
		return
	}

	log.Info().
		Msg("PayPal IPN validated successfully")

	// MUST return 200 OK or PayPal will retry
	c.Status(http.StatusOK)
}
