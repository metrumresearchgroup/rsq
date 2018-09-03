package internal

import (
	"testing"
	"time"

	"reflect"

	"github.com/metrumresearchgroup/rsq/server"
)

func TestMarshalJob(t *testing.T) {
	u1, err := uuid.NewV4()
	testJob := server.Job{
		ID:     u1,
		Status: "COMPLETED",
		RunInfo: server.RunInfo{
			QueueTime: time.Now().AddDate(0, 0, -1),
			StartTime: time.Now().AddDate(0, 0, -1),
			EndTime:   time.Now().AddDate(0, 0, -1),
			Error:     "success",
		},
	}

	var result server.Job
	if buf, err := MarshalModel(&testJob); err != nil {
		t.Fatal(err)
	} else if err := UnmarshalModel(buf, &result); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(testJob, result) {
		t.Fatalf("unexpected copy: %#v", result)
	}
}
