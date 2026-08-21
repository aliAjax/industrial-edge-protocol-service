package s6_test

import (
	"industrial-edge-protocol/internal/aggregation"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/rollup"
	"industrial-edge-protocol/internal/sampling"
	"industrial-edge-protocol/internal/telemetry"
	"reflect"
	"testing"
	"time"
)

func TestS6Resample(t *testing.T) {
	b := time.Unix(100, 0)
	v := []domain.Reading{{PointID: "p1", Value: 2, ObservedAt: b.Add(time.Second)}, {PointID: "p1", Value: 1, ObservedAt: b}}
	w := append([]domain.Reading(nil), v...)
	if len((sampling.Resampler{Interval: time.Second}).Apply(v)) != 2 || !reflect.DeepEqual(v, w) {
		t.Fatal("resampler")
	}
}
func TestS6Aggregate(t *testing.T) {
	b := time.Unix(200, 0)
	v := []domain.Reading{{Value: 7, Quality: domain.QualityGood, ObservedAt: b.Add(2 * time.Second)}, {Value: 3, Quality: domain.QualityGood, ObservedAt: b}}
	w := append([]domain.Reading(nil), v...)
	if len(aggregation.Build(v, b, b.Add(4*time.Second), time.Second)) != 4 || !reflect.DeepEqual(v, w) {
		t.Fatal("aggregate")
	}
}
func TestS6Window(t *testing.T) {
	w := rollup.New(time.Minute)
	w.Add(rollup.Sample{At: time.Unix(300, 0), Value: 10, Valid: true})
	s := w.Snapshot()
	s[0].Value = 99
	if w.Snapshot()[0].Value != 10 {
		t.Fatal("snapshot")
	}
}
func TestS6Filter(t *testing.T) {
	v := []domain.Reading{{PointID: "p1", ObservedAt: time.Unix(400, 0)}}
	if len(telemetry.FilterReadings(v, telemetry.Filter{PointID: "missing"})) != 0 {
		t.Fatal("empty")
	}
	v = []domain.Reading{{PointID: "skip"}, {PointID: "keep"}}
	if len(telemetry.FilterReadings(v, telemetry.Filter{PointID: "keep"})) != 1 || v[0].PointID != "skip" {
		t.Fatal("alias")
	}
}
