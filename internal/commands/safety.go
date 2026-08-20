package commands

import (
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"strings"
	"sync"
	"time"
)

type Interlock func(domain.Command) error
type AuditEvent struct {
	CommandID domain.ID
	Action    string
	Actor     string
	At        time.Time
	Reason    string
}
type Guard struct {
	mu         sync.Mutex
	interlocks map[string][]Interlock
	audits     []AuditEvent
}

func NewGuard() *Guard { return &Guard{interlocks: map[string][]Interlock{}} }
func (g *Guard) Register(name string, check Interlock) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.interlocks[name] = append(g.interlocks[name], check)
}
func (g *Guard) Check(c domain.Command) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if c.DryRun {
		return nil
	}
	if strings.TrimSpace(c.Interlock) == "" {
		return domain.ErrUnsafeCommand
	}
	for _, check := range g.interlocks[c.Interlock] {
		if err := check(c); err != nil {
			return fmt.Errorf("interlock: %w", err)
		}
	}
	return nil
}
func (g *Guard) Record(id domain.ID, action, actor, reason string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.audits = append(g.audits, AuditEvent{CommandID: id, Action: action, Actor: actor, At: time.Now().UTC(), Reason: reason})
}
func (g *Guard) Audit() []AuditEvent {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]AuditEvent, len(g.audits))
	copy(out, g.audits)
	return out
}
func RequireDeviceEnabled(enabled bool) Interlock {
	return func(domain.Command) error {
		if !enabled {
			return fmt.Errorf("device disabled")
		}
		return nil
	}
}
func RequireRange(min, max float64) Interlock {
	return func(c domain.Command) error {
		if c.Value < min || c.Value > max {
			return fmt.Errorf("value outside range")
		}
		return nil
	}
}
