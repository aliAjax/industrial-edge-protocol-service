package repository

import (
	"sync"
	"time"
)

type Lease struct {
	Owner     string
	ExpiresAt time.Time
	Epoch     uint64
}
type LeaseStore struct {
	mu     sync.Mutex
	values map[string]Lease
}

func NewLeaseStore() *LeaseStore { return &LeaseStore{values: map[string]Lease{}} }
func (l *LeaseStore) Acquire(key, owner string, duration time.Duration) (Lease, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	current, ok := l.values[key]
	if ok && current.ExpiresAt.After(now) && current.Owner != owner {
		return Lease{}, false
	}
	current = Lease{Owner: owner, ExpiresAt: now.Add(duration), Epoch: current.Epoch + 1}
	l.values[key] = current
	return current, true
}
func (l *LeaseStore) Renew(key, owner string, duration time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.values[key]
	if !ok || v.Owner != owner || v.ExpiresAt.Before(time.Now()) {
		return false
	}
	v.ExpiresAt = time.Now().Add(duration)
	l.values[key] = v
	return true
}
func (l *LeaseStore) Release(key, owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.values[key]; ok && v.Owner == owner {
		delete(l.values, key)
	}
}
