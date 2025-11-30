package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNullHandler_Enabled(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  bool
	}{
		{
			name:  "debug level",
			level: slog.LevelDebug,
			want:  false,
		},
		{
			name:  "info level",
			level: slog.LevelInfo,
			want:  false,
		},
		{
			name:  "warn level",
			level: slog.LevelWarn,
			want:  false,
		},
		{
			name:  "error level",
			level: slog.LevelError,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewNullHandler()
			got := h.Enabled(context.Background(), tt.level)
			if got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullHandler_Handle(t *testing.T) {
	h := NewNullHandler()
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	err := h.Handle(context.Background(), r)
	if err != nil {
		t.Errorf("Handle() error = %v, want nil", err)
	}
}

func TestNullHandler_WithAttrs(t *testing.T) {
	h := NewNullHandler()
	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	}

	h2 := h.WithAttrs(attrs)
	if h2 != h {
		t.Error("WithAttrs should return the same handler")
	}
}

func TestNullHandler_WithGroup(t *testing.T) {
	h := NewNullHandler()

	h2 := h.WithGroup("test-group")
	if h2 != h {
		t.Error("WithGroup should return the same handler")
	}
}

func TestNullHandler_Integration(t *testing.T) {
	// Create logger with NullHandler
	logger := slog.New(NewNullHandler())

	// Capture output to ensure nothing is written
	var buf bytes.Buffer

	// Log various messages
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// With attributes
	logger.With("key", "value").Info("message with attrs")

	// With group
	logger.WithGroup("group").Info("message in group")

	// Buffer should remain empty since NullHandler discards everything
	if buf.Len() > 0 {
		t.Error("NullHandler should not write any output")
	}
}

func TestNewLogger_WithNullOption(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantNil bool
	}{
		{
			name: "null handler enabled",
			opts: Options{
				Null: true,
			},
			wantNil: false,
		},
		{
			name: "null handler disabled",
			opts: Options{
				Null:   false,
				Level:  "info",
				Format: "json",
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts)
			if logger == nil {
				t.Error("NewLogger() returned nil")
			}

			// Verify that with Null option, we get NullHandler
			if tt.opts.Null {
				if _, ok := logger.Handler().(*NullHandler); !ok {
					t.Error("Expected NullHandler when Null option is true")
				}
			}
		})
	}
}

func BenchmarkNullHandler_Handle(b *testing.B) {
	h := NewNullHandler()
	ctx := context.Background()
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "benchmark message", 0)
	r.AddAttrs(
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
		slog.Bool("key3", true),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.Handle(ctx, r)
	}
}
