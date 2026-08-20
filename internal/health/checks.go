package health

import (
	"context"
	"time"
)

type Check struct {
	Name    string
	Fn      func(context.Context) error
	Timeout time.Duration
}
type Result struct {
	Name    string
	Healthy bool
	Latency time.Duration
	Error   string
}

func Run(ctx context.Context, checks []Check) []Result {
	out := make([]Result, 0, len(checks))
	for _, check := range checks {
		timeout := check.Timeout
		if timeout <= 0 {
			timeout = 2 * time.Second
		}
		child, cancel := context.WithTimeout(ctx, timeout)
		started := time.Now()
		err := check.Fn(child)
		cancel()
		r := Result{Name: check.Name, Healthy: err == nil, Latency: time.Since(started)}
		if err != nil {
			r.Error = err.Error()
		}
		out = append(out, r)
	}
	return out
}
func Healthy(results []Result) bool {
	if len(results) == 0 {
		return true
	}
	for _, r := range results {
		if !r.Healthy {
			return false
		}
	}
	return true
}
