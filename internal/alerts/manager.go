package alerts

import (
	"industrial-edge-protocol/internal/domain"
	"sync"
	"time"
)

type Manager struct {
	mu      sync.Mutex
	alarms  map[domain.ID]domain.Alarm
	history []domain.Alarm
}

func New() *Manager { return &Manager{alarms: map[domain.ID]domain.Alarm{}} }
func (m *Manager) Open(a domain.Alarm) {
	m.mu.Lock()
	m.alarms[a.ID] = a
	m.history = append(m.history, a)
	m.mu.Unlock()
}
func (m *Manager) Clear(id domain.ID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if a, ok := m.alarms[id]; ok {
		now := time.Now()
		a.State = domain.RuleCleared
		a.ClearedAt = &now
		m.alarms[id] = a
		m.history = append(m.history, a)
	}
}
func (m *Manager) Active() []domain.Alarm {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []domain.Alarm{}
	for _, a := range m.alarms {
		if a.State == domain.RuleAlarm || a.State == domain.RulePending {
			out = append(out, a)
		}
	}
	return out
}
func (m *Manager) History() []domain.Alarm {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]domain.Alarm, len(m.history))
	copy(out, m.history)
	return out
}
