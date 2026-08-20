package retention

import "time"

type Policy struct {
	Raw       time.Duration
	Aggregate time.Duration
	Audit     time.Duration
}

func (p Policy) Expired(at, now time.Time, kind string) bool {
	age := p.Raw
	switch kind {
	case "aggregate":
		age = p.Aggregate
	case "audit":
		age = p.Audit
	}
	return age > 0 && now.Sub(at) > age
}
func (p Policy) NextArchive(now time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		interval = time.Hour
	}
	return now.Truncate(interval).Add(interval)
}
