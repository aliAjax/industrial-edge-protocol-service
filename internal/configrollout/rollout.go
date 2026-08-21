package configrollout

import (
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

type Manager struct {
	mu               sync.Mutex
	rollouts         map[domain.ID]domain.Rollout
	acknowledgements map[domain.ID]map[domain.ID]bool
}

func NewManager() *Manager {
	return &Manager{rollouts: map[domain.ID]domain.Rollout{}, acknowledgements: map[domain.ID]map[domain.ID]bool{}}
}
func (m *Manager) Start(r domain.Rollout) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.Percent < 1 || r.Percent > 100 {
		return fmt.Errorf("invalid percentage")
	}
	if _, ok := m.rollouts[r.ID]; ok {
		return fmt.Errorf("rollout exists")
	}
	if m.rollouts == nil {
		m.rollouts = map[domain.ID]domain.Rollout{}
	}
	if m.acknowledgements == nil {
		m.acknowledgements = map[domain.ID]map[domain.ID]bool{}
	}
	r.State = "running"
	r.CreatedAt = time.Now().UTC()
	m.rollouts[r.ID] = r
	m.acknowledgements[r.ID] = map[domain.ID]bool{}
	return nil
}
func (m *Manager) Ack(id, gateway domain.ID, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.acknowledgements[id]; !exists {
		return
	}
	m.acknowledgements[id][gateway] = ok
	m.reconcileLocked(id)
}
func (m *Manager) reconcileLocked(id domain.ID) {
	r := m.rollouts[id]
	total := len(r.GatewayIDs)
	yes := 0
	for _, g := range r.GatewayIDs {
		if m.acknowledgements[id][g] {
			yes++
		}
	}
	if total > 0 && yes == total {
		r.State = "completed"
		m.rollouts[id] = r
	}
}
func (m *Manager) Get(id domain.ID) (domain.Rollout, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rollouts[id]
	return v, ok
}
func (m *Manager) Abort(id domain.ID, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.rollouts[id]; ok {
		r.State = "aborted"
		r.Error = reason
		m.rollouts[id] = r
	}
}
