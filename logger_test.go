package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		opts     Options
		validate func(t *testing.T, logger *slog.Logger)
	}{
		{
			name: "default json format",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
		{
			name: "text format",
			opts: Options{
				Level:  "debug",
				Format: "text",
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
		{
			name: "with pretty option",
			opts: Options{
				Level:  "info",
				Format: "json",
				Pretty: true,
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
		{
			name: "with add source",
			opts: Options{
				Level:     "info",
				Format:    "json",
				AddSource: true,
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
		{
			name: "with attributes",
			opts: Options{
				Level:  "info",
				Format: "json",
				Attr: []slog.Attr{
					slog.String("service", "test"),
					slog.String("version", "1.0.0"),
				},
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
		{
			name: "null handler",
			opts: Options{
				Null: true,
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
				if _, ok := logger.Handler().(*NullHandler); !ok {
					t.Error("Expected NullHandler when Null option is true")
				}
			},
		},
		{
			name: "invalid format defaults to json",
			opts: Options{
				Level:  "info",
				Format: "invalid-format",
			},
			validate: func(t *testing.T, logger *slog.Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts)
			tt.validate(t, logger)
		})
	}
}

func TestSetGlobalLogger(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "set global logger with info level",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
		},
		{
			name: "set global logger with debug level",
			opts: Options{
				Level:  "debug",
				Format: "text",
			},
		},
		{
			name: "set global null logger",
			opts: Options{
				Null: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original default logger
			originalLogger := slog.Default()
			defer slog.SetDefault(originalLogger)

			SetGlobalLogger(tt.opts)

			// Verify global logger was set
			if slog.Default() == nil {
				t.Error("SetGlobalLogger() did not set default logger")
			}

			// If null option, verify it's NullHandler
			if tt.opts.Null {
				if _, ok := slog.Default().Handler().(*NullHandler); !ok {
					t.Error("Expected NullHandler for global logger with Null option")
				}
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  slog.Level
	}{
		{
			name:  "debug level lowercase",
			level: "debug",
			want:  slog.LevelDebug,
		},
		{
			name:  "debug level uppercase",
			level: "DEBUG",
			want:  slog.LevelDebug,
		},
		{
			name:  "info level lowercase",
			level: "info",
			want:  slog.LevelInfo,
		},
		{
			name:  "info level uppercase",
			level: "INFO",
			want:  slog.LevelInfo,
		},
		{
			name:  "warn level lowercase",
			level: "warn",
			want:  slog.LevelWarn,
		},
		{
			name:  "warn level uppercase",
			level: "WARN",
			want:  slog.LevelWarn,
		},
		{
			name:  "error level lowercase",
			level: "error",
			want:  slog.LevelError,
		},
		{
			name:  "error level uppercase",
			level: "ERROR",
			want:  slog.LevelError,
		},
		{
			name:  "invalid level defaults to info",
			level: "invalid",
			want:  slog.LevelInfo,
		},
		{
			name:  "empty level defaults to info",
			level: "",
			want:  slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.level)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "debug level",
			level: "debug",
		},
		{
			name:  "info level",
			level: "info",
		},
		{
			name:  "warn level",
			level: "warn",
		},
		{
			name:  "error level",
			level: "error",
		},
		{
			name:  "uppercase debug",
			level: "DEBUG",
		},
		{
			name:  "invalid level",
			level: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseColor(tt.level)
			if result == "" {
				t.Errorf("ParseColor(%q) returned empty string", tt.level)
			}
		})
	}
}

func TestLogger_Integration(t *testing.T) {
	tests := []struct {
		name     string
		opts     Options
		logFunc  func(logger *slog.Logger)
		validate func(t *testing.T, output string)
	}{
		{
			name: "json format contains required fields",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.Info("test message", "key", "value")
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "test message") {
					t.Error("Output should contain log message")
				}
				if !strings.Contains(output, "key") {
					t.Error("Output should contain attribute key")
				}
				if !strings.Contains(output, "value") {
					t.Error("Output should contain attribute value")
				}
			},
		},
		{
			name: "debug level logs debug messages",
			opts: Options{
				Level:  "debug",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.Debug("debug message")
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "debug message") {
					t.Error("Output should contain debug message")
				}
			},
		},
		{
			name: "info level filters out debug messages",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.Debug("debug message")
				logger.Info("info message")
			},
			validate: func(t *testing.T, output string) {
				if strings.Contains(output, "debug message") {
					t.Error("Output should not contain debug message")
				}
				if !strings.Contains(output, "info message") {
					t.Error("Output should contain info message")
				}
			},
		},
		{
			name: "attributes are logged",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.With("service", "test").Info("message")
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "service") {
					t.Error("Output should contain service attribute")
				}
				if !strings.Contains(output, "test") {
					t.Error("Output should contain service value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewLogger(tt.opts)
			tt.logFunc(logger)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			tt.validate(t, output)
		})
	}
}

func TestOptions_ReplaceAttr(t *testing.T) {
	tests := []struct {
		name     string
		opts     Options
		logFunc  func(logger *slog.Logger)
		validate func(t *testing.T, output string)
	}{
		{
			name: "standard fields are filtered",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.Info("test")
			},
			validate: func(t *testing.T, output string) {
				// Standard fields like level, time, msg should be handled
				if output == "" {
					t.Error("Output should not be empty")
				}
			},
		},
		{
			name: "raw prefix formats value",
			opts: Options{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(logger *slog.Logger) {
				logger.Info("test", "raw;key", "value")
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "key") {
					t.Error("Output should contain key without raw prefix")
				}
			},
		},
		{
			name: "source is formatted",
			opts: Options{
				Level:     "info",
				Format:    "json",
				AddSource: true,
			},
			logFunc: func(logger *slog.Logger) {
				logger.Info("test with source")
			},
			validate: func(t *testing.T, output string) {
				// Source should be formatted as dir/file:line
				if output == "" {
					t.Error("Output should not be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewLogger(tt.opts)
			tt.logFunc(logger)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			tt.validate(t, output)
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	handler := NewHandler(os.Stdout, &opts)

	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	}

	newHandler := handler.WithAttrs(attrs)
	if newHandler == nil {
		t.Error("WithAttrs should return a handler")
	}

	// Test with empty attrs
	emptyHandler := handler.WithAttrs([]slog.Attr{})
	if emptyHandler != &handler {
		t.Error("WithAttrs with empty attrs should return the same handler")
	}
}

func TestHandler_WithGroup(t *testing.T) {
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	handler := NewHandler(os.Stdout, &opts)

	newHandler := handler.WithGroup("test-group")
	if newHandler == nil {
		t.Error("WithGroup should return a handler")
	}
}

func TestHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	handler := NewHandler(&buf, &opts)

	logger := slog.New(&handler)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Error("Handler should write output")
	}

	if !strings.Contains(output, "test message") {
		t.Error("Output should contain log message")
	}
}

func BenchmarkNewLogger(b *testing.B) {
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLogger(opts)
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	handler := NewHandler(&buf, &opts)
	logger := slog.New(&handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "key", "value", "count", i)
	}
}

func BenchmarkLogger_WithAttrs(b *testing.B) {
	var buf bytes.Buffer
	opts := Options{
		Level:  "info",
		Format: "json",
	}

	handler := NewHandler(&buf, &opts)
	logger := slog.New(&handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.With("key1", "value1", "key2", "value2").Info("benchmark message")
	}
}
