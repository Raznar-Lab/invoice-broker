package paypal_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
)

type CreateInvoicePayload struct {
	PaymentConfig *configs.GatewayConfig
	ID            string
	Description   string
	Amount        float64
	ReturnURL     string
	CancelURL     string
}

func (p *PaypalService) getAccessToken(cfg *configs.GatewayConfig) (string, error) {
	reqBody := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", paypalEndpoint(cfg)+"/v1/oauth2/token", reqBody)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(cfg.Paypal.ApiID, cfg.Paypal.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (p *PaypalService) CreateOrder(payload CreateInvoicePayload) (string, error) {
	logger := log.With().
		Str("service", "paypal_create_order").
		Str("external_id", payload.ID).
		Logger()

	fee := payload.Amount * 0.01
	total := payload.Amount + fee

	token, err := p.getAccessToken(payload.PaymentConfig)
	if err != nil {
		return "", err
	}

	body := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": payload.ID,
				"description":  payload.Description,
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", total),
				},
			},
		},
		"application_context": map[string]string{
			"return_url": payload.ReturnURL,
			"cancel_url": payload.CancelURL,
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", paypalEndpoint(payload.PaymentConfig)+"/v2/checkout/orders", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ID    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, link := range result.Links {
		if link.Rel == "approve" {
			logger.Info().
				Str("order_id", result.ID).
				Str("approval_url", link.Href).
				Msg("paypal order created")
			return link.Href, nil
		}
	}

	return "", fmt.Errorf("approval link not found")
}
