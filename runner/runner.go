package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

	var cmd *exec.Cmd
	// run_file won't work, because go run's subprocess won't be killed
	if runFile() != "" {
		cmd = exec.Command("go", "run", runFile())
		runnerLog("Execute %s", cmd.Args)
	} else {
		cmd = exec.Command(buildPath())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
