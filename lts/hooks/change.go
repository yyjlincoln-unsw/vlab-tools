package hooks

import (
	"fmt"
	"os"
	"strings"
)

var ELIGIBLE_FILE_EXTENSIONS = []string{
	"c",
	"h",
	"go",
	"py",
	"sh",
	"js",
	"ts",
	"tsx",
	"jsx",
	"css",
	"scss",
	"sass",
	"html",
}

func GetFileName(path string) string {
	split := strings.Split(path, fmt.Sprintf("%v", os.PathSeparator))
	return split[len(split)-1]
}

func GetExtension(path string) string {
	fileName := GetFileName(path)
	split := strings.Split(fileName, ".")
	if len(split) == 1 {
		return ""
	}
	return split[len(split)-1]
}

func GetFileEligibility(path string) bool {
	ext := GetExtension(path)
	for _, e := range ELIGIBLE_FILE_EXTENSIONS {
		if e == ext {
			return true
		}
	}
	return false
}
