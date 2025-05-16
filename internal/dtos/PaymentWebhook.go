package dtos

import "raznar.id/invoice-broker/internal/pkg/constants"

type WebhookInvoiceData struct {
	InvoiceId   string
	ExternalId  string
	Amount      float64
	Fee         float64
	Gateway     constants.Gateway
	Currency    string
	Status      string
	URL         string
	Description string
}
