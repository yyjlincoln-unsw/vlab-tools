package hooks

import (
	"fmt"
	"os"
	"strings"
)

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
