package types

import (
	"bytes"
	"context"
	"sync"
	"time"
)

type TimedBuffer struct {
	Buffer           bytes.Buffer
	LastModification time.Time
	Lock             sync.Mutex
}

func NewTimedBuffer() *TimedBuffer {
	return &TimedBuffer{
		Buffer:           bytes.Buffer{},
		LastModification: time.Now(),
		Lock:             sync.Mutex{},
	}
}

func PurgeClientTimedBuffers(registry ClientRegistry, bufferLifetime time.Duration) {
	purgeBefore := time.Now().Add(-bufferLifetime)

	for _, descriptor := range registry {
		bufferMap := descriptor.MessageBufferMap()
		PurgeTimedBuffers(bufferMap, purgeBefore)
	}
}

func PurgeTimedBuffers(bufferMap BufferMap, purgeBefore time.Time) {
	for _, tb := range bufferMap {
		tb.Lock.Lock()
		if tb.Buffer.Len() > 0 && tb.LastModification.Before(purgeBefore) {
			tb.Buffer.Reset()
			tb.Buffer = bytes.Buffer{}
			tb.LastModification = time.Now()
		}
		tb.Lock.Unlock()
	}
}

func PurgeJob(ctx context.Context, ticker *time.Ticker, registry ClientRegistry, bufferLifetime time.Duration) {
	select {
	case <-ctx.Done():
	case <-ticker.C:
		PurgeClientTimedBuffers(registry, bufferLifetime)
	}
}
