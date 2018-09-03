package db

import (
	"encoding/binary"

	"github.com/dgraph-io/badger"
	"github.com/metrumresearchgroup/rsq/server"
	"github.com/metrumresearchgroup/rsq/server/db/internal"
)

// make sure JobService implements server.JobService
var _ server.JobService = &JobService{}

// JobService represents a service for managing jobs
type JobService struct {
	client *Client
}

// GetJobs returns all jobs in the db
func (m *JobService) GetJobs() ([]server.Job, error) {
	var jobs []server.Job
	m.client.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// jobs bucket created when db initialized
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			//k := item.Key()
			v, err := item.Value()
			if err != nil {
				// TODO: do something
			}
			var job server.Job
			internal.UnmarshalJob(v, &job)
			jobs = append(jobs, job)
		}
		return nil
	})
	return jobs, nil
}

// GetJobsByStatus returns all jobs in the db
func (m *JobService) GetJobsByStatus(status string) ([]server.Job, error) {
	var jobs []server.Job
	return jobs, nil
}

// GetJobByID returns details about a specific Job
func (m *JobService) GetJobByID(jobID int) (server.Job, error) {
	var job server.Job
	return job, nil
}

// CreateJob adds a job to the db
func (m *JobService) CreateJob(job *server.Job) error {
	return nil
}

// CreateJobs adds an array of jobs to the db in a single batch transaction
func (m *JobService) CreateJobs(jobs []server.Job) ([]server.Job, error) {
	return jobs, nil
}

// AcquireNextQueuedJob returns the next job with status QUEUED while also changing the value to RUNNING
func (m *JobService) AcquireNextQueuedJob() (server.Job, error) {
	var nextJob server.Job
	return nextJob, nil
}

// UpdateJob updates the job status
func (m *JobService) UpdateJob(job *server.Job) error {
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
