package jobqueue

import (
	"fmt"
)

// WorkRequest encapsulates the job requested to be run
type WorkRequest struct {
	JobID int64
}

//WorkQueue is a global
//var WorkQueue = make(chan WorkRequest)

// Worker does work
type Worker struct {
	ID        int
	WorkQueue chan WorkRequest
	Quit      chan bool
}

// Start starts a worker
func (w *Worker) Start() {
	fmt.Printf("starting worker %v \n", w.ID)
	go func() {
		for {
			select {
			case work := <-w.WorkQueue:
				// Receive a work request.
				fmt.Printf("worker%d: Getting Job, %v!\n", w.ID, work.JobID)

			case <-w.Quit:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// JobQueue represents a new job queue
type JobQueue struct {
	WorkQueue chan WorkRequest
	Workers   []Worker
}

// NewJobQueue provides a new Job queue with a number of workers
func (j *JobQueue) NewJobQueue(n int) JobQueue {
	wrc := make(chan WorkRequest)
	jc := JobQueue{
		WorkQueue: wrc,
	}
	for i := 0; i < n; i++ {
		jc.RegisterNewWorker(i + 1)
	}
	return jc
}

// RegisterNewWorker registers new workers
func (j *JobQueue) RegisterNewWorker(id int) {
	// Create, and return the worker.
	worker := Worker{
		ID:        id,
		WorkQueue: j.WorkQueue,
		Quit:      make(chan bool),
	}
	worker.Start()
	j.Workers = append(j.Workers, worker)
	return
}

// Push adds work to the JobQueue
func (j *JobQueue) Push(w WorkRequest) {
	j.WorkQueue <- w
}
