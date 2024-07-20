package invoice

import (
	"slices"

	"raznar.id/invoice-broker/configs"
)

type XenditHeader struct {
	CallbackToken string `json:"callback_token"`
	WebhookID     string `json:"webhook_id"`
}

func IsCBValid(payment configs.PaymentConfig, header XenditHeader) bool {
	
	return slices.Contains(payment.WebhookTokens, header.CallbackToken)
}
