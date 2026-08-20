package tenant

import "sync"

type Quota struct {
	Points   int
	Gateways int
	Bytes    int64
}
type Usage struct {
	Points   int
	Gateways int
	Bytes    int64
}
type Registry struct {
	mu     sync.Mutex
	quotas map[string]Quota
	usage  map[string]Usage
}

func New() *Registry                       { return &Registry{quotas: map[string]Quota{}, usage: map[string]Usage{}} }
func (r *Registry) Set(id string, q Quota) { r.mu.Lock(); r.quotas[id] = q; r.mu.Unlock() }
func (r *Registry) Reserve(id string, u Usage) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	q := r.quotas[id]
	current := r.usage[id]
	if q.Points > 0 && current.Points+u.Points > q.Points {
		return false
	}
	if q.Gateways > 0 && current.Gateways+u.Gateways > q.Gateways {
		return false
	}
	if q.Bytes > 0 && current.Bytes+u.Bytes > q.Bytes {
		return false
	}
	current.Points += u.Points
	current.Gateways += u.Gateways
	current.Bytes += u.Bytes
	r.usage[id] = current
	return true
}
func (r *Registry) Release(id string, u Usage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.usage[id]
	current.Points -= u.Points
	current.Gateways -= u.Gateways
	current.Bytes -= u.Bytes
	if current.Points < 0 {
		current.Points = 0
	}
	if current.Gateways < 0 {
		current.Gateways = 0
	}
	if current.Bytes < 0 {
		current.Bytes = 0
	}
	r.usage[id] = current
}
func (r *Registry) Usage(id string) Usage { r.mu.Lock(); defer r.mu.Unlock(); return r.usage[id] }
