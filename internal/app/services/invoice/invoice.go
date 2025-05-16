package services_invoice

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		content := w.ApplyVariables(w.Content, data)
		for vk, vv := range w.Variables {
			vv = w.ApplyVariables(vv, data)
			content = strings.ReplaceAll(content, fmt.Sprintf("{%s}", vk), vv)
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
