package workers

import (
	"sync"

	"raznar.id/invoice-broker/internal/jobs"
)

var (
	instance *Pool
	once     sync.Once
)

type Pool struct {
	queue chan jobs.JobWorker
}

// Init initializes the global worker pool once
func Init(bufferSize int, workerCount int) {
	once.Do(func() {
		instance = &Pool{
			queue: make(chan jobs.JobWorker, bufferSize),
		}
		instance.start(workerCount)
	})
}

// Enqueue provides a global entry point to the pool
func Enqueue(job jobs.JobWorker) bool {
	if instance == nil {
		return false
	}
	return instance.enqueue(job)
}

func (p *Pool) start(count int) {
	for i := 0; i < count; i++ {
		go func(id int) {
			for job := range p.queue {
				_ = job.Run(id)
			}
		}(i)
	}
}

func (p *Pool) enqueue(job jobs.JobWorker) bool {
	select {
	case p.queue <- job:
		return true
	default:
		return false
	}
}
