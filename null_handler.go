package logger

import (
	"context"
	"log/slog"
)

// NullHandler is a slog.Handler that discards all log records.
// It implements the slog.Handler interface but performs no operations,
// making it useful for testing or disabling logging entirely.
type NullHandler struct{}

// NewNullHandler creates a new NullHandler instance.
func NewNullHandler() *NullHandler {
	return &NullHandler{}
}

// Enabled always returns false, indicating that no log levels are enabled.
// This prevents any log records from being processed.
func (h *NullHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle discards the log record without processing it.
// Always returns nil to indicate successful (no-op) handling.
func (h *NullHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same NullHandler, ignoring any attributes.
// This maintains the no-op behavior even when attributes are added.
func (h *NullHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same NullHandler, ignoring any group name.
// This maintains the no-op behavior even when groups are specified.
func (h *NullHandler) WithGroup(_ string) slog.Handler {
	return h
}
