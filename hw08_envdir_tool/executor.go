package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	commandName := cmd[0]
	commandArgs := cmd[1:]

	command := exec.Command(commandName, commandArgs...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for name, envValue := range env {
		if envValue.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				return 1
			}
		} else {
			if err := os.Setenv(name, envValue.Value); err != nil {
				return 1
			}
		}
	}

	command.Env = os.Environ()

	err := command.Run()
	if err != nil {
		var cmdErr *exec.ExitError
		if errors.As(err, &cmdErr) {
			return cmdErr.ExitCode()
		}
		return 1
	}

	return 0
}
