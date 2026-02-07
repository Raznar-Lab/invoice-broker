package paypal_service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/events"
)

type WebhookValidationPayload struct {
	PaymentConfig *configs.PaymentConfig
	Headers       http.Header
	RawBody       []byte
}

type PaypalWebhookBody struct {
	ID           string `json:"id"`            // WH-...
	EventType    string `json:"event_type"`    // PAYMENT.CAPTURE.COMPLETED
	ResourceType string `json:"resource_type"` // capture
	Summary      string `json:"summary"`

	Resource struct {
		ID     string `json:"id"`     // Capture ID (8HW...)
		Status string `json:"status"` // COMPLETED

		Amount struct {
			Value        float64 `json:"value,string"`
			CurrencyCode string  `json:"currency_code"`
		} `json:"amount"`

		SellerReceivableBreakdown struct {
			GrossAmount struct {
				Value        float64 `json:"value,string"`
				CurrencyCode string  `json:"currency_code"`
			} `json:"gross_amount"`

			NetAmount struct {
				Value        float64 `json:"value,string"`
				CurrencyCode string  `json:"currency_code"`
			} `json:"net_amount"`

			PaypalFee struct {
				Value        float64 `json:"value,string"`
				CurrencyCode string  `json:"currency_code"`
			} `json:"paypal_fee"`
		} `json:"seller_receivable_breakdown"`
	} `json:"resource"`
}


func paypalAPI(cfg *configs.PaymentConfig) string {
	if cfg.Sandbox {
		return "https://api-m.sandbox.paypal.com"
	}
	return "https://api-m.paypal.com"
}

func (p *PaypalService) ValidateWebhook(payload WebhookValidationPayload) bool {
	paypalConfig := payload.PaymentConfig
	webhookID := paypalConfig.WebhookToken
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

	var paypalBody PaypalWebhookBody
	if err := json.Unmarshal(payload.RawBody, &paypalBody); err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal paypal body")
		return false
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

	event := &events.PaymentWebhookEvent{
		ID:            paypalBody.Resource.ID, // The actual Transaction ID
		ExternalID:    paypalBody.ID,          // The Webhook ID
		Status:        strings.ToLower(paypalBody.Resource.Status),
		Amount:        paypalBody.Resource.Amount.Value,
		Currency:      paypalBody.Resource.Amount.CurrencyCode,
		Gateway:       "Paypal",
		Description:   paypalBody.Summary,
		PaymentConfig: paypalConfig,
		Raw:           payload.RawBody,
		Fee:           paypalBody.Resource.SellerReceivableBreakdown.PaypalFee.Value,
	}
	err = events.Emit(event)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to emit event")
		return false
	}

	logger.Info().Msg("webhook validated and queued for forwarding")
	return true
}
