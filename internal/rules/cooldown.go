package rules

import "time"

type Cooldown struct {
	duration time.Duration
	last     time.Time
}

func NewCooldown(duration time.Duration) *Cooldown { return &Cooldown{duration: duration} }
func (c *Cooldown) Allow(now time.Time) bool {
	if c.last.IsZero() || now.Sub(c.last) >= c.duration {
		c.last = now
		return true
	}
	return false
}
func (c *Cooldown) Remaining(now time.Time) time.Duration {
	if c.last.IsZero() {
		return 0
	}
	remaining := c.duration - now.Sub(c.last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
