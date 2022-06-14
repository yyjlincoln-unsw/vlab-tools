package main

import (
	"autotest/commands"
	"fmt"
	"math/rand"
	"os"
	"os/user"

	"github.com/fatih/color"
)

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

func GetCurrentUser() string {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return ""
	}
	return u.Username
}

var randNumber = fmt.Sprintf("%v", rand.Intn(1000000))

var AUTOTEST_MAP = []*CourseInformation{
	Class_CS1511,
	Class_CS2521,
}

var ErrorOutput = func(format string, args ...interface{}) {
	color.Set(color.FgRed)
	fmt.Printf(format, args...)
	color.Unset()
}

var TitleOutput = func(format string, args ...interface{}) {
	fmt.Printf("\n\n")
	color.Set(color.BgWhite).Add(color.FgBlack)
	fmt.Printf(format, args...)
	color.Unset()
	fmt.Printf("\n\n")
}

var Class_CS1511 = &CourseInformation{
	Identifier: "cs1511",
	CourseName: "COMP1511 - Term 2, 2022",
	Tasks: []*Task{
		{
			Identifier: "asm0-test",
			Name:       "Assignment 0 - Test your code",
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
					// Cleanup
					{
						Command: "rm",
						Args: []string{
							"/tmp/cs_chardle_build_" + randNumber,
						},
					},
				},
			},
		},
		{
			Identifier: "asm0-reference",
			Name:       "Assignment 0 - Generate Reference Outputs",
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
		{
			Identifier: "asm0-clear",
			Name:       "Assignment 0 - Clear my previous attempts",
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
		{
			Identifier: "asm0-style",
			Name:       "Assignment 0 - Style Check",
			CommandSet: &commands.CommandSet{
				Commands: []commands.Command{
					{
						Command: "1511",
						Args: []string{
							"style",
							"cs_chardle.c",
						},
					},
				},
			},
		},
	},
}

var Class_CS2521 = &CourseInformation{
	Identifier: "cs2521",
	CourseName: "COMP2521 - Term 2, 2022",
	Tasks: []*Task{
		{
			Identifier: "asm1-test",
			Name:       "Assignment 1 - Test your code",
			CommandSet: &commands.CommandSet{
				Commands: []commands.Command{
					{
						Command: "echo",
						Args: []string{
							"This feature is not available at the moment.",
						},
					},
				},
			},
		},
	},
}
