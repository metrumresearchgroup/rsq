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
	JobID uint64
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
	js          server.JobService
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
			case jid := <-w.WorkQueue:
				// Receive a work request.
				lg.WithFields(logrus.Fields{
					"WID": w.ID,
					"JID": jid.JobID,
				}).Debug("received work request")
				work, err := w.js.GetJobByID(jid.JobID)
				lg.WithFields(logrus.Fields{
					"WID": w.ID,
					"JID": jid.JobID,
					"job": work,
				}).Debug("got job")
				if err != nil {
					lg.WithFields(logrus.Fields{
						"error": err,
						"JID":   jid.JobID,
						"job":   work,
					}).Error("error getting job, aborting...")
					continue
				}
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
				w.UpdateQueue <- &JobUpdate{
					Job:          work,
					msg:          "starting job",
					ShouldUpdate: true,
				}

				lg.WithFields(logrus.Fields{
					"WID": w.ID,
					"JID": jid.JobID,
					"job": work,
				}).Debug("starting Rscript")
				result, err, exitCode := runner.RunRscript(appFS, rs, es, lg)
				work.RunDetails.EndTime = time.Now().UTC()
				if err != nil {
					work.RunDetails.Error = err.Error()
				}
				lg.WithFields(logrus.Fields{
					"WID":      w.ID,
					"JID":      work.ID,
					"Duration": work.RunDetails.EndTime.Sub(work.RunDetails.StartTime),
				}).Debug("completed job")
				work.Result.Output = result
				work.Result.ExitCode = int32(exitCode)
				if exitCode == 0 {
					work.Status = "COMPLETED"
				} else {
					work.Status = "ERROR"
				}
				w.UpdateQueue <- &JobUpdate{
					Job:          work,
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
func NewJobQueue(js server.JobService, n int, updateFunc func(*server.Job), lg *logrus.Logger) *JobQueue {
	wrc := make(chan *WorkRequest, 2000)
	uq := make(chan *JobUpdate, 50)
	jc := JobQueue{
		WorkQueue:   wrc,
		UpdateQueue: uq,
	}
	for i := 0; i < n; i++ {
		jc.RegisterNewWorker(i+1, js, lg)
	}
	go jc.HandleUpdates(updateFunc)
	return &jc
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
func (j *JobQueue) RegisterNewWorker(id int, js server.JobService, lg *logrus.Logger) {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkQueue:   j.WorkQueue,
		UpdateQueue: j.UpdateQueue,
		Quit:        make(chan bool),
		js:          js,
	}
	worker.Start(lg)
	j.Workers = append(j.Workers, worker)
	return
}

// Push adds work to the JobQueue
func (j *JobQueue) Push(w uint64) {
	j.WorkQueue <- &WorkRequest{JobID: w}
}
