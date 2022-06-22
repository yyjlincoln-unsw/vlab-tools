package main

import (
	"automake/maker"
	"automake/watcher"
	"fmt"
	"os"
	"strings"
)

var WATCH_EXTENSIONS []string = []string{
	"c",
	"h",
}

func main() {
	if len(os.Args) < 2 {
		help := "Automatically make when file changes, and execute a file IN THE CURRENT DIRECTORY (must be a name or starts with \"./\").\nUsage: %v <filename> [args...]\n"
		fmt.Fprintf(os.Stderr, help, os.Args[0])
		os.Exit(1)
	}

	targetFileName := GetFileName(os.Args[1])
	args := os.Args[1:]

	watcher, err := watcher.New("./")
	done := make(chan int)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not spawn a watcher: %v\n", watcher)
		done <- 1
		os.Exit(1)
	}
	_, currentKillSig := maker.ReMakeWithThrottle(targetFileName, args)

	if _, err := watcher.OnChange(func(file string) {
		if IsFileEligible(file, targetFileName) {
			fmt.Printf("Change: %v\n", file)
			if currentKillSig != nil {
				// Kill previous
				select {
				case currentKillSig <- 1:
				default:
				}
			}
			if ran, killSig := maker.ReMakeWithThrottle(targetFileName, args); ran {
				currentKillSig = killSig
			}
		}
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create onChange handler: %v", err)
	}
	<-done
	defer watcher.Destroy()
}

func IsFileEligible(file string, targetFileName string) bool {
	ext := GetExtension(file)
	for _, curr := range WATCH_EXTENSIONS {
		if curr == ext {
			return true
		}
	}
	return false
	// if GetFileName(file) == targetFileName {
	// 	return false
	// }
	// if ext == "" {
	// 	return false
	// }
	// if ext == "o" {
	// 	return false
	// }
	// if ext == "tmp" {
	// 	return false
	// }
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
