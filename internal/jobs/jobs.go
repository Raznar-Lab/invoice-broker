package jobs

// JobWorker is the interface that all background tasks must implement
type JobWorker interface {
	Run(workerID int) error
}

// Job is the interface that all background tasks must implement
type Job interface {
	Enqueue() bool
	GenerateId()
}


func Enqueue(job Job) bool {
	job.GenerateId()
	return job.Enqueue()
}
