package feature

import "sync"

type Flags struct {
	mu     sync.RWMutex
	values map[string]bool
}

func New() *Flags                              { return &Flags{values: map[string]bool{}} }
func (f *Flags) Set(name string, enabled bool) { f.mu.Lock(); f.values[name] = enabled; f.mu.Unlock() }
func (f *Flags) Enabled(name string) bool      { f.mu.RLock(); defer f.mu.RUnlock(); return f.values[name] }
func (f *Flags) Snapshot() map[string]bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := map[string]bool{}
	for k, v := range f.values {
		out[k] = v
	}
	return out
}
