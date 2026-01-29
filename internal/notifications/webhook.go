package notifications

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	webhook_job "raznar.id/invoice-broker/internal/jobs/webhook"
)

func SendWebhook(content any, urlList []string, header string, token string) error {
	body, err := json.Marshal(content)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("Failed to marshal webhook content")
		return err
	}

	log.Debug().
		Int("url_count", len(urlList)).
		Msg("Dispatching webhook jobs")

	for _, url := range urlList {
		log.Debug().
			Str("url", url).
			Str("header", header).
			Msg("Enqueuing webhook job")

		job := webhook_job.New(webhook_job.Payload{
			Content: body,
			URL:     url,
			Header:  header,
			Token:   token,
		})

		job.Enqueue()
	}

	return nil
}
