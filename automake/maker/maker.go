package maker

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
)

var makeThrottle bool = false

func Execute(file string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(file, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("run command: %v", err)
	}

	return cmd, nil
}

// Returns done
func ExecWithKillSig(killSig chan int, next func(), command string, args ...string) chan int {
	cmd, err := Execute(command, args...)

	done := make(chan int)
	exitCode := make(chan int)
	if err != nil {
		fmt.Printf("Error: could not execute: %v", err)
		return killSig
	}

	// Wait for the command to finish
	// And put it in the exitCode.
	go func() {
		// Wait async
		err := cmd.Wait()
		if err == nil {
			exitCode <- 0
		} else {
			code, ok := err.(*exec.ExitError)
			if ok {
				exitCode <- code.ExitCode()
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
		close(exitCode)
	}()

	// Wait for either the killSig
	// or command finish (exitCode)
	go func() {
		abort := false
		select {
		case <-killSig:
			// Killed
			cmd.Process.Kill()
			cmd.Process.Release()
			// fmt.Printf("Aborted.\n")
			abort = true
		case code := <-exitCode:
			if abort {
				return
			}
			if code != 0 {
				fmt.Printf("%v\n", color.RedString(fmt.Sprintf("Error: process exited with status %d.", code)))
			} else {
				isKilled := false
				select {
				case <-killSig:
					isKilled = true
				default:
				}
				if !isKilled {
					next()
				}
			}
		}
		// Clear kill sig
		select {
		case <-killSig:
		default:
		}
	}()

	return done
}

// Returns kill signal
// fileName and args are post-made commands
func ReMake(fileName string, args []string) chan int {
	screen.Clear()
	screen.MoveTopLeft()
	fmt.Printf("%v\n", color.YellowString("Building..."))

	killSig := make(chan int)
	currentPwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", color.RedString(fmt.Sprintf("An error occurred: %v", err)))
		panic("Could not get pwd.")
	}
	filePath := fmt.Sprintf("%v%c%v", currentPwd, os.PathSeparator, fileName)

	RunProgramCommand := func() {
		ExecWithKillSig(killSig, func() {
		}, filePath, args...)
	}

	MakeCommand := func() {
		ExecWithKillSig(killSig, RunProgramCommand, "make")
	}
	MakeCommand()
	return killSig
}

func ReMakeWithThrottle(fileName string, args []string) (bool, chan int) {
	if makeThrottle {
		// fmt.Printf("Block\n")
		return false, nil
	}

	makeThrottle = true

	screen.Clear()
	screen.MoveTopLeft()
	fmt.Printf("%v\n", color.YellowString("Building..."))

	killSig := make(chan int)
	currentPwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", color.RedString(fmt.Sprintf("An error occurred: %v", err)))
		panic("Could not get pwd.")
	}
	filePath := fmt.Sprintf("%v%c%v", currentPwd, os.PathSeparator, fileName)

	go func() {
		time.Sleep(1 * time.Second)
		RunProgramCommand := func() {
			screen.Clear()
			screen.MoveTopLeft()
			fmt.Printf("%v\n\n", color.GreenString("Successfully built, running the program..."))
			ExecWithKillSig(killSig, func() {
				// That's done
			}, filePath, args...)
		}

		MakeCommand := func() {
			ExecWithKillSig(killSig, RunProgramCommand, "make")
		}
		MakeCommand()
		makeThrottle = false
	}()
	return true, killSig
}
