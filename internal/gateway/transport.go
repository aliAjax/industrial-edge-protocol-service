package gateway

import (
	"context"
	"errors"
	"industrial-edge-protocol/internal/domain"
	"sync"
)

var ErrTransportClosed = errors.New("gateway transport closed")

type Upload struct {
	GatewayID domain.ID
	SegmentID string
	Sequence  uint64
	Readings  []domain.Reading
	Digest    string
}
type Transport interface {
	UploadSegment(context.Context, Upload) error
	ApplyConfig(context.Context, domain.Gateway, []byte) error
	Heartbeat(context.Context, domain.Gateway) error
}
type MemoryTransport struct {
	mu      sync.Mutex
	closed  bool
	uploads []Upload
}

func NewMemoryTransport() *MemoryTransport { return &MemoryTransport{} }
func (m *MemoryTransport) UploadSegment(_ context.Context, u Upload) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrTransportClosed
	}
	m.uploads = append(m.uploads, u)
	return nil
}
func (m *MemoryTransport) ApplyConfig(_ context.Context, _ domain.Gateway, _ []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrTransportClosed
	}
	return nil
}
func (m *MemoryTransport) Heartbeat(_ context.Context, _ domain.Gateway) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrTransportClosed
	}
	return nil
}
func (m *MemoryTransport) Close() { m.mu.Lock(); m.closed = true; m.mu.Unlock() }
func (m *MemoryTransport) Uploads() []Upload {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Upload, len(m.uploads))
	copy(out, m.uploads)
	return out
}
