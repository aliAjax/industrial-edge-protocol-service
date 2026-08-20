package limits

import (
	"sync"
	"time"
)

type Bucket struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	updated  time.Time
}

func NewBucket(capacity, rate float64) *Bucket {
	return &Bucket{capacity: capacity, tokens: capacity, rate: rate, updated: time.Now()}
}
func (b *Bucket) Allow(cost float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens += now.Sub(b.updated).Seconds() * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.updated = now
	if b.tokens < cost {
		return false
	}
	b.tokens -= cost
	return true
}
func (b *Bucket) Available() float64 { b.mu.Lock(); defer b.mu.Unlock(); return b.tokens }
func (b *Bucket) Reset()             { b.mu.Lock(); b.tokens = b.capacity; b.updated = time.Now(); b.mu.Unlock() }
