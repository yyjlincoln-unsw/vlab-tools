package commands

import (
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	Command string
	Args    []string
}

type CommandSet struct {
	Commands []Command
}

// Runs the command in interactive mode, returns the return code
// and error (not the program's error).
func runCommand(executable string, args []string) (int, error) {
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

func ExecuteCommandSet(CommandSet *CommandSet) (int, error) {
	for _, command := range CommandSet.Commands {
		code, err := runCommand(command.Command, command.Args)
		if err != nil {
			return code, err
		}
		if code != 0 {
			return code, err
		}
	}
	return 0, nil
}
