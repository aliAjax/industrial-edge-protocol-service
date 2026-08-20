package telemetry

import "time"

type Clock struct{ maxFuture time.Duration }

func NewClock(maxFuture time.Duration) *Clock {
	if maxFuture <= 0 {
		maxFuture = time.Minute
	}
	return &Clock{maxFuture: maxFuture}
}
func (c *Clock) Normalize(observed, received time.Time) time.Time {
	if observed.IsZero() {
		return received
	}
	if observed.After(received.Add(c.maxFuture)) {
		return received
	}
	return observed
}
func (c *Clock) Accept(observed, received time.Time) bool {
	return !observed.After(received.Add(c.maxFuture)) && !observed.Before(received.Add(-24*time.Hour))
}
