package xendit_service

import (
	"slices"

	"github.com/rs/zerolog/log"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/constants"
	"raznar.id/invoice-broker/internal/notifications"
)

type ValidationPayload struct {
	PaymentConfig *configs.GatewayConfig
	Invoice       *xenInvoice.Invoice
	CallbackToken string
}

func (x *XenditService) ValidateWebhook(payload ValidationPayload) bool {
	// Create a sub-logger for this validation context
	logger := log.With().
		Str("service", "xendit").
		Str("external_id", payload.Invoice.ExternalId).
		Logger()

	valid := slices.Contains(payload.PaymentConfig.Xendit.WebhookTokens, payload.CallbackToken)
	if !valid {
		logger.Warn().
			Str("received_token", payload.CallbackToken).
			Msg("invalid callback token received")
		return false
	}

	// Forward data to the background queue (this is non-blocking now)
	err := notifications.SendWebhook(
		notifications.WebhookPayload{
			Content: payload.Invoice,
			URLS:    payload.PaymentConfig.Xendit.CallbackURLs,
			Header:  constants.XENDIT_WEBHOOK_HEADER,
			Token:   payload.CallbackToken,
		},
	)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to queue webhook forwarding")
		return false
	}

	logger.Info().Msg("webhook validated and queued for forwarding")
	return true
}
