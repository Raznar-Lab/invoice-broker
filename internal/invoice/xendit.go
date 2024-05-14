package invoice

import (
	"slices"

	"raznar.id/invoice-broker/config"
)

type XenditHeader struct {
	CallbackToken string `json:"callback_token"`
	WebhookID     string `json:"webhook_id"`
}

func IsCBValid(payment config.PaymentConfig, header XenditHeader) (statusCode int) {
	statusCode = 200
	if !slices.Contains(payment.WebhookTokens, header.CallbackToken) {
		statusCode = 401
		return
	}

	return
}
