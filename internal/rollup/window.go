package rollup

import (
	"sync"
	"time"
)

type Sample struct {
	At    time.Time
	Value float64
	Valid bool
}
type Window struct {
	mu     sync.Mutex
	size   time.Duration
	values []Sample
}

func New(size time.Duration) *Window {
	if size <= 0 {
		size = time.Minute
	}
	return &Window{size: size}
}
func (w *Window) Add(s Sample) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.values = append(w.values, s)
	cut := s.At.Add(-w.size)
	start := 0
	for start < len(w.values) && w.values[start].At.Before(cut) {
		start++
	}
	w.values = w.values[start:]
}
func (w *Window) Snapshot() []Sample {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]Sample, len(w.values))
	copy(out, w.values)
	return out
}
func (w *Window) Count() int { w.mu.Lock(); defer w.mu.Unlock(); return len(w.values) }
