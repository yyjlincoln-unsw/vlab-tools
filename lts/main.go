package main

import (
	"fmt"
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
	if len(Args) != 0 {
		ShowHelp(LTSExecutableName)
		os.Exit(1)
	}
	path, err := finder.FindLTS()
	fmt.Printf("%v %v\n", path, err)
}

func ShowHelp(ExecutableName string) {
	fmt.Printf("Usage:\n")
	fmt.Printf("\t%v [name]\n\n", ExecutableName)
	fmt.Printf("This program will try and read lts.json from the current directory or the parent directories, then execute the script using shell.\n")
}
