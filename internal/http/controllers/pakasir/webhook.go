package pakasir_controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	pakasir_service "raznar.id/invoice-broker/internal/services/pakasir"
)

func (x *PakasirController) ValidateWebhook(c *gin.Context) {
	log.Debug().
		Str("ip", c.ClientIP()).
		Str("path", c.FullPath()).
		Msg("PakKasir webhook received")

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to read pakasir webhook body")

		c.Status(http.StatusBadRequest)
		return
	}

	payload := pakasir_service.WebhookValidationPayload{
		PaymentConfig: x.paymentConfig,
		RawBody:       rawBody,
	}

	if !x.Services.Pakasir.ValidateWebhook(payload) {
		log.Warn().
			Msg("pakasir webhook validation failed")

		c.Status(http.StatusBadRequest)
		return
	}

	log.Info().
		Msg("pakasir webhook validated successfully")

	c.Status(http.StatusOK)
}
