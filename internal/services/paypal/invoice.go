package paypal_service

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/configs"
	"raznar.id/invoice-broker/internal/constants"
)

type CreateInvoicePayload struct {
	PaymentConfig *configs.GatewayConfig
	ID            string
	Description   string
	Amount        float64
	ReturnURL     string
	CancelURL     string
}

func (p *PaypalService) CreateInvoice(payload CreateInvoicePayload) (string, error) {
	logger := log.With().
		Str("service", "paypal_create").
		Str("external_id", payload.ID).
		Logger()

	logger.Debug().
		Float64("amount", payload.Amount).
		Str("return_url", payload.ReturnURL).
		Str("cancel_url", payload.CancelURL).
		Msg("creating paypal invoice")

	fee := payload.Amount * 0.01
	total := payload.Amount + fee

	logger.Debug().
		Float64("fee", fee).
		Float64("total_amount", total).
		Msg("paypal fee calculated")

	values := url.Values{}
	values.Set("cmd", "_xclick")
	values.Set("business", payload.PaymentConfig.Paypal.APIKey)
	values.Set("item_name", strings.TrimSpace(payload.Description))

	values.Set("invoice", payload.ID)
	values.Set("amount", fmt.Sprintf("%.2f", total))
	values.Set("currency_code", "USD")
	values.Set("return", payload.ReturnURL)
	values.Set("cancel_return", payload.CancelURL)
	values.Set("no_shipping", "1")
	values.Set("rm", "2")

	endpoint := constants.PAYPAL_LIVE_ENDPOINT
	env := "live"
	if payload.PaymentConfig.Paypal.Sandbox {
		endpoint = constants.PAYPAL_SANDBOX_ENDPOINT
		env = "sandbox"
	}

	logger.Debug().
		Str("environment", env).
		Str("endpoint", endpoint).
		Msg("paypal environment selected")

	approvalURL := endpoint + "/cgi-bin/webscr?" + values.Encode()
	decodedURL, _ := url.QueryUnescape(approvalURL)
	logger.Info().Str("approval_url", decodedURL).Msg("PayPal invoice created")
	logger.Info().
		Str("approval_url", approvalURL).
		Msg("paypal invoice created (redirect only)")

	return approvalURL, nil
}
