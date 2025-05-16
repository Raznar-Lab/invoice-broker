package configs

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"raznar.id/invoice-broker/internal/dtos"
)

type GatewayConfig struct {
	Xendit PaymentConfig `yaml:"xendit"`
	Paypal PaymentConfig `yaml:"paypal"`
}

type PaymentConfig struct {
	APIKey         string                          `yaml:"api_key"`
	Label          string                          `yaml:"label"`
	Webhooks       map[string]PaymentWebhookConfig `yaml:"webhooks"`
	CallbackTokens []string                        `yaml:"callback_tokens"`
	CallbackURLS   []string                        `yaml:"callback_urls"`
}

type PaymentWebhookConfig struct {
	URL         string            `yaml:"url"`
	ContentType string            `yaml:"content_type"`
	Content     string            `yaml:"content"`
	Variables   map[string]string `yaml:"variables"`
}

func (w PaymentWebhookConfig) ApplyVariables(content string, data dtos.WebhookInvoiceData) string {
	content = strings.ReplaceAll(content, "{invoice_id}", data.InvoiceId)
	content = strings.ReplaceAll(content, "{external_id}", data.ExternalId)
	content = strings.ReplaceAll(content, "{amount}", humanize.Commaf(data.Amount))
	content = strings.ReplaceAll(content, "{fee}", humanize.Commaf(data.Fee))
	content = strings.ReplaceAll(content, "{gateway}", fmt.Sprintf("%s", data.Gateway))
	content = strings.ReplaceAll(content, "{currency}", data.Currency)
	content = strings.ReplaceAll(content, "{status}", data.Status)
	content = strings.ReplaceAll(content, "{url}", data.URL)
	content = strings.ReplaceAll(content, "{description}", data.Description)

	return content
}
