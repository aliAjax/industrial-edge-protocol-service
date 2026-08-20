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
func (b *Bus) Count(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}
