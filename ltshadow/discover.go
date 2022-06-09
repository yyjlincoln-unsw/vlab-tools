package main

import (
	"fmt"
	"os"
)

var FILE_PATHS []string = []string{
	"/home/z5423219/local-public/internals",
	"/home/z5423219/local-public",
}

func SearchForFile(fileName string) (string, bool) {
	for _, path := range FILE_PATHS {
		filePath := fmt.Sprintf("%v%c%v", path, os.PathSeparator, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, true
		}
	}
	return "", false
}

func SearchForExecutable(execName string) (string, bool) {
	return SearchForFile(execName)
}

func SearchForShell(execName string) (string, bool) {
	return SearchForFile(execName + ".sh")
}

func SearchForPython(execName string) (string, bool) {
	return SearchForFile(execName + ".py")
}
