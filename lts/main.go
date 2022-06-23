package main

import (
	"fmt"
	"lts/executor"
	"lts/finder"
	"lts/logging"
	"os"
)

func main() {
	if len(os.Args) == 0 {
		logging.Errorf("Error: No argument was provided to the program.\n")
		os.Exit(1)
	}
	LTSExecutableName, Args := os.Args[0], os.Args[1:]
	if len(Args) == 0 {
		ShowHelp(LTSExecutableName)
		os.Exit(1)
	}
	CommandName := Args[0]
	// Read the command list
	list, err := finder.ReadCommandList()
	if err != nil {
		logging.Errorf("Error: Unable to read command: %v\n", err)
		os.Exit(1)
	}
	command, err := list.GetCommand(CommandName)
	if err != nil {
		logging.Errorf("Error: Unable to execute %v: %v\n", CommandName, err)
		os.Exit(1)
	}
	code, err := executor.ExecuteShell(command)
	if err != nil {
		logging.Errorf("%v", err)
		os.Exit(1)
	}
	if code != 0 {
		logging.Errorf("Exit status %v\n", code)
	} else {
		logging.Successf("Command exited.\n")
	}
	os.Exit(code)
}

func ShowHelp(ExecutableName string) {
	fmt.Printf("Usage:\n")
	fmt.Printf("\t%v [name]\n\n", ExecutableName)
	fmt.Printf("This program will try and read lts.json from the current directory or the parent directories, then execute the script using shell.\n")
}
