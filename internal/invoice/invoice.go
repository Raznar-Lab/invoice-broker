package invoice

import (
	"bytes"
	"log"
	"net/http"
)

func ForwardWebhookData(body []byte, urlList []string) {
	for _, url := range urlList {
		_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("[WARNING] an error occured when forwarding webhook on http: %s", err.Error())
		}
	}
}