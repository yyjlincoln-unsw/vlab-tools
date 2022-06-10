package main

import (
	"fmt"
	"os"
	"strings"
)

const VERSION = "v1.1"

func main() {
	execNameSplit := strings.Split(os.Args[0], fmt.Sprintf("%c", os.PathSeparator))
	execName := execNameSplit[len(execNameSplit)-1]
	args := os.Args[1:]

	executed, code, err := SearchAndExecute(execName, args, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return
	}
	if executed {
		os.Exit(code)
		return
	}
	executed, code, err = ExecuteCommandFromCloud(execName, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}
	if executed {
		os.Exit(code)
		return
	}
	fmt.Fprintf(os.Stderr, "Error: Command not found.\n")
}
