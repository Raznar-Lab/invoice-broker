package pakasir_service

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/events"
)

const pakasirAPIBase = "https://app.pakasir.com"

type WebhookValidationPayload struct {
	PaymentConfig *configs.PaymentConfig
	RawBody       []byte
}

type PakasirWebhookBody struct {
	OrderID string `json:"order_id"`
	Amount  any    `json:"amount"`
	Status  string `json:"status"`
}

type transactionDetailResponse struct {
	Transaction struct {
		Status string `json:"status"`
	} `json:"transaction"`
}

func (p *PakasirService) ValidateWebhook(payload WebhookValidationPayload) bool {
	logger := log.With().
		Str("service", "pakasir_webhook").
		Logger()

	var body PakasirWebhookBody
	if err := json.Unmarshal(payload.RawBody, &body); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to unmarshal pakasir webhook body")
		return false
	}

	if body.OrderID == "" || body.Amount == nil {
		logger.Warn().
			Msg("pakasir webhook missing order_id or amount")
		return false
	}

	if !strings.Contains(body.OrderID, "-") {
		logger.Warn().
			Str("order_id", body.OrderID).
			Msg("pakasir webhook received invalid order_id format")
		return false
	}

	amount, err := parseAmount(body.Amount)
	if err != nil || amount <= 0 {
		logger.Warn().
			Str("amount", fmt.Sprintf("%v", body.Amount)).
			Msg("pakasir webhook received invalid amount")
		return false
	}

	detail, err := p.transactionDetail(payload.PaymentConfig, body.OrderID, amount)
	if err != nil {
		logger.Error().
			Err(err).
			Str("order_id", body.OrderID).
			Msg("failed to confirm pakasir transaction detail")
		return false
	}

	if detail == nil || detail.Transaction.Status != "completed" {
		logger.Warn().
			Str("order_id", body.OrderID).
			Str("status", detail.Transaction.Status).
			Msg("pakasir webhook transaction not confirmed as completed")
		return false
	}

	event := &events.PaymentWebhookEvent{
		ID:            body.OrderID,
		ExternalID:    body.OrderID,
		Status:        strings.ToLower(body.Status),
		Amount:        amount,
		Currency:      "IDR",
		Gateway:       "Pakasir",
		Description:   "PakKasir payment via Invoice Broker",
		PaymentConfig: payload.PaymentConfig,
		Raw:           json.RawMessage(payload.RawBody),
	}

	if err := events.Emit(event); err != nil {
		logger.Error().
			Err(err).
			Str("order_id", body.OrderID).
			Msg("failed to emit event")
		return false
	}

	logger.Info().
		Str("order_id", body.OrderID).
		Msg("webhook validated and queued for forwarding")
	return true
}

func (p *PakasirService) transactionDetail(cfg *configs.PaymentConfig, orderID string, amount float64) (*transactionDetailResponse, error) {
	req, err := http.NewRequest(http.MethodGet, pakasirAPIBase+"/api/transactiondetail", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("project", cfg.Project)
	q.Set("amount", strconv.Itoa(int(math.Round(amount))))
	q.Set("order_id", orderID)
	q.Set("api_key", cfg.APIKey)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pakasir transactiondetail returned status %d", resp.StatusCode)
	}

	var detail transactionDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

func parseAmount(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		return strconv.ParseFloat(t, 64)
	default:
		return 0, fmt.Errorf("unsupported amount type %T", v)
	}
}
