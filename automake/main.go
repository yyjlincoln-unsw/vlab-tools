package main

import (
	"automake/maker"
	"automake/watcher"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		help := "Automatically make when file changes, and execute a file name.\nUsage: %v <filename> [args...]\n"
		fmt.Fprintf(os.Stderr, help, os.Args[0])
		os.Exit(1)
	}

	fileName := GetFileName(os.Args[1])
	args := os.Args[1:]

	watcher, err := watcher.New("./")
	done := make(chan int)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not spawn a watcher: %v\n", watcher)
		done <- 1
		os.Exit(1)
	}

	var currentKillSig chan int = nil
	if _, err := watcher.OnChange(func(file string) {
		if currentKillSig != nil {
			// Kill previous
			currentKillSig <- 1
		}
		if GetFileName(file) != fileName {
			maker.ReMake(fileName, args)
		} else {
			fmt.Printf("Not remaking.\n")
		}
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create onChange handler: %v", err)
	}
	<-done
	defer watcher.Destroy()
}

func GetFileName(path string) string {
	split := strings.Split(path, fmt.Sprintf("%v", os.PathSeparator))
	return split[len(split)-1]
}
