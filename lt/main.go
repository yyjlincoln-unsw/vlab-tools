package main

import (
	"fmt"
	"os"
)

const VERSION = "v1.1"

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
	if execName == "which" {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: please provide a command name.\n")
			os.Exit(1)
		}
		execName = args[0]
		Which(execName)
		return
	}

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
