package repository

import (
	"sync"
	"time"
)

type Entry struct {
	Key       string
	Response  any
	CreatedAt time.Time
	ExpiresAt time.Time
}
type Idempotency struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
}

func NewIdempotency(ttl time.Duration) *Idempotency {
	return &Idempotency{entries: map[string]Entry{}, ttl: ttl}
}
func (i *Idempotency) Get(key string) (any, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	e, ok := i.entries[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		delete(i.entries, key)
		return nil, false
	}
	return e.Response, true
}
func (i *Idempotency) Put(key string, response any) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.entries[key] = Entry{Key: key, Response: response, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(i.ttl)}
}
func (i *Idempotency) Prune() {
	i.mu.Lock()
	defer i.mu.Unlock()
	now := time.Now()
	for k, v := range i.entries {
		if now.After(v.ExpiresAt) {
			delete(i.entries, k)
		}
	}
}
