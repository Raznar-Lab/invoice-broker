package paypal_service

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/constants"
	"raznar.id/invoice-broker/internal/notifications"
)

type ValidationPayload struct {
	PaymentConfig *configs.GatewayConfig
	RawBody       []byte
}

func (p *PaypalService) ValidateIPN(payload ValidationPayload) bool {
	logger := log.With().
		Str("service", "paypal").
		Logger()

	verifyURL := constants.PAYPAL_IPN_LIVE_VERIFY_URL
	if payload.PaymentConfig.Paypal.Sandbox {
		verifyURL = constants.PAYPAL_IPN_SANDBOX_VERIFY_URL
	}

	// Step 1: Verify IPN with PayPal
	verifyBody := append([]byte("cmd=_notify-validate&"), payload.RawBody...)

	req, err := http.NewRequest(
		http.MethodPost,
		verifyURL,
		bytes.NewBuffer(verifyBody),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create IPN verify request")
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to verify IPN with PayPal")
		return false
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if string(respBody) != "VERIFIED" {
		logger.Warn().
			Str("response", string(respBody)).
			Msg("invalid PayPal IPN")
		return false
	}

	// Step 2: Parse verified payload
	values, err := url.ParseQuery(string(payload.RawBody))
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse IPN body")
		return false
	}

	invoiceID := values.Get("invoice")
	status := values.Get("payment_status")

	logger.Info().
		Str("invoice_id", invoiceID).
		Str("status", status).
		Msg("paypal IPN verified")

	// Step 3: Forward internally (async)
	err = notifications.SendWebhook(
		values, // raw parsed IPN
		payload.PaymentConfig.Paypal.CallbackURLs,
		constants.PAYPAL_WEBHOOK_HEADER,
		"", // no token in PayPal
	)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to queue PayPal webhook forwarding")
		return false
	}

	logger.Info().
		Msg("paypal webhook validated and queued for forwarding")

	return true
}
