package finder

import (
	"fmt"
	"os"
	"strings"
)

func GetParentDirectory(path string) (string, error) {
	// Remove trailing "/"
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	sep := strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
	if len(sep) == 1 {
		return "", fmt.Errorf("no parent directory is available")
	}
	newPath := ""
	sepLast := len(sep) - 1
	for i, v := range sep {
		if i != sepLast {
			newPath += v
			if i != sepLast-1 {
				newPath += string(os.PathSeparator)
			}
		}
	}
	if newPath == "" {
		newPath = "/"
	}
	return newPath, nil
}

func PathForFileAtDirectory(path1 string, path2 string) string {
	return fmt.Sprintf("%v%c%v", path1, os.PathSeparator, path2)
}

func FindFile(fileName string) (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("in FindLTS: %v", err)
	}
	for currentDirectory != "" {
		path := PathForFileAtDirectory(currentDirectory, fileName)
		if _, err = os.Stat(path); err == nil {
			return path, nil
		}
		currentDirectory, err = GetParentDirectory(currentDirectory)
		if err != nil {
			break
		}
	}
	return "", fmt.Errorf("in FindLTS: unable to find %v", fileName)
}

func FindLTS() (string, error) {
	return FindFile("lts.json")
}
