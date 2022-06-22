package maker

import (
	"fmt"
	"os"
	"os/exec"
)

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
		fmt.Printf("Error: could not make: %v", err)
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
			code := err.(*exec.ExitError).ExitCode()
			exitCode <- code
		}
	}()

	// Wait for either the killSig
	// or command finish (exitCode)
	go func() {
		select {
		case <-killSig:
			// Killed
			cmd.Process.Kill()
			fmt.Printf("Aborted.\n")
			close(killSig)
		case code := <-exitCode:
			if code != 0 {
				fmt.Printf("Error: process exited with status %d.\n", code)
			} else {
				next()
			}
		}
		close(exitCode)
	}()

	return done
}

// Returns kill signal
// fileName and args are post-made commands
func ReMake(fileName string, args []string) chan int {
	killSig := make(chan int)

	RunProgramCommand := func() {
		ExecWithKillSig(killSig, func() {}, "echo", "Execute")
	}

	MakeCommand := func() {
		ExecWithKillSig(killSig, RunProgramCommand, "make")
	}

	MakeCommand()
	return killSig
}
