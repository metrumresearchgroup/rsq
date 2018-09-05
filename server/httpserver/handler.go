package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/metrumresearchgroup/rsq/server"
)

type key int

const (
	keyJobID key = iota
)

// JobHandler represents the HTTP API handler for JobService
type JobHandler struct {
	JobService server.JobService
}

// NewJobHandler provides a pointer to a new httpClient
func NewJobHandler(js server.JobService) *JobHandler {
	return &JobHandler{
		JobService: js,
	}
}

// HandleGetJobsByStatus provides all jobs
// accepts query param status with values COMPLETED, QUEUED, RUNNING
func (c *JobHandler) HandleGetJobsByStatus(w http.ResponseWriter, r *http.Request) {
	var jobs []server.Job
	status := r.URL.Query().Get("status")
	fmt.Println("status: ", status)
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

// JobCtx is the context
func (c *JobHandler) JobCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "jobID")
		mid, _ := strconv.ParseInt(jobID, 10, 64)
		job, err := c.JobService.GetJobByID(mid)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), keyJobID, &job)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandleSubmitJob adds a job to the database for workers to execute
func (c *JobHandler) HandleSubmitJob(w http.ResponseWriter, r *http.Request) {
	var job server.Job
	if err := render.DecodeJSON(r.Body, &job); err != nil {
		render.JSON(w, r, err.Error())
		return
	}
	err := c.JobService.CreateJob(&job)
	if err != nil {
		fmt.Printf("Insertion of jobs failed with err: %v", err)
	}
	render.JSON(w, r, job)
}
