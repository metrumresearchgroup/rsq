package internal

import (
	"testing"
	"time"

	"reflect"

	"github.com/gofrs/uuid/v3"
	"github.com/metrumresearchgroup/rsq/server"
)

func TestMarshalJob(t *testing.T) {
	u1, _ := uuid.NewV4()
	testJob := server.Job{
		ID:     u1.String(),
		Status: "COMPLETED",
		RunDetails: server.RunDetails{
			QueueTime: time.Now().AddDate(0, 0, -1).UTC(),
			StartTime: time.Now().AddDate(0, 0, 0).UTC(),
			EndTime:   time.Now().AddDate(0, 0, 1).UTC(),
			Error:     "no error",
		},
	}

	var result server.Job
	if buf, err := MarshalJob(&testJob); err != nil {
		t.Fatal(err)
	} else if err := UnmarshalJob(buf, &result); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(testJob, result) {
		// t.Log("output")
		// t.Logf("%s", testJob)
		// t.Logf("%s", result)
		t.Fatalf("unexpected copy: %#v", result)
	}
}
