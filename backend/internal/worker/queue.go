package worker

import (
	"sync"
)

type JobQueue struct {
	jobs chan uint // video IDs
	mu   sync.Mutex
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: make(chan uint, 100), // Buffer for 100 jobs
	}
}

func (q *JobQueue) Push(videoID uint) {
	q.jobs <- videoID
}

func (q *JobQueue) Pop() uint {
	return <-q.jobs
}

func (q *JobQueue) Jobs() <-chan uint {
	return q.jobs
}
