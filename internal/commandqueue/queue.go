package commandqueue

import (
	"context"
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

type Queue struct {
	mu       sync.Mutex
	pending  []domain.Command
	inflight map[domain.ID]domain.Command
}

func New() *Queue { return &Queue{inflight: map[domain.ID]domain.Command{}} }
func (q *Queue) Enqueue(c domain.Command) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if c.Status != domain.CommandApproved {
		return domain.ErrConflict
	}
	q.pending = append(q.pending, c)
	return nil
}
func (q *Queue) Next(now time.Time) (domain.Command, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, c := range q.pending {
		if c.ExpiresAt.Before(now) {
			c.Status = domain.CommandExpired
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			continue
		}
		c.Status = domain.CommandSent
		q.inflight[c.ID] = c
		q.pending = append(q.pending[:i], q.pending[i+1:]...)
		return c, true
	}
	return domain.Command{}, false
}
func (q *Queue) Confirm(id domain.ID, ok bool) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	c, exists := q.inflight[id]
	if !exists {
		return domain.ErrNotFound
	}
	delete(q.inflight, id)
	if ok {
		c.Status = domain.CommandConfirmed
	} else {
		c.Status = domain.CommandRejected
	}
	return nil
}
func (q *Queue) Drain(ctx context.Context, send func(context.Context, domain.Command) bool) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		c, ok := q.Next(time.Now())
		if !ok {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if !send(ctx, c) {
			_ = q.Confirm(c.ID, false)
		} else {
			_ = q.Confirm(c.ID, true)
		}
	}
}
