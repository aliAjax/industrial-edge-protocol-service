package sampling

import (
	"industrial-edge-protocol/internal/domain"
	"sort"
	"time"
)

type Resampler struct{ Interval time.Duration }

func (r Resampler) Apply(values []domain.Reading) []domain.Reading {
	if r.Interval <= 0 {
		return append([]domain.Reading(nil), values...)
	}
	out := []domain.Reading{}
	sorted := append([]domain.Reading(nil), values...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].ObservedAt.Before(sorted[j].ObservedAt) })
	var bucket time.Time
	var sum float64
	var count int
	var template domain.Reading
	flush := func() {
		if count > 0 {
			template.Value = sum / float64(count)
			template.ObservedAt = bucket
			out = append(out, template)
			sum = 0
			count = 0
		}
	}
	for _, v := range sorted {
		current := v.ObservedAt.Truncate(r.Interval)
		if bucket.IsZero() {
			bucket = current
		}
		if current != bucket {
			flush()
			bucket = current
		}
		sum += v.Value
		count++
		template = v
	}
	flush()
	return out
}
func Align(t time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		return t
	}
	return t.Truncate(interval)
}
