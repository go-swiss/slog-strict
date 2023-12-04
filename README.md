# Slog: Strict

`slogstrict` defines a **stricter** interface for [slog](https://pkg.go.dev/log/slog).

## Why?

Slog is a great logging library, but the `*slog.Logger` contains a lot of methods. This makes it easy to use, but also easy to misuse. The `slogstrict.Logger` in this package does this:

1. Always require a context for log messages. This is very useful if you want to include things like `requestID`, `traceID` or `spanID` in our log messages.
2. Always log with `slog.Attr`. This prevents any mistakes when using `...any` in the `slog.Logger` methods.

```go
type Logger interface {
	With(attrs ...slog.Attr) Logger
	WithGroup(name string) Logger

	Debug(ctx context.Context, msg string, attrs ...slog.Attr)
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	Warn(ctx context.Context, msg string, attrs ...slog.Attr)
	Error(ctx context.Context, msg string, err error, attrs ...slog.Attr)

	// For custom levels
	Log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

// The default implementation of [Logger] also implements [ToSlogger]
// so you can always get the underlying [slog.Logger]
type ToSlogger interface {
    // To yeid the underlying *slog.Logger
    ToSlog() *slog.Logger
}
```

The package includes a `slogstrict.New()` function that creates a new logger from a `slog.Handler`, similar to `slog.New()`.

```go
// Create a new [Logger] from a [slog.Handler]
func New(h slog.Handler) Logger {
	return impl{slog.New(h)}
}
```

The package also includes a `slogstrict.FromSlog()` function that creates a new logger from a `*slog.Logger`.

```go
// Create a new [Logger] from a [*slog.Logger]
func FromSlog(l *slog.Logger) Logger {
	return logger{l}
}
```
