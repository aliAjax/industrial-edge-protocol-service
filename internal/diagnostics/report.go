package diagnostics

import (
	"sort"
	"time"
)

type SampleMetric struct {
	Name  string
	Value float64
	Unit  string
	At    time.Time
}
type Report struct {
	GatewayID   string
	GeneratedAt time.Time
	Metrics     []SampleMetric
	Warnings    []string
}

func BuildReport(gateway string, metrics []SampleMetric, now time.Time) Report {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	out := append([]SampleMetric(nil), metrics...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	r := Report{GatewayID: gateway, GeneratedAt: now, Metrics: out}
	for _, m := range out {
		if m.Name == "disk_free_bytes" && m.Value < 100<<20 {
			r.Warnings = append(r.Warnings, "low disk space")
		}
		if m.Name == "cpu_ratio" && m.Value > .9 {
			r.Warnings = append(r.Warnings, "high cpu")
		}
		if m.Name == "memory_ratio" && m.Value > .9 {
			r.Warnings = append(r.Warnings, "high memory")
		}
	}
	return r
}
func (r Report) Healthy() bool { return len(r.Warnings) == 0 }
func (r Report) Metric(name string) (SampleMetric, bool) {
	for _, m := range r.Metrics {
		if m.Name == name {
			return m, true
		}
	}
	return SampleMetric{}, false
}
func (r Report) Age(now time.Time) time.Duration                { return now.Sub(r.GeneratedAt) }
func (r Report) Fresh(now time.Time, maxAge time.Duration) bool { return r.Age(now) <= maxAge }
