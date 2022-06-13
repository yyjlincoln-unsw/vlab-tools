package main

import (
	"autotest/commands"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
)

const VERSION = "v1.0"

func GetCurrentUser() string {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return ""
	}
	return u.Username
}

// Priorities:
// - PurgeOverride
// - CmdOverride
// - RunAsUser (provides a username that's used in conjunction with TaskName)
// - TaskName

type Task struct {
	Name       string
	CommandSet *commands.CommandSet
}

type CourseInformation struct {
	ID         string
	CourseName string
	Tasks      map[string]*Task
}

var randNumber = fmt.Sprintf("%v", rand.Intn(1000000))

var AUTOTEST_MAP = map[string]*CourseInformation{
	"cs1511": {
		ID:         "cs1511",
		CourseName: "COMP1511",
		Tasks: map[string]*Task{
			"asm0-test": {
				Name: "Assignment 0 - Test your code",
				CommandSet: &commands.CommandSet{
					Commands: []commands.Command{
						{
							Command: "dcc",
							Args: []string{
								"-o",
								"/tmp/cs_chardle_build_" + randNumber,
								"cs_chardle.c",
							},
						},
						{
							Command: "lt",
							Args: []string{
								"cloud-autotest",
								"cs1511_22t2_asm0",
								GetCurrentUser(),
								"/tmp/cs_chardle_build_" + randNumber,
							},
						},
					},
				},
			},
			"asm0-reference": {
				Name: "Assignment 0 - Generate Reference Outputs",
				CommandSet: &commands.CommandSet{
					Commands: []commands.Command{
						{
							Command: "lt",
							Args: []string{
								"cloud-autotest",
								"cs1511_22t2_asm0",
								"ref",
								"lt cs_chardle",
							},
						},
					},
				},
			},
			"asm0-clear": {
				Name: "Assignment 0 - Clear my previous attempts",
				CommandSet: &commands.CommandSet{
					Commands: []commands.Command{
						{
							Command: "lt",
							Args: []string{
								"cloud-autotest-admin",
								"--taskId",
								"cs1511_22t2_asm0",
								"--workerId",
								GetCurrentUser(),
								"--purge-data",
							},
						},
					},
				},
			},
		},
	},
}

func main() {
	fmt.Printf("\nLT Autotest %v\nA wrapper of cloud-autotest that provides a nice UI.\n\n", VERSION)
	fmt.Printf("Please select one of the courses: \n")
	for key := range AUTOTEST_MAP {
		fmt.Printf("\t%v:\t%v\n", key, AUTOTEST_MAP[key].CourseName)
	}
	var courseId string
	fmt.Printf("\nCourse: ")
	fmt.Scanln(&courseId)
	courseInfo := AUTOTEST_MAP[courseId]
	if courseInfo == nil {
		fmt.Printf("Invalid course ID.\n")
		os.Exit(1)
	}
	fmt.Printf("\n\nPlease select one of the tasks: \n")
	for key := range courseInfo.Tasks {
		fmt.Printf("\t%v:\t%v\n", key, courseInfo.Tasks[key].Name)
	}
	fmt.Printf("\nTask: ")
	var taskId string
	fmt.Scanln(&taskId)
	taskInfo := courseInfo.Tasks[taskId]
	if taskInfo == nil {
		fmt.Printf("Invalid task ID.\n")
		os.Exit(1)
	}
	RunTask(taskInfo)
}

func RunTask(taskInfo *Task) {
	if taskInfo == nil {
		fmt.Printf("Invalid task info.\n")
		os.Exit(1)
	}

	fmt.Printf("\nRunning %v...\n", taskInfo.Name)

	code, err := commands.ExecuteCommandSet(taskInfo.CommandSet)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	os.Exit(code)
}

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
