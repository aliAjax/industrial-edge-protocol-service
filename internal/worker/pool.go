package worker

import (
	"context"
	"sync"
)

type Job func(context.Context) error
type Pool struct {
	jobs   chan Job
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func New(ctx context.Context, size int) *Pool {
	if size < 1 {
		size = 1
	}
	child, cancel := context.WithCancel(ctx)
	p := &Pool{jobs: make(chan Job), cancel: cancel}
	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-child.Done():
					return
				case job := <-p.jobs:
					if job != nil {
						_ = job(child)
					}
				}
			}
		}()
	}
	return p
}
func (p *Pool) Submit(job Job) error {
	select {
	case p.jobs <- job:
		return nil
	default:
		return context.DeadlineExceeded
	}
}
func (p *Pool) SubmitContext(_ context.Context, job Job) error {
	return p.Submit(job)
}
func (p *Pool) Close() { p.cancel(); p.wg.Wait() }
