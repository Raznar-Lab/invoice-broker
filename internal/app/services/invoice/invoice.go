package services_invoice

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/dustin/go-humanize"
	"raznar.id/invoice-broker/configs"
	services_base "raznar.id/invoice-broker/internal/app/services/base"
	"raznar.id/invoice-broker/internal/dtos"
)

type InvoiceService struct {
	services_base.BaseService
}

func (s InvoiceService) ForwardCallbackData(body []byte, urlList []string, header []string, token string) {
	for _, url := range urlList {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		for _, v := range header {
			req.Header.Set(v, token)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[WARNING] an error occured when forwarding callback on http: %s", err.Error())
		}

		if res.StatusCode != 200 && res.StatusCode != 201 {
			log.Printf("[WARNING] Unsucessful forward %s", res.Body)
		}
	}
}

func (s InvoiceService) ForwardWebhookData(data dtos.WebhookInvoiceData, webhooks map[string]configs.PaymentWebhookConfig) {
	for _, w := range webhooks {
		content := w.Content
		for vk, vv := range w.Variables {
			vv = strings.ReplaceAll(vv, "{transaction_id}", data.InvoiceId)
			vv = strings.ReplaceAll(vv, "{external_id}", data.ExternalId)
			vv = strings.ReplaceAll(vv, "{amount}", humanize.Commaf(data.Amount))
			vv = strings.ReplaceAll(vv, "{fee}", humanize.Commaf(data.Fee))
			vv = strings.ReplaceAll(vv, "{gateway}", data.Gateway)
			vv = strings.ReplaceAll(vv, "{currency}", data.Currency)
			content = strings.ReplaceAll(content, vk, vv)
		}

		req, _ := http.NewRequest("POST", w.URL, bytes.NewBuffer([]byte(content)))
		req.Header.Set("Content-Type", w.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[WARNING] an error occured when forwarding webhook on http: %s", err.Error())
		}

		if res.StatusCode != 200 && res.StatusCode != 201 {
			log.Printf("[WARNING] Unsucessful webhook %s", res.Body)
		}
	}
}
