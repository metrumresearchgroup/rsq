package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const defaultFailedCode = 1
const defaultSuccessCode = 0

// RunR launches an interactive R console
func RunR(
	fs afero.Fs,
	rs RSettings,
	rdir string, // this should be put into RSettings
	lg *logrus.Logger,
) error {

	cmdArgs := []string{
		"--no-save",
		"--no-restore-data",
	}

	envVars := os.Environ()
	ok, rLibsSite := rs.LibPathsEnv()
	if ok {
		envVars = append(envVars, rLibsSite, "R_LIBS=''")
	}

	lg.WithFields(
		logrus.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       rLibsSite,
		}).Debug("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		rs.R(),
		cmdArgs...,
	)

	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	cmd.Env = envVars
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RunRscript runs an Rscript and returns the result, along with the any errors and exit code
// exit code 0 - success, 1 - error, 2 - script not found`
func RunRscript(
	fs afero.Fs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger,
) (string, error, int) {

	cmdArgs := []string{
		"--no-save",
		"--no-restore-data",
		es.Rfile,
	}

	envVars := configureEnv(os.Environ(), rs)

	lg.WithFields(
		logrus.Fields{
			"cmdArgs":      cmdArgs,
			"RSettings":    rs,
			"ExecSettings": es,
		}).Debug("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		fmt.Sprintf("%sscript", rs.R()),
		cmdArgs...,
	)

	rdir := es.WorkDir
	if rdir == "" {
		rdir, _ = os.Getwd()
		lg.WithFields(
			logrus.Fields{"rdir": rdir},
		).Debug("launch dir")
	}
	cmd.Dir = rdir
	cmd.Env = envVars
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	exitCode := cmd.ProcessState.ExitCode()
	lg.WithFields(
		logrus.Fields{
			"stdout":   stdout,
			"stderr":   stderr,
			"exitCode": exitCode,
		}).Info("cmd output")
	return stdout, err, exitCode
}
