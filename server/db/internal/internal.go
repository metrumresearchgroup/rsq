package internal

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/metrumresearchgroup/rsq/server"
)

// using to gogoproto had problems with timestamps
// so just using regular as super performance doesn't matter

//go:generate protoc --go_out=. internal.proto

// MarshalJob encodes a model to binay format.
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

	queueTime, _ := ptypes.TimestampProto(runDetails.QueueTime)
	startTime, _ := ptypes.TimestampProto(runDetails.StartTime)
	endTime, _ := ptypes.TimestampProto(runDetails.EndTime)
	return proto.Marshal(&Job{
		Id:     m.ID,
		Status: status,
		RunDetails: &RunDetails{
			QueueTime: queueTime,
			StartTime: startTime,
			EndTime:   endTime,
			Error:     runDetails.Error,
		},
		Context: m.Context,
		Result: &RscriptResult{
			Output:   m.Result.Output,
			ExitCode: m.Result.ExitCode,
		},
		Rscript: &Rscript{
			RPath:       m.Rscript.RPath,
			WorkDir:     m.Rscript.WorkDir,
			RscriptPath: m.Rscript.RscriptPath,
			Renv:        m.Rscript.Renv,
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

	queueTime, _ := ptypes.Timestamp(runDetails.QueueTime)
	startTime, _ := ptypes.Timestamp(runDetails.StartTime)
	endTime, _ := ptypes.Timestamp(runDetails.EndTime)
	m.RunDetails = server.RunDetails{
		QueueTime: queueTime,
		StartTime: startTime,
		EndTime:   endTime,
		Error:     runDetails.Error,
	}
	m.Context = pb.Context
	m.Result = server.Result{
		Output:   pb.Result.Output,
		ExitCode: pb.Result.ExitCode,
	}
	m.Rscript = server.Rscript{
		RPath:       pb.Rscript.RPath,
		WorkDir:     pb.Rscript.WorkDir,
		RscriptPath: pb.Rscript.RscriptPath,
		Renv:        pb.Rscript.Renv,
	}

	return nil
}
