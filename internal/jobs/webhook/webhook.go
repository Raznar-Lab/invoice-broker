package webhook_job

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"raznar.id/invoice-broker/internal/workers"
)

const (
	RetryDelay = 1 * time.Hour
	MaxRetries = 24 * 4 // 4 days, 96 retries
)

var activeJobs sync.Map
type WebhookJob struct {
	ID       string
	Content  []byte
	URL      string
	Headers  map[string]string
	Attempts int
}

func (j *WebhookJob) GenerateId() {
	h := sha256.New()
	h.Write([]byte(j.URL))
	h.Write(j.Content)
	j.ID = hex.EncodeToString(h.Sum(nil))
}

// ---------- enqueue ----------
func (j *WebhookJob) Enqueue() bool {
	if _, exists := activeJobs.LoadOrStore(j.ID, j); exists {
		log.Warn().
			Str("job_id", j.ID).
			Int("attempts", j.Attempts).
			Msg("Webhook job already active, skipping enqueue")
		return false
	}

	if ok := workers.Enqueue(j); !ok {
		activeJobs.Delete(j.ID)
		log.Error().
			Str("job_id", j.ID).
			Msg("Failed to enqueue webhook job")
		return false
	}

	log.Info().
		Str("job_id", j.ID).
		Int("attempts", j.Attempts).
		Msg("Webhook job enqueued successfully")

	return true
}

// ---------- worker execution ----------
func (j *WebhookJob) Run(workerID int) error {
	logger := log.With().
		Int("worker_id", workerID).
		Str("job_id", j.ID).
		Str("url", j.URL).
		Int("attempts", j.Attempts).
		Logger()

	req, err := http.NewRequest("POST", j.URL, bytes.NewBuffer(j.Content))
	if err != nil {
		activeJobs.Delete(j.ID)
		logger.Error().Err(err).Msg("Failed to create HTTP request for webhook")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range j.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)

	// ---------- failure path ----------
	if err != nil || res == nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		statusCode := 0
		if res != nil {
			statusCode = res.StatusCode
			res.Body.Close()
		}

		j.Attempts++

		if j.Attempts > MaxRetries {
			activeJobs.Delete(j.ID)
			logger.Error().
				Int("attempts", j.Attempts).
				Msg("Max retries reached, abandoning webhook")
			return nil
		}

		logger.Warn().
			Int("attempts", j.Attempts).
			Int("status_code", statusCode).
			Msg("Webhook failed, scheduling retry")

		time.AfterFunc(RetryDelay, func() {
			j.Enqueue()
			log.Debug().
				Str("job_id", j.ID).
				Time("next_retry_at", time.Now().Add(RetryDelay)).
				Msg("Scheduled webhook retry")
		})

		activeJobs.Delete(j.ID)
		return nil
	}

	// ---------- success ----------
	res.Body.Close()
	activeJobs.Delete(j.ID)

	logger.Info().
		Int("status_code", res.StatusCode).
		Int("attempts", j.Attempts).
		Msg("Webhook delivered successfully")

	return nil
}
