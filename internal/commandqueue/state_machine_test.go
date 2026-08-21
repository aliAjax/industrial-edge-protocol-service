package commandqueue

import (
	"testing"
	"time"

	"industrial-edge-protocol/internal/domain"
)

func TestCommandRetryUsesLegalTransition(t *testing.T) {
	if !domain.CommandConfirmed.Terminal() {
		t.Fatal("confirmed command must be terminal")
	}
	if domain.CommandSent.Terminal() || !domain.CommandSent.CanRetry() {
		t.Fatal("sent command must remain retryable and non-terminal")
	}
}

func TestCommandRetrySuccessReachesConfirmed(t *testing.T) {
	machine := domain.NewMachine(string(domain.CommandSent), []domain.Transition{
		{From: string(domain.CommandSent), Event: "reject", To: string(domain.CommandRejected)},
		{From: string(domain.CommandRejected), Event: "retry", To: string(domain.CommandRetrying)},
		{From: string(domain.CommandRetrying), Event: "confirm", To: string(domain.CommandConfirmed)},
	})
	for _, event := range []string{"reject", "retry", "confirm"} {
		if err := machine.Apply(event); err != nil {
			t.Fatalf("apply %s: %v", event, err)
		}
	}
	if machine.State() != string(domain.CommandConfirmed) {
		t.Fatalf("unexpected terminal state: %s", machine.State())
	}
}

func TestConfirmedCommandLeavesInflight(t *testing.T) {
	queue := New()
	command := domain.Command{ID: "cmd-1", Status: domain.CommandApproved, ExpiresAt: time.Now().Add(time.Minute)}
	if err := queue.Enqueue(command); err != nil {
		t.Fatal(err)
	}
	leased, ok := queue.Next(time.Now())
	if !ok {
		t.Fatal("command was not leased")
	}
	if err := queue.Confirm(leased.ID, true); err != nil {
		t.Fatal(err)
	}
	status, ok := queue.Status(leased.ID)
	if !ok || status != domain.CommandConfirmed {
		t.Fatalf("unexpected persisted status: %s", status)
	}
	if len(queue.inflight) != 0 {
		t.Fatalf("confirmed command remains inflight: %d", len(queue.inflight))
	}
}

func TestRetryResetClearsIntermediateState(t *testing.T) {
	retry := &Retry{Attempts: 3, Max: 5, Next: time.Now().Add(time.Hour), Backoff: time.Second}
	retry.Reset()
	if retry.Attempts != 0 || !retry.Next.IsZero() {
		t.Fatalf("retry state not fully reset: %+v", retry)
	}
}
