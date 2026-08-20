package reconnect

import "time"

type Policy struct {
	Initial time.Duration
	Max     time.Duration
	Factor  float64
	Jitter  time.Duration
}

func (p Policy) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if p.Initial <= 0 {
		p.Initial = time.Second
	}
	if p.Max <= 0 {
		p.Max = time.Minute
	}
	d := float64(p.Initial)
	for i := 0; i < attempt; i++ {
		d *= p.Factor
		if d >= float64(p.Max) {
			return p.Max
		}
	}
	if time.Duration(d) > p.Max {
		return p.Max
	}
	return time.Duration(d)
}
func (p Policy) Reset() int { return 0 }
