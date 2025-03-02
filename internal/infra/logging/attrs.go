package logging

import "log/slog"

// Err converts error to slog.Attr.
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
