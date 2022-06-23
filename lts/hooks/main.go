package hooks

import (
	"fmt"
	"lts/logging"
	"lts/watcher"

	"github.com/inancgumus/screen"
)

// Returns a blocking channel that, when the hook is dead, sends.
func RegisterHook(name string, callback func()) chan int {
	done := make(chan int)

	go func() {
		switch name {
		case "change":
			// File system change.
			w, err := watcher.New("./")

			if err != nil {
				logging.Errorf("Could not register change hook: %v", err)
				break
			}

			done := make(chan int)
			w.OnChange(func(file string) {
				logging.Infof("Change: %v\n", file)
				DebounceCallback(func() {
					screen.Clear()
					screen.MoveTopLeft()
					callback()
				})
			})
			<-done
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
