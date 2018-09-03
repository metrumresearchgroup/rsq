package main

import (
	"os"

	"github.com/metrumresearchgroup/rsq/runner"
	"github.com/sirupsen/logrus"

	"github.com/spf13/afero"
)

func main() {
	appFS := afero.NewOsFs()
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)
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
