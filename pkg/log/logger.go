package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const (
	TypeLogfmt        = "logfmt"
	TypeJSON          = "json"
	TypePretty        = "pretty"
	TypeCloudFunction = "cloudfunction"
)

type ctxKey struct{}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return New(os.Stderr, TypePretty, false)
}

func New(w io.Writer, format string, verbose bool) *slog.Logger {
	opt := &slog.HandlerOptions{Level: slog.LevelInfo}
	if verbose {
		opt.Level = slog.LevelDebug
		opt.AddSource = true
	}

	switch format {
	case TypeLogfmt:
		return slog.New(NewLogfmtHandler(w, opt))
	case TypeJSON:
		return slog.New(slog.NewJSONHandler(w, opt))
	case TypePretty:
		return slog.New(NewPrettyHandler(w, opt))
	case TypeCloudFunction:
		return slog.New(NewCloudFunctionHandler(w, opt))
	default:
		return slog.New(NewPrettyHandler(w, opt))
	}
}
