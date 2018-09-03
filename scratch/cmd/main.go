package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/metrumresearchgroup/rsq/runner"
	"github.com/sirupsen/logrus"

	"github.com/metrumresearchgroup/rsq/server/db"
	"github.com/spf13/afero"
)

func main() {
	appFS := afero.NewOsFs()
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)
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
