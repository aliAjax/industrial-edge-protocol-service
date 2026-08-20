package diagnostics

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Span struct {
	ID      string
	Parent  string
	Name    string
	Started time.Time
	Ended   time.Time
	Attrs   map[string]string
}
type Tracer struct {
	mu    sync.Mutex
	spans []Span
}

func NewTracer() *Tracer { return &Tracer{} }
func (t *Tracer) Start(parent, name string) Span {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return Span{ID: hex.EncodeToString(b), Parent: parent, Name: name, Started: time.Now(), Attrs: map[string]string{}}
}
func (t *Tracer) End(span Span) {
	span.Ended = time.Now()
	t.mu.Lock()
	t.spans = append(t.spans, span)
	t.mu.Unlock()
}
func (t *Tracer) Spans() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Span, len(t.spans))
	copy(out, t.spans)
	return out
}
