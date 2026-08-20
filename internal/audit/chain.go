package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Sequence uint64    `json:"sequence"`
	Actor    string    `json:"actor"`
	Action   string    `json:"action"`
	Resource string    `json:"resource"`
	At       time.Time `json:"at"`
	Previous string    `json:"previous"`
	Digest   string    `json:"digest"`
}
type Chain struct {
	mu     sync.Mutex
	events []Event
}

func NewChain() *Chain { return &Chain{events: []Event{}} }
func (c *Chain) Append(actor, action, resource string, meta any) Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := Event{Sequence: uint64(len(c.events) + 1), Actor: actor, Action: action, Resource: resource, At: time.Now().UTC()}
	if len(c.events) > 0 {
		e.Previous = c.events[len(c.events)-1].Digest
	}
	raw, _ := json.Marshal(struct {
		E Event
		M any
	}{e, meta})
	sum := sha256.Sum256(raw)
	e.Digest = hex.EncodeToString(sum[:])
	c.events = append(c.events, e)
	return e
}
func (c *Chain) Verify() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	previous := ""
	for i, e := range c.events {
		if e.Sequence != uint64(i+1) || e.Previous != previous {
			return fmt.Errorf("audit chain discontinuity at %d", e.Sequence)
		}
		previous = e.Digest
	}
	return nil
}
func (c *Chain) Export() []Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Event, len(c.events))
	copy(out, c.events)
	return out
}
