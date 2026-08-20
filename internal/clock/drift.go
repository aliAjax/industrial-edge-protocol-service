package clock

import "time"

type Sample struct {
	Device  time.Time
	Gateway time.Time
}
type Drift struct {
	Offset time.Duration
	Skew   float64
	Valid  bool
}

func Estimate(samples []Sample) Drift {
	if len(samples) == 0 {
		return Drift{}
	}
	var total time.Duration
	for _, s := range samples {
		total += s.Gateway.Sub(s.Device)
	}
	offset := total / time.Duration(len(samples))
	return Drift{Offset: offset, Skew: float64(offset) / float64(time.Second), Valid: offset > -5*time.Minute && offset < 5*time.Minute}
}
func Correct(at time.Time, d Drift) time.Time {
	if !d.Valid {
		return at
	}
	return at.Add(d.Offset)
}
