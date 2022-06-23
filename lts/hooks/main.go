package hooks

import "fmt"

// Returns a blocking channel that, when the hook is dead, sends.
func RegisterHook(name string, callback func()) chan int {
	done := make(chan int)

	go func() {
		switch name {
		case "change":
			// File system change.
			for {
				// Monitor for FS Change, and run the command
			}
		}
		fmt.Printf("Close %v\n", name)
		done <- 1
		close(done)
	}()

	return done
}

func WaitForAllHooks(hooks []chan int) {
	for _, v := range hooks {
		<-v
	}
}
