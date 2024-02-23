package slogstrict

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// A limited set of method from slog.Logger
// This forces adding context.Context and using slog.Attr
type Logger interface {
	With(attrs ...slog.Attr) Logger
	WithGroup(name string) Logger

	Debug(ctx context.Context, msg string, attrs ...slog.Attr)
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	Warn(ctx context.Context, msg string, attrs ...slog.Attr)
	Error(ctx context.Context, msg string, err error, attrs ...slog.Attr)

	// For custom levels
	Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)

	// To retrieve a *slog.Logger
	ToSlog() *slog.Logger
}

// Create a new [Logger] from a [slog.Handler]
func New(h slog.Handler) Logger {
	return logger{slog.New(h)}
}

// Create a new [Logger] from a [*slog.Logger]
func FromSlog(l *slog.Logger) Logger {
	return logger{l}
}

type logger struct{ *slog.Logger }

func (s logger) With(attrs ...slog.Attr) Logger {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return logger{s.Logger.With(args...)}
}

func (s logger) WithGroup(name string) Logger {
	return logger{s.Logger.WithGroup(name)}
}

func (s logger) Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, level, msg, attrs...)
}

func (s logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

func (s logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func (s logger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	s.logAttrs(ctx, slog.LevelWarn, msg, attrs...)
}

func (s logger) Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append([]slog.Attr{slog.String("err", err.Error())}, attrs...)
	}
	s.logAttrs(ctx, slog.LevelError, msg, attrs...)
}

func (l logger) logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}

	_ = l.Handler().Handle(ctx, r)
}

func (s logger) ToSlog() *slog.Logger {
	return s.Logger
}
