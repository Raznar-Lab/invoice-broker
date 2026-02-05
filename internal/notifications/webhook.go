package notifications

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	webhook_job "raznar.id/invoice-broker/internal/jobs/webhook"
)

type WebhookPayload struct {
	Content any
	URLS    []string
	Header  string
	Token   string
}

func SendWebhook(payload WebhookPayload) error {
	body, err := json.Marshal(payload.Content)
	if err != nil {
		log.Debug().
			Err(err).
			Interface("content", payload.Content).
			Msg("Failed to marshal webhook content")
		return err
	}

	preview := string(body)
	if len(preview) > 200 {
		preview = preview[:200] + "..." // truncate for debug
	}

	log.Debug().
		Int("url_count", len(payload.URLS)).
		Str("content_preview", preview).
		Str("header", payload.Header).
		Msg("Dispatching webhook jobs")

	for i, url := range payload.URLS {
		log.Debug().
			Int("index", i).
			Str("url", url).
			Str("header", payload.Header).
			Msg("Enqueuing webhook job")

		job := webhook_job.New(webhook_job.Payload{
			Content: body,
			URL:     url,
			Header:  payload.Header,
			Token:   payload.Token,
		})
		ok := job.Enqueue()
		if ok {
			log.Info().
				Str("url", url).
				Msg("Failed to enqueue webhook job")
		} else {
			log.Info().
				Str("url", url).
				Msg("Webhook job enqueued successfully")
		}
	}

	log.Debug().
		Int("total_urls", len(payload.URLS)).
		Msg("Finished dispatching webhook jobs")

	return nil
}
