package rules

import "time"

type Debouncer struct {
	delay  time.Duration
	active bool
	since  time.Time
}

func NewDebouncer(delay time.Duration) *Debouncer { return &Debouncer{delay: delay} }
func (d *Debouncer) Update(condition bool, now time.Time) bool {
	if !condition {
		d.active = false
		d.since = time.Time{}
		return false
	}
	if !d.active {
		d.active = true
		d.since = now
	}
	return now.Sub(d.since) >= d.delay
}

type Hysteresis struct {
	High   float64
	Low    float64
	active bool
}

func (h *Hysteresis) Update(value float64) bool {
	if h.active {
		if value <= h.Low {
			h.active = false
		}
	} else if value >= h.High {
		h.active = true
	}
	return h.active
}
