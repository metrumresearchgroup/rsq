package db

import (
	"log"

	"os"

	"github.com/dgraph-io/badger"
	"github.com/metrumresearchgroup/rsq/server"
)

// Client represents a client to the underlying BoltDB instance
type Client struct {
	// Filepath to the BoltDB database
	Path string

	// Services
	jobService JobService

	db *badger.DB
}

// JobService provides an interface for getting jobs
// type JobService interface {
// 	GetJobs() ([]*Job, error)
// 	GetJob(mID int) (*Job, error)
// 	CreateJob(md Job) error
// }

// NewClient creates a new client bound to the jobService
func NewClient() *Client {
	c := &Client{}
	c.jobService.client = c
	return c
}

// Open opens and initializes the Badger database
func (c *Client) Open() error {
	// Open database file.
	path := c.Path
	if path == "" {
		path = "/tmp/badger"
	}

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	c.db = db
	return err
}

// Close closes then underlying BoltDB database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ResetDb will delete and reinitialize the DB
func (c *Client) ResetDb() error {
	// if file exists delete
	if _, err := os.Stat(c.Path); err == nil {
		err = os.Remove(c.Path)
		if err != nil {
			return err
		}
	}
	c.Open()
	return nil
}

//JobService returns the jobService associated with the client
func (c *Client) JobService() server.JobService {
	return &c.jobService
}
