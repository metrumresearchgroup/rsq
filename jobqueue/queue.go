package jobqueue

import (
	"time"

	"github.com/metrumresearchgroup/rsq/runner"
	"github.com/metrumresearchgroup/rsq/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// WorkRequest encapsulates the job requested to be run
type WorkRequest struct {
	JobID int64
}

// Worker does work
type Worker struct {
	ID          int
	WorkQueue   <-chan server.Job
	UpdateQueue chan<- server.Job
	Quit        chan bool
}

// Start starts a worker
func (w *Worker) Start(lg *logrus.Logger) {
	lg.Infof("starting worker %v \n", w.ID)
	go func() {
		appFS := afero.NewOsFs()
		for {
			select {
			case work := <-w.WorkQueue:
				// Receive a work request.
				lg.WithFields(logrus.Fields{
					"WID": w.ID,
					"JID": work.ID,
				}).Debug("getting job")
				rs := runner.RSettings{
					Rpath:   work.Rscript.RPath,
					EnvVars: work.Rscript.Renv,
				}
				es := runner.ExecSettings{
					WorkDir: work.Rscript.WorkDir,
					Rfile:   work.Rscript.RscriptPath,
				}
				work.Status = "RUNNING"
				work.RunDetails.StartTime = time.Now().UTC()
				w.UpdateQueue <- work
				result, _ := runner.RunRscript(appFS, rs, es, lg)
				work.RunDetails.EndTime = time.Now().UTC()
				lg.WithFields(logrus.Fields{
					"WID":      w.ID,
					"JID":      work.ID,
					"Duration": work.RunDetails.EndTime.Sub(work.RunDetails.StartTime),
				}).Debug("completed job")
				work.Result.Output = result
				work.Status = "COMPLETED"
				work.Result.ExitCode = 0
				w.UpdateQueue <- work

			case <-w.Quit:
				// We have been asked to stop.
				lg.Printf("worker%d stopping\n", w.ID)
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
	WorkQueue   chan server.Job
	UpdateQueue chan server.Job
	Workers     []Worker
}

// NewJobQueue provides a new Job queue with a number of workers
func NewJobQueue(n int, updateFunc func(server.Job), lg *logrus.Logger) JobQueue {
	wrc := make(chan server.Job, 200)
	uq := make(chan server.Job, 5)
	jc := JobQueue{
		WorkQueue:   wrc,
		UpdateQueue: uq,
	}
	for i := 0; i < n; i++ {
		jc.RegisterNewWorker(i+1, lg)
	}
	go jc.HandleUpdates(updateFunc)
	return jc
}

// HandleUpdates handles updates
func (j *JobQueue) HandleUpdates(fn func(server.Job)) {
	for {
		fn(<-j.UpdateQueue)
	}
}

// RegisterNewWorker registers new workers
func (j *JobQueue) RegisterNewWorker(id int, lg *logrus.Logger) {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkQueue:   j.WorkQueue,
		UpdateQueue: j.UpdateQueue,
		Quit:        make(chan bool),
	}
	worker.Start(lg)
	j.Workers = append(j.Workers, worker)
	return
}

// Push adds work to the JobQueue
func (j *JobQueue) Push(w server.Job) {
	j.WorkQueue <- w
}
