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
	Job *server.Job
}

// JobUpdate provides information about the Job in the queue
type JobUpdate struct {
	Job          *server.Job
	msg          string
	err          error
	ShouldUpdate bool
}

// Worker does work
type Worker struct {
	ID          int
	WorkQueue   <-chan *WorkRequest
	UpdateQueue chan<- *JobUpdate
	Quit        chan bool
}

// JobQueue represents a new job queue
type JobQueue struct {
	WorkQueue   chan *WorkRequest
	UpdateQueue chan *JobUpdate
	Workers     []Worker
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
					"JID": work.Job.ID,
				}).Debug("getting job")
				rs := runner.RSettings{
					Rpath:   work.Job.Rscript.RPath,
					EnvVars: work.Job.Rscript.Renv,
				}
				es := runner.ExecSettings{
					WorkDir: work.Job.Rscript.WorkDir,
					Rfile:   work.Job.Rscript.RscriptPath,
				}
				work.Job.Status = "RUNNING"
				work.Job.RunDetails.StartTime = time.Now().UTC()
				w.UpdateQueue <- &JobUpdate{
					Job:          work.Job,
					msg:          "starting job",
					ShouldUpdate: true,
				}
				result, err, exitCode := runner.RunRscript(appFS, rs, es, lg)
				work.Job.RunDetails.EndTime = time.Now().UTC()
				work.Job.RunDetails.Error = err.Error()
				lg.WithFields(logrus.Fields{
					"WID":      w.ID,
					"JID":      work.Job.ID,
					"Duration": work.Job.RunDetails.EndTime.Sub(work.Job.RunDetails.StartTime),
				}).Debug("completed job")
				work.Job.Result.Output = result
				work.Job.Result.ExitCode = int32(exitCode)
				if exitCode == 0 {
					work.Job.Status = "COMPLETED"
				} else {
					work.Job.Status = "ERROR"
				}
				w.UpdateQueue <- &JobUpdate{
					Job:          work.Job,
					msg:          "completed job",
					err:          err,
					ShouldUpdate: true,
				}

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

// NewJobQueue provides a new Job queue with a number of workers
func NewJobQueue(n int, updateFunc func(*server.Job), lg *logrus.Logger) JobQueue {
	wrc := make(chan *WorkRequest, 2000)
	uq := make(chan *JobUpdate, 50)
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
func (j *JobQueue) HandleUpdates(fn func(*server.Job)) {
	for {
		ju := <-j.UpdateQueue
		if ju.ShouldUpdate {
			fn(ju.Job)
		}
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
func (j *JobQueue) Push(w *server.Job) {
	j.WorkQueue <- &WorkRequest{Job: w}
}
