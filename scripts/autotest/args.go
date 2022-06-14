package main

import (
	"autotest/commands"
	"fmt"
	"os"
)

func DealWithArgs(args []string) {
	if len(args) == 0 {
		return
	}
	for _, arg := range args {
		switch arg {
		case "help":
			fmt.Printf("Help: Use the option -v to show the commands that's being executed.\n-n Execute without reading the profile.")
			os.Exit(0)
		case "-v":
			commands.ShowCommands = true
			fmt.Printf("args: -v was detected, showing command.\n")
		case "version":
			fmt.Printf("%s\n", VERSION)
			os.Exit(0)
		case "-n":
			NoProfile = true
			fmt.Printf("args: -n was detected, not using profile.\n")
		default:
			ErrorOutput("Error: Invalid arguments.")
			fmt.Printf("\n")
			os.Exit(1)
		}
	}
}
