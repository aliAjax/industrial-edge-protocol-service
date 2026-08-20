package opcua

import (
	"context"
	"errors"
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

var ErrSession = errors.New("opc ua session unavailable")

type Node struct {
	ID          string
	DisplayName string
	DataType    string
	Writable    bool
}
type Value struct {
	NodeID     string
	Value      any
	Quality    domain.Quality
	SourceTime time.Time
}
type Session interface {
	Browse(context.Context, string) ([]Node, error)
	Read(context.Context, []string) ([]Value, error)
	Subscribe(context.Context, []string, chan<- Value) error
	Close() error
}
type MockSession struct {
	mu     sync.Mutex
	values map[string]Value
	closed bool
}

func NewMock() *MockSession { return &MockSession{values: map[string]Value{}} }
func (m *MockSession) Browse(context.Context, string) ([]Node, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrSession
	}
	out := make([]Node, 0, len(m.values))
	for id := range m.values {
		out = append(out, Node{ID: id, DisplayName: id, DataType: "float64"})
	}
	return out, nil
}
func (m *MockSession) Read(_ context.Context, ids []string) ([]Value, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrSession
	}
	out := []Value{}
	for _, id := range ids {
		if v, ok := m.values[id]; ok {
			out = append(out, v)
		}
	}
	return out, nil
}
func (m *MockSession) Subscribe(ctx context.Context, ids []string, ch chan<- Value) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, id := range ids {
				select {
				case ch <- Value{NodeID: id, Value: 0, Quality: domain.QualityGood, SourceTime: time.Now().UTC()}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}
func (m *MockSession) Close() error { m.mu.Lock(); m.closed = true; m.mu.Unlock(); return nil }
func (m *MockSession) Set(id string, v any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[id] = Value{NodeID: id, Value: v, Quality: domain.QualityGood, SourceTime: time.Now().UTC()}
}
