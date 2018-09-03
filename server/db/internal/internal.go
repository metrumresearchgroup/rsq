package internal

import (
	"fmt"

	"github.com/dpastoor/babylon/runner"
	"github.com/gogo/protobuf/proto"
	"github.com/metrumresearchgroup/srq/server"
)

// need go get github.com/gogo/protobuf/protoc-gen-gofast

//go:generate protoc --gofast_out=. internal.proto

// MarshalJob encodes a model to binary format.
func MarshalJob(m *server.Job) ([]byte, error) {
	var status Job_StatusType
	runDetails := &m.RunDetails
	switch m.Status {
	case "QUEUED":
		status = Job_QUEUED
	case "RUNNING":
		status = Job_RUNNING
	case "COMPLETED":
		status = Job_COMPLETED
	case "ERROR":
		status = Job_ERROR
	default:
		return nil, fmt.Errorf("unrecognized Job status: %v", m.Status)
	}
	return proto.Marshal(&Job{
		Id:     m.ID,
		Status: status,
		RunDetails: &RunDetails{
			QueueTime: jobDetails.QueueTime,
			StartTime: jobDetails.StartTime,
			Duration:  jobDetails.Duration,
			RunDir:    jobDetails.RunDir,
			Error:     jobDetails.Error,
		},
	})
}

// UnmarshalJob decodes a job from a binary data.
func UnmarshalJob(data []byte, m *server.Job) error {
	var pb Job
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	runDetails := pb.GetRunDetails()

	m.ID = pb.Id
	status := pb.GetStatus()

	switch status {
	case Job_COMPLETED:
		m.Status = "COMPLETED"
	case Job_ERROR:
		m.Status = "ERROR"
	case Job_RUNNING:
		m.Status = "RUNNING"
	case Job_QUEUED:
		m.Status = "QUEUED"
	default:
		return fmt.Errorf("unrecognized job status: %v", status)
	}
	m.RunDetails = server.RunDetails{
		QueueTime: runDetails.QueueTime,
		StartTime: runDetails.StartTime,
		Duration:  runDetails.Duration,
		Error:     runDetails.Error,
	}

	return nil
}
