package metrics

import (
	"sync"
	"time"
)

type Counter struct {
	mu     sync.RWMutex
	values map[string]uint64
}

func NewCounter() *Counter                   { return &Counter{values: map[string]uint64{}} }
func (c *Counter) Add(name string, n uint64) { c.mu.Lock(); c.values[name] += n; c.mu.Unlock() }
func (c *Counter) Inc(name string)           { c.Add(name, 1) }
func (c *Counter) Get(name string) uint64    { c.mu.RLock(); defer c.mu.RUnlock(); return c.values[name] }
func (c *Counter) Snapshot() map[string]uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := map[string]uint64{}
	for k, v := range c.values {
		out[k] = v
	}
	return out
}

type Timer struct{ started time.Time }

func StartTimer() Timer                 { return Timer{started: time.Now()} }
func (t Timer) Duration() time.Duration { return time.Since(t.started) }
