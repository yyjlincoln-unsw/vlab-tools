package main

import (
	"fmt"
	"os"
	"strings"
)

const VERSION = "v1.0"

func main() {
	execNameSplit := strings.Split(os.Args[0], fmt.Sprintf("%c", os.PathSeparator))
	execName := execNameSplit[len(execNameSplit)-1]
	args := os.Args[1:]

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
