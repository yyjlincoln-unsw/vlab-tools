package main

import "fmt"

func Which(execName string) {
	// Without local cache
	if filePath, found := SearchForExecutable(execName, false); found {
		fmt.Printf("%v\n", filePath)
		return
	}
	if filePath, found := SearchForShell(execName, false); found {
		fmt.Printf("%v\n", filePath)
		return
	}
	if filePath, found := SearchForPython(execName, false); found {
		fmt.Printf("%v\n", filePath)
		return
	}

	// With local cache
	if err := DownloadCommandFromCloud(execName); err == nil {
		fmt.Printf("Downloaded cloud command: \n")
		if filePath, found := SearchForExecutable(execName, true); found {
			fmt.Printf("%v\n", filePath)
			return
		}
		if filePath, found := SearchForShell(execName, true); found {
			fmt.Printf("%v\n", filePath)
			return
		}
		if filePath, found := SearchForPython(execName, true); found {
			fmt.Printf("%v\n", filePath)
			return
		}
		return
	}
	fmt.Printf("Command not found\n")
}
