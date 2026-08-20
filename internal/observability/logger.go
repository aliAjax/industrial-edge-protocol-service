package observability

import (
	"context"
	"log/slog"
	"time"
)

type Logger struct{ base *slog.Logger }

func NewLogger(base *slog.Logger) *Logger { return &Logger{base: base} }
func (l *Logger) WithRequest(ctx context.Context, requestID string) *slog.Logger {
	return l.base.With("request_id", requestID, "trace_id", TraceID(ctx))
}
func TraceID(ctx context.Context) string {
	if value, ok := ctx.Value(traceKey{}).(string); ok {
		return value
	}
	return ""
}
func WithTrace(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceKey{}, id)
}

type traceKey struct{}

func (l *Logger) Observe(name string, started time.Time, attrs ...any) {
	args := []any{"duration_ms", time.Since(started).Milliseconds()}
	args = append(args, attrs...)
	l.base.Info(name, args...)
}
