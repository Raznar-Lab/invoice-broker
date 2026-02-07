package xendit_service

import (
	"github.com/rs/zerolog/log"
	xenInvoice "github.com/xendit/xendit-go/v5/invoice"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/events"
)

type ValidationPayload struct {
	PaymentConfig *configs.PaymentConfig
	Invoice       *xenInvoice.Invoice
	CallbackToken string
}

func (x *XenditService) ValidateWebhook(payload ValidationPayload) bool {
	// Create a sub-logger for this validation context
	logger := log.With().
		Str("service", "xendit").
		Str("external_id", payload.Invoice.ExternalId).
		Logger()

	if payload.PaymentConfig.WebhookToken != payload.CallbackToken {
		logger.Warn().
			Str("received_token", payload.CallbackToken).
			Msg("invalid callback token received")
		return false
	}

	fee := float64(0)

	for _, v := range payload.Invoice.GetFees() {
		fee += float64(v.GetValue())
	}
	err := events.Emit(&events.PaymentWebhookEvent{
		PaymentConfig: payload.PaymentConfig,
		Raw:           payload.Invoice,
		ExternalID:    payload.Invoice.ExternalId,
		ID:            payload.Invoice.GetId(),
		Status:        payload.Invoice.Status.String(),
		Currency:      payload.Invoice.Currency.String(),
		Description:   payload.Invoice.GetDescription(),
		Amount:        payload.Invoice.GetAmount(),
		Fee:           fee,
		Gateway:       "Xendit",
	})

	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to emit event")
		return false
	}

	logger.Info().Msg("webhook validated and queued for forwarding")
	return true
}
