package types

import (
	"bytes"
	"cc-rshell-server/model"
	"context"
	"time"
)

type ClientRegistry map[model.ComputerID]ComputerDescriptor

func (r ClientRegistry) PurgeClientTimedBuffers(bufferLifetime time.Duration) {
	purgeBefore := time.Now().Add(-bufferLifetime)

	for _, descriptor := range r {
		bufferMap := descriptor.MessageBufferMap()
		PurgeTimedBuffers(bufferMap, purgeBefore)
	}
}

func PurgeTimedBuffers(bufferMap model.BufferMap, purgeBefore time.Time) {
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

func (r ClientRegistry) PurgeJob(ctx context.Context, ticker *time.Ticker, bufferLifetime time.Duration) {
	select {
	case <-ctx.Done():
	case <-ticker.C:
		r.PurgeClientTimedBuffers(bufferLifetime)
	}
}
