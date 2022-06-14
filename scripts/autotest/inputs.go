package main

import "fmt"

func AskForCourseInfo() *CourseInformation {
	TitleOutput("Course Selection")
	if courseInfoFromProfile := CheckProfileForCourseInfo(); courseInfoFromProfile != nil {
		return courseInfoFromProfile
	}
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
	return courseInfo
}

func AskForTaskInfo(courseInfo *CourseInformation) *Task {
	TitleOutput("Task Selection")
	fmt.Printf("Please select one of the tasks: \n")
	for key := range courseInfo.Tasks {
		fmt.Printf("\t%v:\t%v\n", courseInfo.Tasks[key].Identifier, courseInfo.Tasks[key].Name)
	}
	if len(courseInfo.Tasks) == 0 {
		ErrorOutput("\tNo task is available at the moment.\n")
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
	return taskInfo
}
