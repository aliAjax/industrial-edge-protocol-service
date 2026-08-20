package aggregation

import (
	"industrial-edge-protocol/internal/domain"
	"math"
	"sort"
	"time"
)

type Bucket struct {
	Start time.Time
	End   time.Time
	Count int
	Min   float64
	Max   float64
	Sum   float64
	Good  int
	Bad   int
}

func Build(values []domain.Reading, start, end time.Time, step time.Duration) []Bucket {
	if step <= 0 {
		step = time.Minute
	}
	count := int(end.Sub(start) / step)
	if count < 1 {
		count = 1
	}
	out := make([]Bucket, count)
	for i := range out {
		out[i] = Bucket{Start: start.Add(time.Duration(i) * step), End: start.Add(time.Duration(i+1) * step), Min: math.Inf(1), Max: math.Inf(-1)}
	}
	for _, v := range values {
		if v.ObservedAt.Before(start) || !v.ObservedAt.Before(end) {
			continue
		}
		i := int(v.ObservedAt.Sub(start) / step)
		if i < 0 || i >= len(out) {
			continue
		}
		b := &out[i]
		b.Count++
		b.Sum += v.Value
		if v.Value < b.Min {
			b.Min = v.Value
		}
		if v.Value > b.Max {
			b.Max = v.Value
		}
		if v.Quality == domain.QualityGood {
			b.Good++
		} else {
			b.Bad++
		}
	}
	return out
}
func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	v := append([]float64(nil), values...)
	sort.Float64s(v)
	middle := len(v) / 2
	if len(v)%2 == 1 {
		return v[middle]
	}
	return (v[middle-1] + v[middle]) / 2
}
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	v := append([]float64(nil), values...)
	sort.Float64s(v)
	index := int(float64(len(v)-1) * p)
	return v[index]
}
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
