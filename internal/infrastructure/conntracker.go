package infrastructure

import (
	"sync"
	"time"
)

// ConnectionTracker tracks active SSE connections and signals
// when all connections have been closed for a grace period.
type ConnectionTracker struct {
	mu             sync.Mutex
	count          int
	everConnected  bool
	timer          *time.Timer
	grace          time.Duration
	done           chan struct{}
	closeOnce      sync.Once
}

// NewConnectionTracker creates a new ConnectionTracker.
// The grace duration controls how long to wait after the last connection
// closes before signaling shutdown.
func NewConnectionTracker(grace time.Duration) *ConnectionTracker {
	return &ConnectionTracker{
		grace: grace,
		done:  make(chan struct{}),
	}
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

// Remove decrements the connection count.
// If count reaches 0 and at least one connection was made before,
// a grace period timer starts. When it fires, Done() channel is closed.
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

// Done returns a channel that is closed when all connections have been
// gone for the grace period.
func (ct *ConnectionTracker) Done() <-chan struct{} {
	return ct.done
}
