package acquisition

import (
	"context"
	"errors"
	"industrial-edge-protocol/internal/domain"
	"sort"
	"sync"
	"time"
)

var ErrClosed = errors.New("scheduler closed")

type Reader interface {
	Read(context.Context, domain.Point) (domain.Reading, error)
}
type Task struct {
	Plan     domain.AcquisitionPlan
	Due      time.Time
	Attempts int
	Lease    string
}
type Scheduler struct {
	mu     sync.Mutex
	queue  []Task
	leases map[string]Task
	closed bool
	reader Reader
}

func NewScheduler(reader Reader) *Scheduler {
	return &Scheduler{reader: reader, leases: map[string]Task{}}
}
func (s *Scheduler) Add(plan domain.AcquisitionPlan) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.queue = append(s.queue, Task{Plan: plan, Due: time.Now()})
	sort.SliceStable(s.queue, func(i, j int) bool { return s.queue[i].Due.Before(s.queue[j].Due) })
}
func (s *Scheduler) Lease(now time.Time, worker string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return Task{}, ErrClosed
	}
	for i, t := range s.queue {
		if !t.Due.After(now) {
			t.Lease = worker
			s.leases[worker] = t
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return t, nil
		}
	}
	return Task{}, context.DeadlineExceeded
}
func (s *Scheduler) Ack(worker string, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.leases[worker]
	if !ok {
		return
	}
	delete(s.leases, worker)
	if !success && t.Attempts < 3 {
		t.Attempts++
		t.Due = time.Now().Add(time.Duration(t.Attempts) * time.Second)
		s.queue = append(s.queue, t)
	}
}
func (s *Scheduler) Recover(worker string) { s.Ack(worker, false) }
func (s *Scheduler) Close()                { s.mu.Lock(); s.closed = true; s.mu.Unlock() }
