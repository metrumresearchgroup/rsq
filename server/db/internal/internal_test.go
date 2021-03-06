package internal

import (
	"testing"
	"time"

	"reflect"

	"github.com/metrumresearchgroup/rsq/server"
)

func TestMarshalJob(t *testing.T) {
	testJob := server.Job{
		ID:     uint64(1234),
		Status: "COMPLETED",
		User:   "Ola",
		RunDetails: server.RunDetails{
			QueueTime: time.Now().AddDate(0, 0, -1).UTC(),
			StartTime: time.Now().AddDate(0, 0, 0).UTC(),
			EndTime:   time.Now().AddDate(0, 0, 1).UTC(),
			Error:     "no error",
		},
		Context: "interesting analysis 1",
		Rscript: server.Rscript{
			RPath:       "R",
			WorkDir:     "/some/dir",
			RscriptPath: "/some/dir/script.R",
			Renv:        map[string]string{"R_SITE_LIB": "/some/R/path"},
		},
		Result: server.Result{
			Output:   "some awesome output",
			ExitCode: 1,
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
