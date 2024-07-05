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

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[WARNING] an error occured when forwarding webhook on http: %s", err.Error())
		}

		if res.StatusCode != 200 && res.StatusCode != 201 {
			log.Printf("[WARNING] Unsucessful forward %s", res.Body)
		}
	}
}
