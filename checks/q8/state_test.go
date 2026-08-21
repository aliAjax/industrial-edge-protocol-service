package q8_test

import (
	"industrial-edge-protocol/internal/commandqueue"
	"industrial-edge-protocol/internal/domain"
	"testing"
	"time"
)

func TestQ8Legal(t *testing.T) {
	if !domain.CommandConfirmed.Terminal() || domain.CommandSent.Terminal() || !domain.CommandSent.CanRetry() {
		t.Fatal("status")
	}
}
func TestQ8Transition(t *testing.T) {
	m := domain.NewMachine(string(domain.CommandSent), []domain.Transition{{From: string(domain.CommandSent), Event: "reject", To: string(domain.CommandRejected)}, {From: string(domain.CommandRejected), Event: "retry", To: string(domain.CommandRetrying)}, {From: string(domain.CommandRetrying), Event: "confirm", To: string(domain.CommandConfirmed)}})
	for _, e := range []string{"reject", "retry", "confirm"} {
		if m.Apply(e) != nil {
			t.Fatal(e)
		}
	}
	if m.State() != string(domain.CommandConfirmed) {
		t.Fatal(m.State())
	}
}
func TestQ8Inflight(t *testing.T) {
	q := commandqueue.New()
	if q.Enqueue(domain.Command{ID: "cmd-1", Status: domain.CommandApproved, ExpiresAt: time.Now().Add(time.Minute)}) != nil {
		t.Fatal("enqueue")
	}
	c, ok := q.Next(time.Now())
	if !ok || q.Confirm(c.ID, true) != nil {
		t.Fatal("confirm")
	}
	s, ok := q.Status(c.ID)
	if !ok || s != domain.CommandConfirmed {
		t.Fatal(s)
	}
}
func TestQ8Reset(t *testing.T) {
	r := &commandqueue.Retry{Attempts: 3, Max: 5, Next: time.Now().Add(time.Hour), Backoff: time.Second}
	r.Reset()
	if r.Attempts != 0 || !r.Next.IsZero() {
		t.Fatal("reset")
	}
}
