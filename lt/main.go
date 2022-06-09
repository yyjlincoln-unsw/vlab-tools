package main

import (
	"fmt"
	"os"
)

const VERSION = "v1.0"

func main() {
	selfExec := os.Args[0]
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Lincoln's Tools - Version %v\n", VERSION)
		fmt.Fprintf(os.Stderr, "Error: No arguments provided.\n")
		os.Exit(1)
	}
	execName := args[0]
	args = args[1:]
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
	if filePath, found := SearchForFile(execName); found {
		if code, err := RunCommand(filePath, args); err == nil {
			executed = true
			os.Exit(code)
		}
		return
	}

	if filePath, found := SearchForShell(execName); found {
		args = append([]string{filePath}, args...)
		if code, err := RunCommand("sh", args); err == nil {
			executed = true
			os.Exit(code)
		}
		return
	}

	if filePath, found := SearchForPython(execName); found {
		args = append([]string{filePath}, args...)
		if code, err := RunCommand("python3", args); err == nil {
			executed = true
			os.Exit(code)
		}
	}

	if !executed {
		if err := ExecuteCommandFromCloud(execName, args); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	}
}
