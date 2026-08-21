package replay

import (
	"context"
	"errors"
	"testing"
	"time"

	"industrial-edge-protocol/internal/domain"
)

func TestReplayRollbackKeepsCheckpoint(t *testing.T) {
	sentinel := errors.New("gateway upload rejected")
	ctx, cancel := context.WithCancel(context.Background())
	err := Run(ctx, []domain.Reading{{PointID: "p1", ObservedAt: time.Now()}}, func(context.Context, domain.Reading) error {
		cancel()
		return sentinel
	}, 0)
	if !errors.Is(err, sentinel) {
		t.Fatalf("upload error was overwritten: %v", err)
	}
}
