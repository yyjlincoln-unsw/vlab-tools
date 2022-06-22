package main

import (
	"automake/watcher"
	"fmt"
	"os"
)

func main() {
	watcher, err := watcher.New("./")
	done := make(chan int)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not spawn a watcher: %v\n", watcher)
		done <- 1
		os.Exit(1)
	}
	if _, err := watcher.OnChange(func(file string) {
		fmt.Printf("File changed: %s\n", file)
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create onChange handler: %v", err)
	}
	<-done
	defer watcher.Destroy()
}
