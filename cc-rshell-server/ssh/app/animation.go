package app

import (
	"context"
	"github.com/rivo/tview"
	"sync"
	"time"
)

var globalAnimationMutex = &sync.Mutex{}

// Animate triggers the application to redraw every 100ms
func Animate(ctx context.Context, app *tview.Application) {
	if !globalAnimationMutex.TryLock() {
		return
	}
	defer globalAnimationMutex.Unlock()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			app.QueueUpdateDraw(func() {})
			time.Sleep(100 * time.Millisecond)
		}
	}
}
