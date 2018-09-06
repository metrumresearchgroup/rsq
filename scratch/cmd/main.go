package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/metrumresearchgroup/rsq/jobqueue"
	"github.com/metrumresearchgroup/rsq/runner"
	"github.com/metrumresearchgroup/rsq/server"
	"github.com/metrumresearchgroup/rsq/server/db"
	"github.com/metrumresearchgroup/rsq/server/httpserver"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	//appFS := afero.NewOsFs()
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)
	jc := jobqueue.NewJobQueue(3, func(w server.Job) {
		fmt.Println(w)
	})
	jc.Push(server.Job{ID: uint64(1), Status: "QUEUED"})
	jc.Push(server.Job{ID: uint64(2), Status: "QUEUED"})
	jc.Push(server.Job{ID: uint64(3), Status: "QUEUED"})
	jc.Push(server.Job{ID: uint64(4), Status: "QUEUED"})
	// need to sleep or will exit before goroutines finish
	time.Sleep(5 * time.Second)
	return
}

// runScriptExample(appFS, lg, "add.R")
// runScriptExample(appFS, lg, "add-stop.R")
func runScriptExample(appFS afero.Fs, lg *logrus.Logger, rfile string) {
	runSettings := runner.RSettings{
		LibPaths: []string{},
		Rpath:    "R",
	}
	wd, _ := os.Getwd()
	execSettings := runner.ExecSettings{
		WorkDir: wd,
		Rfile:   rfile,
	}
	// if passes will give exit code 0
	// if stopped will give exit code 1
	// if rscript not found will give exit code 2
	runner.RunRscript(appFS,
		runSettings,
		execSettings,
		lg)
}

func runJobCreationExamples(js server.JobService) {
	testJob := server.Job{
		Status: "COMPLETED",
		RunDetails: server.RunDetails{
			QueueTime: time.Now().AddDate(0, 0, -1).UTC(),
			StartTime: time.Now().AddDate(0, 0, 0).UTC(),
			EndTime:   time.Now().AddDate(0, 0, 1).UTC(),
			Error:     "no error",
		},
		Context: "interesting job1",
	}
	testJob2 := server.Job{
		Status: "QUEUED",
		RunDetails: server.RunDetails{
			QueueTime: time.Now().AddDate(0, 0, -1).UTC(),
			StartTime: time.Now().AddDate(0, 0, 0).UTC(),
			EndTime:   time.Now().AddDate(0, 0, 1).UTC(),
			Error:     "no error",
		},
		Context: "another curious job",
	}
	err := js.CreateJob(&testJob)
	err = js.CreateJob(&testJob2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("job set")
	// fmt.Println(testJob)
	fmt.Println("----------all jobs -----------")
	jobs, err := js.GetJobs()
	fmt.Println(fmt.Sprintf("num jobs: %v", len(jobs)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("-------- job by id ---------")
	job, err := js.GetJobByID(uint64(1))
	fmt.Println(job, err)

	fmt.Println("-------- job by status ---------")
	queuedJobs, err := js.GetJobsByStatus("QUEUED")
	fmt.Println(queuedJobs, err)
	queuedJobs, err = js.GetJobsByStatus("nonsense")
	fmt.Println(queuedJobs, err)
}

func jobServer(appFS afero.Fs) {
	client := db.NewClient()
	wd, _ := os.Getwd()
	badgerPath := filepath.Join(wd, "badger")
	if _, err := os.Stat(badgerPath); os.IsNotExist(err) {
		err := appFS.Mkdir(badgerPath, 0755)
		if err != nil {
			log.Fatalf("could not create folder for db %v", err)
			os.Exit(1)
		}
	}
	client.Path = badgerPath
	err := client.Open()
	defer client.Close()
	if err != nil {
		log.Fatalf("could not open db %v", err)
		os.Exit(1)
	}

	js := client.JobService()
	// fmt.Println("about to set job")
	// fmt.Println(testJob)
	httpserver.NewHTTPServer(js, "0.0.1-alpha", "8999")
	return
}
