package webhook_job

import (
	"bytes"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/internal/workers"
)

const (
	MaxRetries = 40
	// 10 times a day for 4 days = 144 minutes per attempt
	RetryDelay = 144 * time.Minute
)

type WebhookJob struct {
	// Removed Pool pointer - we use the workers singleton now
	Content  []byte
	URL      string
	Header   string
	Token    string
	Attempts int
}

func (j WebhookJob) Run(workerID int) error {
	logger := log.With().
		Int("worker_id", workerID).
		Str("url", j.URL).
		Int("attempt", j.Attempts+1).
		Logger()

	req, err := http.NewRequest("POST", j.URL, bytes.NewBuffer(j.Content))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create request object")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(j.Header, j.Token)

	// Keep timeout tight to free up workers quickly
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)

	// Error handling & Retry logic
	if err != nil || (res != nil && (res.StatusCode < 200 || res.StatusCode >= 300)) {
		if res != nil {
			res.Body.Close()
		}

		if j.Attempts < MaxRetries {
			j.Attempts++

			logger.Warn().
				Msgf("Delivery failed. Retrying in %v (Attempt %d/%d)", RetryDelay, j.Attempts, MaxRetries)

			// Non-blocking background timer
			time.AfterFunc(RetryDelay, func() {
				// Accessing the global Enqueue from the workers package
				if !workers.Enqueue(j) {
					log.Error().Str("url", j.URL).Msg("Retry failed: worker queue is full")
				}
			})
		} else {
			logger.Error().Msg("Max retries reached after 4 days. Abandoning webhook.")
		}
		return nil
	}

	defer res.Body.Close()
	logger.Info().Msg("Webhook delivered successfully")
	return nil
}