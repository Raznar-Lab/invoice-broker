package invoice

import (
	"bytes"
	"log"
	"net/http"
)

func ForwardWebhookData(body []byte, urlList []string, header string, token string) {
	for _, url := range urlList {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(header, token)

		_, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[WARNING] an error occured when forwarding webhook on http: %s", err.Error())
		}
	}
}
