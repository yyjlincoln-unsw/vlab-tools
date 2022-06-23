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
				logging.Successf("Command exited.\n")
			}
		})
	}
	done, kill, err := execute()
	currentKill = &kill
	if err != nil {
		logging.Errorf("Error: %v\n", err)
	}

	// Get hooks
	hooksForCommand := list.GetHooks(CommandName)
	dones := []chan int{}
	for _, v := range hooksForCommand {
		dones = append(dones, hooks.RegisterHook(v, func() {
			if currentKill != nil {
				fn := *(currentKill)
				fn()
			}

			_, kill, err := execute()
			*currentKill = kill

			if err != nil {
				logging.Errorf("Error: %v\n", err)
			}
		}))
	}

	// Cleanup child when killed
	sigKill := make(chan os.Signal, 3)
	signal.Notify(sigKill, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)
	go func() {
		<-sigKill
		if currentKill != nil {
			fn := *currentKill
			go fn()
		}
		logging.Errorf("\nExiting.\n")
		os.Exit(0)
	}()

	hooks.WaitForAllHooks(append(dones, done))
}

func ShowHelp(ExecutableName string) {
	fmt.Printf("Usage:\n")
	fmt.Printf("\t%v [name]\n\n", ExecutableName)
	fmt.Printf("This program will try and read lts.json from the current directory or the parent directories, then execute the script using shell.\n")
}
