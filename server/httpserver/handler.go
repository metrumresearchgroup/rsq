package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/metrumresearchgroup/rsq/jobqueue"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/metrumresearchgroup/rsq/server"
)

// this is because ctx cannot have raw strings
// so instead we can define a unique const that can be used to prevent
// context collisions
// see: https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// and information about context collision at https://blog.golang.org/context#TOC_3.2.
type key int

const (
	keyJobID key = iota
)

// JobHandler represents the HTTP API handler for JobService
type JobHandler struct {
	JobService server.JobService
	Queue      *jobqueue.JobQueue
	Logger     *logrus.Logger
}

// NewJobHandler provides a pointer to a new httpClient
func NewJobHandler(js server.JobService, n int, memory float64, lg *logrus.Logger) *JobHandler {
	return &JobHandler{
		Logger:     lg,
		JobService: js,
		Queue: jobqueue.NewJobQueue(js, n, memory, func(j *server.Job) {
			js.UpdateJob(j)
		}, lg),
	}
}

// HandleGetJobsByStatus provides all jobs
// accepts query param status with values COMPLETED, QUEUED, RUNNING
func (c *JobHandler) HandleGetJobsByStatus(w http.ResponseWriter, r *http.Request) {
	var jobs []*server.Job
	status := r.URL.Query().Get("status")
	if status != "" {
		jobs, _ = c.JobService.GetJobsByStatus(status)
	} else {
		jobs, _ = c.JobService.GetJobs()
	}
	render.JSON(w, r, jobs)
}

// HandleGetJobByID handles providing job by ID taken from context
func (c *JobHandler) HandleGetJobByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	job, ok := ctx.Value(keyJobID).(*server.Job)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	render.JSON(w, r, job)
}

// HandleCancelJob
func (c *JobHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	job, ok := ctx.Value(keyJobID).(*server.Job)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	cancelled, err := c.JobService.CancelJob(job.ID)
	if err != nil {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	if cancelled {
		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, true)
	} else {
		w.WriteHeader(http.StatusNotModified)
		render.JSON(w, r, false)
	}

}

// JobCtx is the context
func (c *JobHandler) JobCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobID")
		mid, _ := strconv.ParseUint(jobID, 10, 64)
		job, err := c.JobService.GetJobByID(mid)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), keyJobID, job)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandleSubmitJob adds a job to the database for workers to execute
func (c *JobHandler) HandleSubmitJob(w http.ResponseWriter, r *http.Request) {
	var job server.Job
	if err := render.DecodeJSON(r.Body, &job); err != nil {
		c.Logger.WithFields(logrus.Fields{
			"body": r.Body,
			"err":  err,
		}).Warn("Decoding JSON from job submission failed")
		render.JSON(w, r, fmt.Sprintf("Error decoding JSON %s", r.Body))
		return
	}
	job.RunDetails.QueueTime = time.Now().UTC()
	err := c.JobService.CreateJob(&job)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Insertion of jobs failed")
		render.JSON(w, r, fmt.Sprintf("error inserting job"))
		return
	}
	c.Queue.Push(job.ID)
	render.JSON(w, r, job)
}
