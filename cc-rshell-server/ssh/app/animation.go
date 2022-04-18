package app

import (
	"github.com/rivo/tview"
	"sync"
	"time"
)

var globalAnimationMutex = &sync.Mutex{}

// Animate triggers the application to redraw every 50ms
func Animate(stop <-chan struct{}, app *tview.Application) {
	if !globalAnimationMutex.TryLock() {
		return
	}
	defer globalAnimationMutex.Unlock()

	for {
		select {
		case <-stop:
			return
		default:
			app.QueueUpdateDraw(func() {})
			time.Sleep(100 * time.Millisecond)
		}
	}
}
