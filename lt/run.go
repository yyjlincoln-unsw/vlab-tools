package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Runs the command in interactive mode, returns the return code
// and error (not the program's error).
func RunCommand(executable string, args []string) (int, error) {
	// Execute the command.
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("run command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		}
	}
	return 0, nil
}

// Returns bool:executed, int:returnCode, error:Error before running
func SearchAndExecute(execName string, args []string, useLocalCache bool) (bool, int, error) {
	// Search from local executable
	if filePath, found := SearchForExecutable(execName, useLocalCache); found {
		code, err := RunCommand(filePath, args)
		if err == nil {
			return true, code, nil
		}
		return false, 0, fmt.Errorf("executing executable: %v", err)
	}

	if filePath, found := SearchForShell(execName, useLocalCache); found {
		args = append([]string{filePath}, args...)
		code, err := RunCommand("sh", args)
		if err == nil {
			return true, code, nil
		}
		return false, 0, fmt.Errorf("executing shell: %v", err)
	}

	if filePath, found := SearchForPython(execName, useLocalCache); found {
		args = append([]string{filePath}, args...)
		code, err := RunCommand("python3", args)
		if err == nil {
			return true, code, nil
		}
		return false, 0, fmt.Errorf("executing python3: %v", err)
	}

	return false, 0, nil
}
