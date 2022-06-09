package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const VERSION = "v1.0"

var EXECUTABLE_PATHS []string = []string{
	"/home/z5423219/local-public/internals",
	"/home/z5423219/local-public/bin",
	"/home/z5423219/local-public",
}

func main() {
	execNameSplit := strings.Split(os.Args[0], fmt.Sprintf("%c", os.PathSeparator))
	execName := execNameSplit[len(execNameSplit)-1]
	args := os.Args[1:]

	executed := false

	defer func(execName string) {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Could not execute %v due to an internal error.\n", execName)
		}
	}(execName)

	// Search from local executable
	for _, path := range EXECUTABLE_PATHS {
		executable := fmt.Sprintf("%v%c%v", path, os.PathSeparator, execName)
		if _, err := os.Stat(executable); err == nil {
			// Execute the command.
			cmd := exec.Command(executable, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err = cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if err = cmd.Wait(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
			}

			executed = true
			break
		}
	}

	if !executed {
		fmt.Fprintf(os.Stderr, "Shadow: Could not find the executable %v.\n", execName)
	}
}
