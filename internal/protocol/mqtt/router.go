package mqtt

import (
	"strings"
	"sync"
	"time"
)

type Message struct {
	Topic      string
	Payload    []byte
	QoS        byte
	Retain     bool
	ReceivedAt time.Time
	MessageID  string
}
type Handler func(Message) error
type Router struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	seen     map[string]time.Time
}

func NewRouter() *Router {
	return &Router{handlers: map[string][]Handler{}, seen: map[string]time.Time{}}
}
func (r *Router) Subscribe(pattern string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[pattern] = append(r.handlers[pattern], h)
}
func match(pattern, topic string) bool {
	pp := strings.Split(pattern, "/")
	tt := strings.Split(topic, "/")
	for i, p := range pp {
		if p == "#" {
			return i <= len(tt)
		}
		if i >= len(tt) {
			return false
		}
		if p != "+" && p != tt[i] {
			return false
		}
	}
	return len(pp) == len(tt)
}
func (r *Router) Publish(m Message) []error {
	r.mu.Lock()
	if m.MessageID != "" {
		if _, ok := r.seen[m.MessageID]; ok {
			r.mu.Unlock()
			return nil
		}
		r.seen[m.MessageID] = time.Now()
	}
	hs := []Handler{}
	for p, v := range r.handlers {
		if match(p, m.Topic) {
			hs = append(hs, v...)
		}
	}
	r.mu.Unlock()
	errs := []error{}
	for _, h := range hs {
		if err := h(m); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func (r *Router) Prune(ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cut := time.Now().Add(-ttl)
	for id, t := range r.seen {
		if t.Before(cut) {
			delete(r.seen, id)
		}
	}
}
