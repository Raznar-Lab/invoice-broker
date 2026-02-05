package paypal_service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/notifications"
)

type WebhookValidationPayload struct {
	PaymentConfig *configs.GatewayConfig
	Headers       http.Header
	RawBody       []byte
}

func paypalAPI(cfg *configs.GatewayConfig) string {
	if cfg.Paypal.Sandbox {
		return "https://api-m.sandbox.paypal.com"
	}
	return "https://api-m.paypal.com"
}

func (p *PaypalService) ValidateWebhook(payload WebhookValidationPayload) bool {
	webhookID := payload.PaymentConfig.Paypal.WebhookToken()
	logger := log.With().
		Str("service", "paypal_webhook").
		Logger()

	logger.Debug().
		Str("webhook_id", webhookID).
		Msg("starting paypal webhook validation")

	token, err := p.getAccessToken(payload.PaymentConfig)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to get paypal access token")
		return false
	}

	logger.Debug().
		Msg("paypal access token acquired")

	body := map[string]interface{}{
		"auth_algo":         payload.Headers.Get("PAYPAL-AUTH-ALGO"),
		"cert_url":          payload.Headers.Get("PAYPAL-CERT-URL"),
		"transmission_id":   payload.Headers.Get("PAYPAL-TRANSMISSION-ID"),
		"transmission_sig":  payload.Headers.Get("PAYPAL-TRANSMISSION-SIG"),
		"transmission_time": payload.Headers.Get("PAYPAL-TRANSMISSION-TIME"),
		"webhook_id":        webhookID,
		"webhook_event":     json.RawMessage(payload.RawBody),
	}

	logger.Debug().
		Interface("verify_payload", body).
		Msg("constructed webhook verification payload")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to marshal webhook verification payload")
		return false
	}

	req, err := http.NewRequest(
		http.MethodPost,
		paypalAPI(payload.PaymentConfig)+"/v1/notifications/verify-webhook-signature",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to create webhook verify request")
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	logger.Debug().
		Str("endpoint", req.URL.String()).
		Msg("sending webhook verification request to paypal")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to call paypal webhook verification API")
		return false
	}
	defer resp.Body.Close()

	logger.Debug().
		Int("http_status", resp.StatusCode).
		Msg("received paypal webhook verification response")

	var result struct {
		VerificationStatus string `json:"verification_status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to decode webhook verification response")
		return false
	}

	logger.Debug().
		Str("verification_status", result.VerificationStatus).
		Msg("paypal webhook verification result")

	if result.VerificationStatus != "SUCCESS" {
		logger.Warn().
			Str("status", result.VerificationStatus).
			Msg("paypal webhook verification failed")
		return false
	}

	logger.Info().
		Msg("paypal webhook verified successfully")

		// Forward data to the background queue (this is non-blocking now)
	err = notifications.SendWebhook(
		notifications.WebhookPayload{
			Content: body,
			URLS:    payload.PaymentConfig.Paypal.CallbackURLs,
			Header:  "",
			Token:   "",
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
