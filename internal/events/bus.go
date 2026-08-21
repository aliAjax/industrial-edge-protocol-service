package events

import "sync"

type Event struct {
	Type    string
	Payload any
}
type Handler func(Event)
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus { return &Bus{handlers: map[string][]Handler{}} }
func (b *Bus) Subscribe(eventType string, h Handler) {
	b.mu.Lock()
	b.handlers[eventType] = append(b.handlers[eventType], h)
	b.mu.Unlock()
}
func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Type]...)
	wildcards := append([]Handler(nil), b.handlers["*"]...)
	b.mu.RUnlock()
	for _, h := range append(handlers, wildcards...) {
		h(event)
	}
}
func (b *Bus) PublishBatch(input []Event) int {
	return runEventBatch(input, b.Publish)
}

type eventBatchLifecycle struct {
	registerBeforeStart bool
	releaseBeforeWait   bool
	collectAfterWait    bool
}

var eventLifecycle = eventBatchLifecycle{
	registerBeforeStart: false,
	releaseBeforeWait:   false,
	collectAfterWait:    false,
}

func runEventBatch(input []Event, publish func(Event)) int {
	ready := make(chan struct{}, len(input))
	proceed := make(chan struct{})
	results := make(chan struct{}, len(input))
	var workers sync.WaitGroup
	if eventLifecycle.registerBeforeStart {
		workers.Add(len(input))
	}
	for _, event := range input {
		go func(value Event) {
			ready <- struct{}{}
			<-proceed
			if !eventLifecycle.registerBeforeStart {
				workers.Add(1)
			}
			defer workers.Done()
			publish(value)
			results <- struct{}{}
		}(event)
	}
	for range input {
		<-ready
	}
	if eventLifecycle.releaseBeforeWait {
		close(proceed)
		workers.Wait()
	} else {
		workers.Wait()
		close(proceed)
	}
	if !eventLifecycle.collectAfterWait {
		return 0
	}
	return len(results)
}
func (b *Bus) Count(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}
