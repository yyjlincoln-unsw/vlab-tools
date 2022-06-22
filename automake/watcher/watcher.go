package watcher

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

type onChangeHandler = func(string)

type TERMINATE_STATUS = int

const (
	TERMINATE_ROUTINE  int = 0
	ROUTINE_TERMINATED int = 1
)

type watcher struct {
	handlers map[string]onChangeHandler
	watcher  *fsnotify.Watcher

	// Routine
	routineStarted   bool
	terminateRoutine chan TERMINATE_STATUS // Setting any value here will terminate the routine
}

func (w *watcher) pullChanges() {
	defer func() {
		w.routineStarted = false
		w.terminateRoutine <- ROUTINE_TERMINATED
	}()
	for {
		select {
		case <-w.terminateRoutine:
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op == fsnotify.Write || event.Op == fsnotify.Create || event.Op == fsnotify.Remove {
				for _, fn := range w.handlers {
					fn(event.Name)
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		}
	}
}

func New(files string) (*watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("new watcher: could not init new watcher: %v", err)
	}
	err = fsWatcher.Add(files)
	if err != nil {
		return nil, err
	}
	// Start the routine
	watcher := &watcher{
		handlers:         map[string]func(string){},
		watcher:          fsWatcher,
		routineStarted:   true,
		terminateRoutine: make(chan int),
	}
	go watcher.pullChanges()
	return watcher, nil
}

func (w *watcher) OnChange(handler onChangeHandler) (string, error) {
	id := uuid.New().String()
	if handler := w.handlers[id]; handler != nil {
		return "", fmt.Errorf("uuid was not unique")
	}
	w.handlers[id] = handler
	return id, nil
}

func (w *watcher) Destroy() {
	defer func() {
		close(w.terminateRoutine)
	}()
	w.terminateRoutine <- TERMINATE_ROUTINE
	status := <-w.terminateRoutine
	if status != ROUTINE_TERMINATED {
		panic(fmt.Errorf("could not terminate the routine"))
	}
	w.handlers = make(map[string]func(string))
}
