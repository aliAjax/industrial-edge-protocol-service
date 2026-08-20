package telemetry

import (
	"industrial-edge-protocol/internal/domain"
	"sort"
	"time"
)

type Filter struct {
	PointID domain.ID
	From    time.Time
	To      time.Time
	Quality domain.Quality
	Min     *float64
	Max     *float64
}

func FilterReadings(values []domain.Reading, f Filter) []domain.Reading {
	out := []domain.Reading{}
	for _, v := range values {
		if f.PointID != "" && v.PointID != f.PointID {
			continue
		}
		if !f.From.IsZero() && v.ObservedAt.Before(f.From) {
			continue
		}
		if !f.To.IsZero() && !v.ObservedAt.Before(f.To) {
			continue
		}
		if f.Quality != "" && v.Quality != f.Quality {
			continue
		}
		if f.Min != nil && v.Value < *f.Min {
			continue
		}
		if f.Max != nil && v.Value > *f.Max {
			continue
		}
		out = append(out, v)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ObservedAt.Before(out[j].ObservedAt) })
	return out
}
func Latest(values []domain.Reading) domain.Reading {
	if len(values) == 0 {
		return domain.Reading{}
	}
	latest := values[0]
	for _, v := range values[1:] {
		if v.ObservedAt.After(latest.ObservedAt) {
			latest = v
		}
	}
	return latest
}
