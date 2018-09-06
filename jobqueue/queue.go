package jobqueue

import (
	"fmt"
	"time"
)

// WorkRequest encapsulates the job requested to be run
type WorkRequest struct {
	JobID int64
}

// Worker does work
type Worker struct {
	ID          int
	WorkQueue   chan WorkRequest
	UpdateQueue chan WorkRequest
	Quit        chan bool
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
				time.Sleep(time.Duration(1 * time.Second))
				w.UpdateQueue <- work

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
	WorkQueue   chan WorkRequest
	UpdateQueue chan WorkRequest
	Workers     []Worker
}

// NewJobQueue provides a new Job queue with a number of workers
func NewJobQueue(n int, updateFunc func(WorkRequest)) JobQueue {
	wrc := make(chan WorkRequest)
	uq := make(chan WorkRequest)
	jc := JobQueue{
		WorkQueue:   wrc,
		UpdateQueue: uq,
	}
	for i := 0; i < n; i++ {
		jc.RegisterNewWorker(i + 1)
	}
	go jc.HandleUpdates(updateFunc)
	return jc
}

// HandleUpdates handles updates
func (j *JobQueue) HandleUpdates(fn func(WorkRequest)) {
	for {
		fn(<-j.UpdateQueue)
	}
}

// RegisterNewWorker registers new workers
func (j *JobQueue) RegisterNewWorker(id int) {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkQueue:   j.WorkQueue,
		UpdateQueue: j.UpdateQueue,
		Quit:        make(chan bool),
	}
	worker.Start()
	j.Workers = append(j.Workers, worker)
	return
}

// Push adds work to the JobQueue
func (j *JobQueue) Push(w WorkRequest) {
	j.WorkQueue <- w
}
