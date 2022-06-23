package hooks

import (
	"time"
)

var WillCallCallback bool = false

func DebounceCallback(callback func()) {
	if WillCallCallback {
		return
	}
	WillCallCallback = true
	go func() {
		time.Sleep(1 * time.Second)
		WillCallCallback = false
		callback()
	}()
}
