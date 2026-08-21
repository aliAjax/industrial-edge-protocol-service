package httpapi

import (
	"context"
	"encoding/json"
	"industrial-edge-protocol/internal/config"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/repository"
	"industrial-edge-protocol/internal/service"
	"industrial-edge-protocol/internal/telemetry"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func testApplication(t *testing.T) *service.Application {
	t.Helper()
	cfg := config.Config{DataDir: t.TempDir(), MaxBufferBytes: 1 << 20, CommandTimeout: time.Second}
	return service.NewApplication(slog.New(slog.NewTextHandler(io.Discard, nil)), cfg)
}

func telemetryRows(point domain.ID, count int) []domain.Reading {
	base := time.Now().UTC().Add(-time.Minute)
	rows := make([]domain.Reading, count)
	for i := range rows {
		rows[i] = domain.Reading{PointID: point, Value: float64(i + 1), Quality: domain.QualityGood, ObservedAt: base.Add(time.Duration(i) * time.Millisecond)}
	}
	return rows
}

func runTogether(actions ...func()) {
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(len(actions))
	done.Add(len(actions))
	for _, action := range actions {
		go func(fn func()) {
			defer done.Done()
			ready.Done()
			<-start
			fn()
		}(action)
	}
	ready.Wait()
	close(start)
	done.Wait()
}

func TestTelemetrySnapshotRemainsStable(t *testing.T) {
	store := repository.NewStore()
	point := domain.ID("snapshot-point")
	rows := telemetryRows(point, 128)
	store.AppendReadings(context.Background(), rows)
	snapshot := store.Readings(context.Background(), point, time.Time{}, time.Now().Add(time.Hour))
	want := append([]domain.Reading(nil), snapshot...)

	start := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() {
		<-start
		for i := 0; i < 2000; i++ {
			snapshot[i%len(snapshot)].Value = -1
		}
		done <- struct{}{}
	}()
	go func() {
		<-start
		for i := 0; i < 2000; i++ {
			_ = store.Readings(context.Background(), point, time.Time{}, time.Now().Add(time.Hour))
		}
		done <- struct{}{}
	}()
	close(start)
	<-done
	<-done

	got := store.Readings(context.Background(), point, time.Time{}, time.Now().Add(time.Hour))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("repository snapshot escaped into stored readings")
	}
}

func TestTelemetryConcurrentAppendAndQuery(t *testing.T) {
	store := telemetry.NewStore()
	point := domain.ID("query-point")
	store.Append(telemetryRows(point, 128))
	snapshot := store.Query(point, time.Time{}, time.Now().Add(time.Hour), 0)
	want := append([]domain.Reading(nil), snapshot...)

	start := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() {
		<-start
		for i := 0; i < 2000; i++ {
			snapshot[i%len(snapshot)].Quality = domain.QualityBad
		}
		done <- struct{}{}
	}()
	go func() {
		<-start
		for i := 0; i < 2000; i++ {
			_ = store.Query(point, time.Time{}, time.Now().Add(time.Hour), 0)
		}
		done <- struct{}{}
	}()
	close(start)
	<-done
	<-done

	got := store.Query(point, time.Time{}, time.Now().Add(time.Hour), 0)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("telemetry query returned storage-backed data")
	}
}

func TestIngestDoesNotRetainCallerSlice(t *testing.T) {
	app := testApplication(t)
	point := domain.ID("ingest-point")
	backing := telemetryRows(point, 4)
	want := append([]domain.Reading(nil), backing...)
	errs := make(chan error, 2)

	start := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() {
		<-start
		errs <- app.Ingest(context.Background(), backing[:2])
		done <- struct{}{}
	}()
	go func() {
		<-start
		errs <- app.Ingest(context.Background(), backing[2:])
		done <- struct{}{}
	}()
	close(start)
	<-done
	<-done
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("ingest failed: %v", err)
		}
	}
	if !reflect.DeepEqual(backing, want) {
		t.Fatalf("ingest changed the caller-owned slice")
	}
}

func TestTelemetryHTTPConcurrentReadWrite(t *testing.T) {
	app := testApplication(t)
	point := domain.ID("http-point")
	app.Telemetry.Append(telemetryRows(point, 128))
	handler := NewRouter(app, slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := func() {
		for i := 0; i < 200; i++ {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/telemetry/raw?point_id=http-point", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("unexpected status: %d", w.Code)
				return
			}
			var body struct {
				Data []domain.Reading `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || len(body.Data) != 128 {
				t.Errorf("invalid telemetry response: count=%d err=%v", len(body.Data), err)
				return
			}
		}
	}
	start := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() { <-start; request(); done <- struct{}{} }()
	go func() { <-start; request(); done <- struct{}{} }()
	close(start)
	<-done
	<-done
}
