package executor

import (
	"os"
	"os/exec"
	"syscall"
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
			exitCode <- 0
			close(exitCode)
			errWhenDone <- nil
			close(errWhenDone)
		}
	}()

	done := make(chan int)
	go func() {
		select {
		case <-kill:
			// Also kill its children
			pgid, err := syscall.Getpgid(cmd.Process.Pid)
			if err == nil {
				syscall.Kill(-pgid, syscall.SIGKILL)
			}
			cmd.Process.Kill()
			cmd.Process.Release()
			cmd.Process.Wait()
			select {
			case done <- 1:
			default:
			}
			close(done)
			return
		case code := <-exitCode:
			err := <-errWhenDone
			onCompletion(code, err)
			select {
			case done <- 1:
			default:
			}
			close(done)
			return
		}
	}()

	return done, func() {
		// Kill the program
		if killed {
			return
		}
		killed = true
		// For some reasons, killSig does not
		// work when there is no hook-triggered
		// reruns. We're going to kill it here anyway.
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
		cmd.Process.Kill()
		cmd.Process.Release()
		cmd.Process.Wait()
		// Sends the KillSig
		kill <- 1
	}
}

// Returns the Kill function
func ExecuteShell(command string, onCompletion func(int, error)) (chan int, func(), error) {
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		done := make(chan int)
		done <- 1
		close(done)
		return done, func() {}, err
	}
	done, kill := WaitForCompletionOrKill(cmd, onCompletion)
	return done, kill, nil
}
