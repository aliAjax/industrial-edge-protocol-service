package commandqueue

import "time"

type Retry struct {
	Attempts int
	Max      int
	Next     time.Time
	Backoff  time.Duration
}

func (r *Retry) Failure(now time.Time) bool {
	r.Attempts++
	if r.Attempts > r.Max {
		return false
	}
	if r.Backoff <= 0 {
		r.Backoff = time.Second
	}
	r.Next = now.Add(time.Duration(r.Attempts) * r.Backoff)
	return true
}
func (r Retry) Ready(now time.Time) bool { return !r.Next.After(now) }
func (r *Retry) Reset() {
	r.Attempts = 0
	r.Next = time.Time{}
	r.Backoff = 0
}
