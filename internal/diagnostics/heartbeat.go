package diagnostics

import (
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

type Heartbeat struct {
	GatewayID     domain.ID
	At            time.Time
	CPU           float64
	MemoryBytes   uint64
	DiskFreeBytes uint64
	Version       string
	Healthy       bool
}
type Registry struct {
	mu     sync.RWMutex
	values map[domain.ID]Heartbeat
}

func NewRegistry() *Registry { return &Registry{values: map[domain.ID]Heartbeat{}} }
func (r *Registry) Record(h Heartbeat) {
	r.mu.Lock()
	defer r.mu.Unlock()
	h.At = time.Now().UTC()
	h.Healthy = h.CPU < 0.95 && h.DiskFreeBytes > 100<<20
	r.values[h.GatewayID] = h
}
func (r *Registry) Get(id domain.ID) (Heartbeat, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.values[id]
	return h, ok
}
func (r *Registry) Stale(maxAge time.Duration) []domain.ID {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.ID{}
	cut := time.Now().Add(-maxAge)
	for id, h := range r.values {
		if h.At.Before(cut) {
			out = append(out, id)
		}
	}
	return out
}
func (r *Registry) Snapshot() []Heartbeat {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Heartbeat, 0, len(r.values))
	for _, h := range r.values {
		out = append(out, h)
	}
	return out
}
