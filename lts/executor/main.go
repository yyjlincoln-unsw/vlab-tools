package executor

import (
	"os"
	"os/exec"
)

func ExecuteShell(command string) (int, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return err.ExitCode(), nil
		}
		return 0, err
	}
	return 0, nil
}
