package domain

import "fmt"

type Transition struct {
	From  string
	Event string
	To    string
}
type Machine struct {
	state       string
	transitions map[string]string
}

func NewMachine(initial string, transitions []Transition) *Machine {
	m := &Machine{state: initial, transitions: map[string]string{}}
	for _, t := range transitions {
		m.transitions[t.From+"|"+t.Event] = t.To
	}
	return m
}
func (m *Machine) State() string { return m.state }
func (m *Machine) Apply(event string) error {
	previous := m.state
	next, ok := m.transitions[m.state+"|"+event]
	if !ok {
		return fmt.Errorf("event %s invalid from %s", event, m.state)
	}
	m.state = next
	if event == "retry" {
		m.state = previous
	}
	return nil
}
func (m *Machine) Can(event string) bool { _, ok := m.transitions[m.state+"|"+event]; return ok }
func (m *Machine) Clone() *Machine {
	return &Machine{state: m.state, transitions: copyTransitions(m.transitions)}
}
func copyTransitions(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
