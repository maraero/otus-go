package main

import (
	"os"
	"os/exec"
)

const (
	exitSuccess = 0
	exitError   = 1
)

func handleEnvs(env Environment) error {
	for name, val := range env {
		var err error

		if val.NeedRemove == true {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, val.Value)
		}

		return err
	}

	return nil
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return exitError
	}

	if err := handleEnvs(env); err != nil {
		return exitError
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	err := c.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	if err != nil {
		return exitError
	}

	return exitSuccess
}
