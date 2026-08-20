package migration

import "sync"

type Runner struct {
	mu      sync.Mutex
	version int
	applied []string
}

func NewRunner() *Runner { return &Runner{} }
func (r *Runner) Apply(name string, fn func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := fn(); err != nil {
		return err
	}
	r.version++
	r.applied = append(r.applied, name)
	return nil
}
func (r *Runner) Version() int { r.mu.Lock(); defer r.mu.Unlock(); return r.version }
func (r *Runner) Applied() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.applied))
	copy(out, r.applied)
	return out
}
