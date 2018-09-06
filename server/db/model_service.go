package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

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
			k := item.Key()
			// badger creates an additional key called job1 since we have that as the
			// monotonic increasing prefix - I don't know why it does and should
			// hunt this down more, but for now, we can just not try to unmarshal
			// this one
			// TODO: followup on incrementing
			if bytes.Compare(k, []byte("job1")) == 0 {
				continue
			}
			v, err := item.Value()
			if err != nil {
				fmt.Println("error")
				fmt.Println(err)
				// TODO: do something better
				continue
			}
			var job server.Job
			err = internal.UnmarshalJob(v, &job)
			if err != nil {
				fmt.Println("error unmarshalling")
				fmt.Println(item, k, v)
				fmt.Println(err)
				continue
			} else {
				jobs = append(jobs, job)
			}
		}
		return nil
	})
	return jobs, nil
}

// GetJobsByStatus returns all jobs in the db
func (m *JobService) GetJobsByStatus(status string) ([]server.Job, error) {
	var jobs []server.Job
	err := m.client.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// jobs bucket created when db initialized
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			// badger creates an additional key called job1 since we have that as the
			// monotonic increasing prefix - I don't know why it does and should
			// hunt this down more, but for now, we can just not try to unmarshal
			// this one
			// TODO: followup on incrementing
			if bytes.Compare(k, []byte("job1")) == 0 {
				continue
			}
			v, err := item.Value()
			if err != nil {
				fmt.Println("error")
				fmt.Println(err)
				// TODO: do something better
				continue
			}
			var job server.Job
			err = internal.UnmarshalJob(v, &job)
			if err != nil {
				fmt.Println("error unmarshalling")
				fmt.Println(item, k, v)
				fmt.Println(err)
				continue
			} else {
				if job.Status == status {
					jobs = append(jobs, job)
				}
			}
		}
		return nil
	})
	if len(jobs) == 0 {
		err = errors.New("no jobs found with status: " + status)
	}
	return jobs, err
}

// GetJobByID returns details about a specific Job
func (m *JobService) GetJobByID(jobID uint64) (server.Job, error) {
	var job server.Job
	err := m.client.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// jobs bucket created when db initialized
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			// badger creates an additional key called job1 since we have that as the
			// monotonic increasing prefix - I don't know why it does and should
			// hunt this down more, but for now, we can just not try to unmarshal
			// this one
			// TODO: followup on incrementing
			if bytes.Compare(k, []byte("job1")) == 0 {
				continue
			}
			v, err := item.Value()
			if err != nil {
				fmt.Println("error")
				fmt.Println(err)
				// TODO: do something better
				continue
			}
			var newjob server.Job
			err = internal.UnmarshalJob(v, &newjob)
			if err != nil {
				fmt.Println("error unmarshalling")
				fmt.Println(item, k, v)
				fmt.Println(err)
				continue
			}
			if newjob.ID == jobID {
				job = newjob
				return nil
			}
		}
		return nil
	})
	if job.ID == uint64(0) {
		err = errors.New("job ID not found: " + strconv.Itoa(int((jobID))))
	}
	return job, err
}

// CreateJob adds a job to the db
func (m *JobService) CreateJob(job *server.Job) error {
	seq, err := m.client.db.GetSequence([]byte("job1"), 1)
	defer seq.Release()
	id, err := seq.Next()
	if err != nil {
		// TODO: handle error better
		fmt.Println("error creating sequence")
		return nil
	}
	// don't want to ever have a 0 id since makes it hard to tell
	// if the db was actually storing the job, or if the job ID was
	// set to default 0 value
	job.ID = uint64(id + 1)
	buf, err := internal.MarshalJob(job)
	if err != nil {
		fmt.Println("error marshalling")
		return err
	}
	err = m.client.db.Update(func(txn *badger.Txn) error {
		err = txn.Set(uint64ToBytes(id), buf)
		if err != nil {
			// TODO: handle error
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		// TODO: handle error better
		fmt.Println(err)
		return err
	}
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
	buf, err := internal.MarshalJob(job)
	if err != nil {
		fmt.Println("error marshalling")
		return err
	}
	err = m.client.db.Update(func(txn *badger.Txn) error {
		// job id is one less than key since never want key of 0
		err = txn.Set(uint64ToBytes(job.ID-1), buf)
		if err != nil {
			// TODO: handle error
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		// TODO: handle error better
		fmt.Println(err)
		return err
	}
	return nil
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
