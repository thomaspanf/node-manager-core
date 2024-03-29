package log

import (
	"log/slog"
	"time"
)

// Prints an error to a log line
func Err(err error) slog.Attr {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return slog.String(ErrorKey, msg)
}

// Replaces the default time formatting (RFC3339) in a logger with an easier to read format
func replaceTime(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Time()
		return slog.String(slog.TimeKey, t.UTC().Format(time.DateTime))
	}
	return a
}
