package rules

import (
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Action interface{ Execute(domain.Alarm) error }
type Callback func(domain.Alarm) error

func (c Callback) Execute(a domain.Alarm) error { return c(a) }

type Engine struct {
	mu         sync.Mutex
	windows    map[domain.ID][]domain.Reading
	states     map[domain.ID]domain.RuleState
	lastAction map[domain.ID]time.Time
	actions    map[domain.ID][]Action
}

func NewEngine() *Engine {
	return &Engine{windows: map[domain.ID][]domain.Reading{}, states: map[domain.ID]domain.RuleState{}, lastAction: map[domain.ID]time.Time{}, actions: map[domain.ID][]Action{}}
}
func (e *Engine) AddAction(id domain.ID, a Action) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.actions[id] = append(e.actions[id], a)
}
func (e *Engine) Evaluate(rule domain.Rule, reading domain.Reading) (domain.Alarm, bool, error) {
	if err := rule.Validate(); err != nil {
		return domain.Alarm{}, false, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	w := append(e.windows[rule.ID], reading)
	if len(w) > rule.Window {
		w = w[len(w)-rule.Window:]
	}
	e.windows[rule.ID] = w
	state := e.states[rule.ID]
	if state == "" {
		state = domain.RuleNormal
	}
	trigger, err := evaluate(rule.Expression, w)
	if err != nil {
		return domain.Alarm{}, false, err
	}
	alarm := domain.Alarm{ID: domain.ID("alarm-" + string(rule.ID)), RuleID: rule.ID, PointID: reading.PointID, State: state, LastValue: reading.Value, Count: len(w)}
	if rule.Maintenance {
		alarm.State = domain.RuleSuppressed
		return alarm, true, nil
	}
	if trigger && state != domain.RuleAlarm {
		alarm.State = domain.RulePending
		if len(w) >= rule.Window {
			alarm.State = domain.RuleAlarm
			alarm.TriggeredAt = time.Now()
			e.states[rule.ID] = alarm.State
		}
	} else if !trigger && state == domain.RuleAlarm && reading.Value < rule.Threshold-rule.Hysteresis {
		alarm.State = domain.RuleCleared
		t := time.Now()
		alarm.ClearedAt = &t
		e.states[rule.ID] = alarm.State
	}
	return alarm, true, nil
}
func evaluate(expr string, w []domain.Reading) (bool, error) {
	if len(w) == 0 {
		return false, nil
	}
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid expression")
	}
	v, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return false, err
	}
	current := w[len(w)-1].Value
	switch parts[1] {
	case ">":
		return current > v, nil
	case ">=":
		return current >= v, nil
	case "<":
		return current < v, nil
	case "<=":
		return current <= v, nil
	case "==":
		return math.Abs(current-v) < 1e-9, nil
	default:
		return false, fmt.Errorf("unsupported operator")
	}
}
func (e *Engine) Snapshot() map[domain.ID]domain.RuleState {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := map[domain.ID]domain.RuleState{}
	for k, v := range e.states {
		out[k] = v
	}
	return out
}
