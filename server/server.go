package server

import "time"

// RunDetails stores details about a script being run
// Queue time represents the time a request was added to the Queue
// StartTime is the time the worker starts execution of the code for processing steps
// Duration is the time, in milliseconds from StartTime to the run completing
// RunDir is the (sub)-directory where the script was executed
// Error is the string representation of the error that stopped the run if an error was present
// as a unix timestamp
type RunDetails struct {
	QueueTime time.Time `json:"queue_time,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Job represents information about the job queue
type Job struct {
	ID         uint64
	Status     string
	RunDetails RunDetails
	// some information about the job like the title
	Context string
}

// Client creates a connection to services
type Client interface {
	JobQueueService() JobService
}

// JobService describes the interface to interact with models
type JobService interface {
	GetJobs() ([]Job, error)
	GetJobsByStatus(status string) ([]Job, error)
	GetJobByID(jobID uint64) (Job, error)
	CreateJob(m *Job) error
	CreateJobs(job []Job) ([]Job, error)
	AcquireNextQueuedJob() (Job, error)
	UpdateJob(m *Job) error
}
