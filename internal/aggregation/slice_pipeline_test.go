package aggregation_test

import (
	"reflect"
	"testing"
	"time"

	"industrial-edge-protocol/internal/aggregation"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/rollup"
	"industrial-edge-protocol/internal/sampling"
	"industrial-edge-protocol/internal/telemetry"
)

func TestResamplerDoesNotMutateInput(t *testing.T) {
	base := time.Unix(100, 0)
	values := []domain.Reading{
		{PointID: "p1", Value: 2, ObservedAt: base.Add(time.Second)},
		{PointID: "p1", Value: 1, ObservedAt: base},
	}
	want := append([]domain.Reading(nil), values...)
	result := (sampling.Resampler{Interval: time.Second}).Apply(values)
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("resampler reordered caller input: %+v", values)
	}
	if len(result) != 2 || result[0].ObservedAt.After(result[1].ObservedAt) {
		t.Fatalf("resampler emitted buckets out of order: %+v", result)
	}
}

func TestAggregationKeepsIndependentBuckets(t *testing.T) {
	base := time.Unix(200, 0)
	values := []domain.Reading{
		{Value: 7, Quality: domain.QualityGood, ObservedAt: base.Add(2 * time.Second)},
		{Value: 3, Quality: domain.QualityGood, ObservedAt: base},
	}
	want := append([]domain.Reading(nil), values...)
	buckets := aggregation.Build(values, base, base.Add(4*time.Second), time.Second)
	if len(buckets) != 4 || !reflect.DeepEqual(values, want) {
		t.Fatalf("aggregation reused or reordered its input: buckets=%d values=%+v", len(buckets), values)
	}
}

func TestRollupSnapshotSurvivesAppend(t *testing.T) {
	window := rollup.New(time.Minute)
	window.Add(rollup.Sample{At: time.Unix(300, 0), Value: 10, Valid: true})
	snapshot := window.Snapshot()
	snapshot[0].Value = 99
	again := window.Snapshot()
	if again[0].Value != 10 {
		t.Fatalf("snapshot mutation leaked into window: %+v", again)
	}
}

func TestTelemetryFilterBoundaryDoesNotPanic(t *testing.T) {
	values := []domain.Reading{{PointID: "p1", Value: 4, ObservedAt: time.Unix(400, 0)}}
	filtered := telemetry.FilterReadings(values, telemetry.Filter{PointID: "missing"})
	if len(filtered) != 0 {
		t.Fatalf("expected empty filtered result: %+v", filtered)
	}
	if values[0].PointID != "p1" {
		t.Fatalf("filter rewrote caller input: %+v", values)
	}
	values = []domain.Reading{
		{PointID: "skip", Value: 1, ObservedAt: time.Unix(401, 0)},
		{PointID: "keep", Value: 2, ObservedAt: time.Unix(402, 0)},
	}
	filtered = telemetry.FilterReadings(values, telemetry.Filter{PointID: "keep"})
	if len(filtered) != 1 || values[0].PointID != "skip" {
		t.Fatalf("filter reused caller backing array: result=%+v input=%+v", filtered, values)
	}
}
