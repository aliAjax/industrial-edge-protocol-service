package acquisition_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"industrial-edge-protocol/internal/acquisition"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/health"
	"industrial-edge-protocol/internal/worker"
)

type waitingReader struct{}

func (waitingReader) Read(ctx context.Context, _ domain.Point) (domain.Reading, error) {
	<-ctx.Done()
	return domain.Reading{}, ctx.Err()
}

func TestAcquisitionDeadlineReachesReader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	_, err := acquisition.ReadWithRetry(ctx, waitingReader{}, []domain.Point{{ID: "p1"}}, 200*time.Millisecond, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected read error: %v", err)
	}
	if elapsed := time.Since(started); elapsed > 80*time.Millisecond {
		t.Fatalf("parent cancellation did not reach reader promptly: %v", elapsed)
	}
}

func TestAcquisitionCancelStopsRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (acquisition.MockReader{Seed: 4}).Read(ctx, domain.Point{ID: "p1", Scale: 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("mock reader ignored cancellation: %v", err)
	}
}

func TestWorkerDoesNotReuseExpiredContext(t *testing.T) {
	pool := worker.New(context.Background(), 1)
	defer pool.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	observed := make(chan error, 1)
	job := func(jobCtx context.Context) error {
		observed <- jobCtx.Err()
		return nil
	}
	deadline := time.Now().Add(time.Second)
	for {
		if err := pool.SubmitContext(ctx, job); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("worker never accepted job")
		}
		time.Sleep(time.Millisecond)
	}
	select {
	case err := <-observed:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("worker used pool context instead of request context: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("worker did not run job")
	}
}

func TestHealthCheckHonorsParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results := health.Run(ctx, []health.Check{{
		Name: "gateway-link",
		Fn:   func(checkCtx context.Context) error { return checkCtx.Err() },
	}})
	if len(results) != 1 || results[0].Healthy || results[0].Error == "" {
		t.Fatalf("canceled health check reported healthy: %+v", results)
	}
}
