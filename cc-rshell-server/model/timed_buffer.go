package model

import (
	"bytes"
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
