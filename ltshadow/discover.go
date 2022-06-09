package main

import (
	"fmt"
	"os"
)

var FILE_PATHS []string = []string{
	"/home/z5423219/local-public/internals",
	"/home/z5423219/local-public/bin",
	"/home/z5423219/local-public",
}

func SearchForFile(fileName string, useLocalCache bool) (string, bool) {
	AdditionalPaths := []string{}

	cache := GetLocalServiceCacheDirectory()
	if cache != "" && useLocalCache {
		AdditionalPaths = append(AdditionalPaths, cache)
	}

	for _, path := range append(AdditionalPaths, FILE_PATHS...) {
		filePath := fmt.Sprintf("%v%c%v", path, os.PathSeparator, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, true
		}
	}
	return "", false
}

func SearchForExecutable(execName string, useLocalCache bool) (string, bool) {
	return SearchForFile(execName, useLocalCache)
}

func SearchForShell(execName string, useLocalCache bool) (string, bool) {
	return SearchForFile(execName+".sh", useLocalCache)
}

func SearchForPython(execName string, useLocalCache bool) (string, bool) {
	return SearchForFile(execName+".py", useLocalCache)
}
