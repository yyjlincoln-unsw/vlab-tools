package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Runs the command in interactive mode, returns the return code
// and error (not the program's error).
func RunCommand(executable string, args []string) (int, error) {
	// Execute the command.
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("error: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		}
	}
	return 0, nil
}
