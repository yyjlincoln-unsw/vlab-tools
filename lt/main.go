package main

import (
	"fmt"
	"os"
	"os/exec"
)

const VERSION = "v1.0"

var EXECUTABLE_PATHS []string = []string{
	"/home/z5423219/local-public/internals",
	"/home/z5423219/local-public/bin",
	"/home/z5423219/local-public",
}

func main() {
	selfExec := os.Args[0]
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Lincoln's Tools - Version %v\n", VERSION)
		fmt.Fprintf(os.Stderr, "Error: No arguments provided.\n")
		os.Exit(1)
	}
	execName := args[0]
	if execName == "version" {
		fmt.Printf("%v\n", VERSION)
		os.Exit(0)
	}
	if execName == "help" {
		fmt.Printf("Lincoln's Tools - Version %v\n", VERSION)
		fmt.Printf("Usage: %v <command> \n", selfExec)
		fmt.Printf("\t<command> - The command to run.\n")
		os.Exit(0)
	}

	executed := false

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
				// fmt.Fprintf(os.Stderr, "%v\n", err)
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
			}

			executed = true
			break
		}
	}

	if !executed {
		if err := ExecuteCommandFromCloud(execName, args); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
}
