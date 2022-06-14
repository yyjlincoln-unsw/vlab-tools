package main

import (
	"autotest/commands"
	"fmt"
	"os"
)

const VERSION = "v1.2.2"

func main() {
	args := os.Args[1:]
	DealWithArgs(args)

	fmt.Printf("LT Autotest %v\nA wrapper of cloud-autotest that provides a nice UI.\n", VERSION)
	courseInfo := AskForCourseInfo()
	if courseInfo == nil {
		ErrorOutput("Invalid course ID.\n")
		os.Exit(1)
	}
	taskInfo := AskForTaskInfo(courseInfo)
	if taskInfo == nil {
		ErrorOutput("Invalid task ID.\n")
		os.Exit(1)
	}
	RunTask(taskInfo)
}

func RunTask(taskInfo *Task) {
	if taskInfo == nil {
		ErrorOutput("Invalid task info.\n")
		os.Exit(1)
	}

	fmt.Printf("\nRunning %v...\n", taskInfo.Name)

	code, err := commands.ExecuteCommandSet(taskInfo.CommandSet)
	if err != nil {
		ErrorOutput("Error: %v\n", err)
	}
	os.Exit(code)
}
