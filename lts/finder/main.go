package finder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// LTS Command List ADT

type LTSCommandList struct {
	commands map[string]string
}

func newLTSCommandList() LTSCommandList {
	return LTSCommandList{
		commands: map[string]string{},
	}
}

func (list LTSCommandList) GetCommand(name string) (string, error) {
	cmd, ok := list.commands[name]
	if ok {
		return cmd, nil
	}
	return "", fmt.Errorf("command not found: %v", name)
}

func (list LTSCommandList) addCommand(name string, command string) {
	list.commands[name] = command
}

func ReadCommandList() (*LTSCommandList, error) {
	// Try and find LTS
	file, err := FindLTS()
	if err != nil {
		return nil, fmt.Errorf("in RadCommandList: %v", err)
	}
	// Read the file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("in RadCommandList: %v", err)
	}
	// Unmarshal the JSON
	commands := map[string]string{}
	if err := json.Unmarshal(content, &commands); err != nil {
		return nil, fmt.Errorf("in RadCommandList: %v", err)
	}
	// Load them into LTSCommandList
	list := newLTSCommandList()
	for k, v := range commands {
		list.addCommand(k, v)
	}
	return &list, nil
}
