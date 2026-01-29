package notifications

import (
	"encoding/json"

	webhook_job "raznar.id/invoice-broker/internal/jobs/webhook"
	"raznar.id/invoice-broker/internal/workers"
)

func SendWebhook(content any, urlList []string, header string, token string) error {
	body, _ := json.Marshal(content)

	for _, url := range urlList {
		job := webhook_job.WebhookJob{
			// Note: The job struct still needs a way to re-enqueue itself
			// We can modify the job to use workers.Enqueue() instead of j.Pool.Enqueue()
			Content: body,
			URL:     url,
			Header:  header,
			Token:   token,
		}

		workers.Enqueue(job)
	}
	return nil
}