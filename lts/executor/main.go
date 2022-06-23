package executor

import (
	"fmt"
	"os"
	"os/exec"
)

// Returns a Kill function
func WaitForCompletionOrKill(cmd *exec.Cmd, onCompletion func(int, error)) (chan int, func()) {

	kill := make(chan int)
	killed := false

	exitCode := make(chan int)
	errWhenDone := make(chan error)

	go func() {
		if err := cmd.Wait(); err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				exitCode <- err.ExitCode()
				close(exitCode)
				errWhenDone <- nil
				close(errWhenDone)
			} else {
				exitCode <- 0
				close(exitCode)
				errWhenDone <- err
				close(errWhenDone)
			}
		} else {
			errWhenDone <- nil
			exitCode <- 0
		}
	}()

	done := make(chan int)
	go func() {
		select {
		case <-kill:
			cmd.Process.Kill()
			cmd.Process.Release()
			done <- 1
			fmt.Printf("Kill")
			close(done)
			return
		case code := <-exitCode:
			done <- 1
			fmt.Printf("Done")
			close(done)
			onCompletion(code, <-errWhenDone)
			return
		}
	}()

	return done, func() {
		if killed {
			return
		}
		// Kill the program
		killed = true
		kill <- 1
	}
}

// Returns the Kill function
func ExecuteShell(command string, onCompletion func(int, error)) (chan int, func(), error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		done := make(chan int)
		done <- 1
		close(done)
		return done, func() {}, err
	}
	done, kill := WaitForCompletionOrKill(cmd, onCompletion)
	return done, kill, nil
}
