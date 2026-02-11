package server

import (
	"sync"
	"time"
)

// ConnectionTracker tracks active SSE connections and signals shutdown readiness.
type ConnectionTracker struct {
	mu            sync.Mutex
	count         int
	everConnected bool
	timer         *time.Timer
	grace         time.Duration
	done          chan struct{}
	closeOnce     sync.Once
}

// NewConnectionTracker creates a new ConnectionTracker.
func NewConnectionTracker(grace time.Duration) *ConnectionTracker {
	return &ConnectionTracker{grace: grace, done: make(chan struct{})}
}

// Add increments the connection count.
func (ct *ConnectionTracker) Add() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.count++
	ct.everConnected = true
	if ct.timer != nil {
		ct.timer.Stop()
		ct.timer = nil
	}
}

// Remove decrements the connection count and starts grace timer if no active connections remain.
func (ct *ConnectionTracker) Remove() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.count--
	if ct.count <= 0 {
		ct.count = 0
		if ct.everConnected {
			ct.timer = time.AfterFunc(ct.grace, func() {
				ct.closeOnce.Do(func() {
					close(ct.done)
				})
			})
		}
	}
}

// Done returns a channel that is closed when shutdown can proceed.
func (ct *ConnectionTracker) Done() <-chan struct{} {
	return ct.done
}
