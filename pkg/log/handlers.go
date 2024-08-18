package log

import (
	"io"
	"log/slog"
	"strings"

	"github.com/lmittmann/tint"
)

// NewCloudFunctionHandler returns a new log handler for Google Cloud Functions.
func NewCloudFunctionHandler(w io.Writer, opt *slog.HandlerOptions) slog.Handler {
	opt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "level" {
			a.Key = "severity"
			a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
		}
		if a.Key == "msg" {
			a.Key = "message"
		}
		return a
	}

	return slog.NewJSONHandler(w, opt)
}

// NewPrettyHandler returns a new log handler for pretty printing.
func NewPrettyHandler(w io.Writer, opt *slog.HandlerOptions) slog.Handler {
	topt := &tint.Options{
		Level:     opt.Level,
		AddSource: opt.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}

	return tint.NewHandler(w, topt)
}

// NewLogfmtHandler returns a new log handler for logfmt printing.
func NewLogfmtHandler(w io.Writer, opt *slog.HandlerOptions) slog.Handler {
	return slog.NewTextHandler(w, opt)
}
