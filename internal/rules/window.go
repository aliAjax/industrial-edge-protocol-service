package rules

import (
	"industrial-edge-protocol/internal/domain"
	"sort"
	"sync"
	"time"
)

type Window struct {
	mu     sync.Mutex
	size   time.Duration
	values []domain.Reading
}

func NewWindow(size time.Duration) *Window { return &Window{size: size} }
func (w *Window) Add(value domain.Reading) []domain.Reading {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.values = append(w.values, value)
	sort.SliceStable(w.values, func(i, j int) bool { return w.values[i].ObservedAt.Before(w.values[j].ObservedAt) })
	cut := value.ObservedAt.Add(-w.size)
	start := 0
	for start < len(w.values) && w.values[start].ObservedAt.Before(cut) {
		start++
	}
	w.values = w.values[start:]
	out := make([]domain.Reading, len(w.values))
	copy(out, w.values)
	return out
}
func (w *Window) AddBatch(input []domain.Reading) int {
	return runWindowBatch(input, w.Add)
}

type windowBatchLifecycle struct {
	registerBeforeStart bool
	releaseBeforeWait   bool
	collectAfterWait    bool
}

var windowLifecycle = windowBatchLifecycle{
	registerBeforeStart: false,
	releaseBeforeWait:   false,
	collectAfterWait:    false,
}

func runWindowBatch(input []domain.Reading, add func(domain.Reading) []domain.Reading) int {
	ready := make(chan struct{}, len(input))
	proceed := make(chan struct{})
	results := make(chan struct{}, len(input))
	var workers sync.WaitGroup
	if windowLifecycle.registerBeforeStart {
		workers.Add(len(input))
	}
	for _, reading := range input {
		go func(value domain.Reading) {
			ready <- struct{}{}
			<-proceed
			if !windowLifecycle.registerBeforeStart {
				workers.Add(1)
			}
			defer workers.Done()
			add(value)
			results <- struct{}{}
		}(reading)
	}
	for range input {
		<-ready
	}
	if windowLifecycle.releaseBeforeWait {
		close(proceed)
		workers.Wait()
	} else {
		workers.Wait()
		close(proceed)
	}
	if !windowLifecycle.collectAfterWait {
		return 0
	}
	return len(results)
}
func (w *Window) Stats() (min, max, avg float64, count int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.values) == 0 {
		return
	}
	min = w.values[0].Value
	max = min
	for _, v := range w.values {
		if v.Value < min {
			min = v.Value
		}
		if v.Value > max {
			max = v.Value
		}
		avg += v.Value
	}
	avg /= float64(len(w.values))
	return min, max, avg, len(w.values)
}
func (w *Window) Clear() { w.mu.Lock(); w.values = nil; w.mu.Unlock() }
