package main

import (
	"autotest/commands"
	"fmt"
	"math/rand"
)

var randNumber = fmt.Sprintf("%v", rand.Intn(1000000))

var AUTOTEST_MAP = map[string]*CourseInformation{
	"cs1511": Class_CS1511,
	"cs2521": Class_CS2521,
}

var Class_CS1511 = &CourseInformation{
	ID:         "cs1511",
	CourseName: "COMP1511 - Term 2, 2022",
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
}

var Class_CS2521 = &CourseInformation{
	ID:         "cs2521",
	CourseName: "COMP2521 - Term 2, 2022",
	Tasks:      map[string]*Task{
		
	},
}
