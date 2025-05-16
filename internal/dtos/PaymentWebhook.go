package dtos

type WebhookInvoiceData struct {
	InvoiceId  string
	ExternalId string
	Amount     float64
	Fee        float64
	Gateway    string
	Currency   string
}
