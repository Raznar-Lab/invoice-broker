package webhook

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/internal/jobs"
	webhook_job "raznar.id/invoice-broker/internal/jobs/webhook"
)

type Payload struct {
	Content any
	URLS    []string
	Headers map[string]string
}

func Send(payload Payload) error {
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
		Interface("header", payload.Headers).
		Msg("Dispatching webhook jobs")

	for i, url := range payload.URLS {
		log.Debug().
			Int("index", i).
			Str("url", url).
			Interface("header", payload.Headers).
			Msg("Enqueuing webhook job")

		job := &webhook_job.WebhookJob{
			Content: body,
			URL:     url,
			Headers: payload.Headers,
		}
		
		jobs.Enqueue(job)
	}

	log.Debug().
		Int("total_urls", len(payload.URLS)).
		Msg("Finished dispatching webhook jobs")

	return nil
}
