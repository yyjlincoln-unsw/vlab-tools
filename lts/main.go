package main

import (
	"fmt"
	"lts/executor"
	"lts/finder"
	"lts/hooks"
	"lts/logging"
	"os"
	"os/signal"
	"syscall"

	"github.com/inancgumus/screen"
)

const VERSION = "v1.1.3"

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
	// Check if it's internal
	if HandleBuiltInCommand(CommandName) {
		return
	}
	// Read the command list
	list, err := finder.ReadCommandList()
	if err != nil {
		logging.Errorf("Error: Unable to read command: %v\n", err)
		os.Exit(1)
	}
	// Get command
	command, err := list.GetCommand(CommandName)
	if err != nil {
		logging.Errorf("Error: Unable to execute %v: %v\n", CommandName, err)
		os.Exit(1)
	}

	// Execute the command
	var currentKill *func()
	execute := func() (chan int, func(), error) {
		return executor.ExecuteShell(command, func(code int, err error) {
			if err != nil {
				logging.Errorf("%v", err)
				return
			}
			if code != 0 {
				logging.Errorf("Exit status %v\n", code)
			} else {
				logging.Successf("Command exited with code 0.\n")
			}
		})
	}

	// Get hooks
	hooksForCommand := list.GetHooks(CommandName)
	dones := []chan int{}
	if len(hooksForCommand) > 0 {
		screen.Clear()
		screen.MoveTopLeft()
		logging.Successf("Running with hooks: %v\n", hooksForCommand)
	}
	for _, v := range hooksForCommand {
		dones = append(dones, hooks.RegisterHook(v, func() {
			if currentKill != nil {
				fn := *(currentKill)
				go fn()
			}

			_, kill, err := execute()
			*currentKill = kill

			if err != nil {
				logging.Errorf("Error: %v\n", err)
			}
		}))
	}

	// Now, execute it
	done, kill, err := execute()
	currentKill = &kill
	if err != nil {
		logging.Errorf("Error: %v\n", err)
	}

	// Cleanup child when killed
	sigKill := make(chan os.Signal, 3)
	signal.Notify(sigKill, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)
	go func() {
		<-sigKill
		logging.Errorf("\n\nCleaning up...\n")
		if currentKill != nil {
			fn := *currentKill
			go fn()
		}

		logging.Errorf("Exiting.\n")
		os.Exit(0)
	}()

	hooks.WaitForAllHooks(append(dones, done))
}

func HandleBuiltInCommand(cmd string) bool {
	switch cmd {
	case "list":
		fmt.Printf("Looking for lts.json...\n")
		path, err := finder.FindLTS()
		if err != nil {
			logging.Errorf("LTS lookup failure: %v\n", err)
			os.Exit(1)
			return true
		}
		logging.Infof("The configuration file is found at:\n")
		logging.Infof("%v\n\n", path)
		logging.Infof("Available commands...\n")
		list, err := finder.ReadCommandList()
		if err != nil {
			logging.Errorf("Parse failure: %v\n", err)
			os.Exit(1)
			return true
		}
		for name, exec := range list.Scripts {
			logging.Successf("%v\t", name)
			fmt.Printf("Hooks=")
			hooks := list.GetHooks(name)
			logging.Warnf("%v\t", hooks)
			fmt.Printf("Command=")
			fmt.Printf("%v\n", exec)
		}
		return true
	default:
		return false
	}
}

func ShowHelp(ExecutableName string) {
	logging.Warnf("LTS - Scripts Service (as a part of Lincoln's Tools)\n")
	logging.Warnf("Version: %s\n\n", VERSION)
	logging.Infof("Usage:\n")
	fmt.Printf("\t%v <name>|list\n\n", ExecutableName)
	fmt.Printf("This program will try and read lts.json from the current directory or the parent directories, then execute the script using shell.\n\n")
	logging.Infof("Built in commands:\n")
	fmt.Printf("list:\n")
	fmt.Printf("\tLists all user-defined commands and the file that defined them.\n\n")

	logging.Infof("Hooks:\n")
	fmt.Printf("change:\n \tListens for file changes in the current directory, and if a change is detected, kill the current command process (and child processes) then rerun the same command.\n")
	fmt.Printf("\tSupported extensions: %v\n", hooks.ELIGIBLE_FILE_EXTENSIONS)
}
