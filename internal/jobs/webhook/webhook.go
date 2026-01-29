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
	// Retry once per hour for 4 days (96 retries)
	MaxRetries = 24 * 4
)

// job registry (ID -> *WebhookJob)
var activeJobs sync.Map

type Payload struct {
	Content []byte
	URL     string
	Header  string
	Token   string
}

func (p Payload) ID() string {
	h := sha256.New()
	h.Write([]byte(p.URL))
	h.Write(p.Content)
	return hex.EncodeToString(h.Sum(nil))
}

type WebhookJob struct {
	ID       string
	Content  []byte
	URL      string
	Header   string
	Token    string
	Attempts int
}

// ---------- constructor ----------

func New(payload Payload) *WebhookJob {
	return &WebhookJob{
		ID:      payload.ID(),
		Content: payload.Content,
		URL:     payload.URL,
		Header:  payload.Header,
		Token:   payload.Token,
	}
}

// ---------- enqueue ----------

func (j *WebhookJob) Enqueue() bool {
	// Deduplicate by job ID
	if _, exists := activeJobs.LoadOrStore(j.ID, j); exists {
		log.Warn().
			Str("job_id", j.ID).
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
		Int("attempt", j.Attempts+1).
		Msg("Webhook job enqueued")

	return true
}

// ---------- worker execution ----------

func (j *WebhookJob) Run(workerID int) error {
	logger := log.With().
		Int("worker_id", workerID).
		Str("job_id", j.ID).
		Str("url", j.URL).
		Int("attempt", j.Attempts+1).
		Logger()

	req, err := http.NewRequest("POST", j.URL, bytes.NewBuffer(j.Content))
	if err != nil {
		activeJobs.Delete(j.ID)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(j.Header, j.Token)

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)

	// ---------- failure path ----------
	if err != nil || (res.StatusCode >= 300 || res.StatusCode < 200) {
		if res != nil {
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

		// IMPORTANT: remove before retry
		activeJobs.Delete(j.ID)

		logger.Warn().
			Int("attempts", j.Attempts).
			Int("code", res.StatusCode).
			Msg("Webhook failed, scheduling retry")

		time.AfterFunc(RetryDelay, func() {
			j.Enqueue()
		})

		return nil
	}

	// ---------- success ----------
	res.Body.Close()
	activeJobs.Delete(j.ID)

	logger.Info().Int("code", res.StatusCode).
		Msg("Webhook delivered successfully")

	return nil
}
