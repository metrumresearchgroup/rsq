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
	runSettings := runner.RSettings{
		LibPaths: []string{},
		Rpath:    "R",
	}
	wd, _ := os.Getwd()
	execSettings := runner.ExecSettings{
		WorkDir: wd,
		Rfile:   "add-stop.R",
	}
	runner.RunRscript(appFS,
		runSettings,
		execSettings,
		lg)
}
