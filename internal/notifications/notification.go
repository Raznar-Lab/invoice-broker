package notifications

import "raznar.id/invoice-broker/internal/notifications/webhook"

type Payload struct {
	Content any
	URLS    []string
	Headers map[string]string
}

func Send(payload Payload) error {
	err := webhook.Send(webhook.Payload{
		Content: payload.Content,
		URLS:    payload.URLS,
		Headers: payload.Headers,
	})
	return err
}
