package main

import (
	"autotest/commands"
	"fmt"
	"os"
	"os/user"
)

const VERSION = "v1.1"

func GetCurrentUser() string {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return ""
	}
	return u.Username
}

type Task struct {
	Name       string
	Identifier string
	CommandSet *commands.CommandSet
}

type CourseInformation struct {
	Identifier string
	CourseName string
	Tasks      []*Task
}

func main() {
	fmt.Printf("LT Autotest %v\nA wrapper of cloud-autotest that provides a nice UI.\n\n", VERSION)
	fmt.Printf("Please select one of the courses: \n")
	for key := range AUTOTEST_MAP {
		fmt.Printf("\t%v:\t%v\n", AUTOTEST_MAP[key].Identifier, AUTOTEST_MAP[key].CourseName)
	}
	var courseId string
	fmt.Printf("\nCourse: ")
	fmt.Scanln(&courseId)
	var courseInfo *CourseInformation = nil
	for i := range AUTOTEST_MAP {
		if AUTOTEST_MAP[i].Identifier == courseId {
			courseInfo = AUTOTEST_MAP[i]
		}
	}
	if courseInfo == nil {
		fmt.Printf("Invalid course ID.\n")
		os.Exit(1)
	}
	fmt.Printf("\n\nPlease select one of the tasks: \n")
	for key := range courseInfo.Tasks {
		fmt.Printf("\t%v:\t%v\n", courseInfo.Tasks[key].Identifier, courseInfo.Tasks[key].Name)
	}
	if len(courseInfo.Tasks) == 0 {
		fmt.Printf("\tNo task is available at the moment.\n")
	}
	fmt.Printf("\nTask: ")
	var taskId string
	fmt.Scanln(&taskId)
	var taskInfo *Task = nil
	for key := range courseInfo.Tasks {
		if courseInfo.Tasks[key].Identifier == taskId {
			taskInfo = courseInfo.Tasks[key]
			break
		}
	}
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
