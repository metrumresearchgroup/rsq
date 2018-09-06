package server

import "time"

// RScript contains information needed to run an Rscript
type Rscript struct {
	RPath       string            `json:"r_path,omitempty"`
	WorkDir     string            `json:"work_dir,omitempty"`
	RscriptPath string            `json:"rscript_path,omitempty"`
	Renv        map[string]string `json:"renv,omitempty"`
}

// Result stores the results from Rscript
type Result struct {
	Output   string `json:"output,omitempty"`
	ExitCode int32  `json:"exit_code,omitempty"`
}

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
	Rscript Rscript
	Result  Result
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
