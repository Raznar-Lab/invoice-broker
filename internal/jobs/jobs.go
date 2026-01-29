package jobs

// Job is the interface that all background tasks must implement
type Job interface {
	Run(workerID int) error
}