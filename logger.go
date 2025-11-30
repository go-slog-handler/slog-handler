package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Options configures the logger behavior including format, level, and output options.
// It extends slog.HandlerOptions with additional fields for customization.
type Options struct {
	*slog.HandlerOptions

	AddSource bool        // AddSource includes source file and line number in log output
	Attr      []slog.Attr // Attr is a list of attributes to add to every log record
	Format    string      // Format specifies output format: "json" or "text"
	Level     string      // Level sets minimum log level: "debug", "info", "warn", or "error"
	Pretty    bool        // Pretty enables JSON pretty-printing with indentation
	Null      bool        // Null uses NullHandler to discard all logs (useful for testing)
}

// NewLogger creates a new slog.Logger with the specified options.
// If Null option is true, returns a logger with NullHandler that discards all output.
// Otherwise, creates a custom handler with the configured format, level, and attributes.
func NewLogger(opts Options) *slog.Logger {
	// If Null option is set, return a logger with NullHandler
	if opts.Null {
		return slog.New(NewNullHandler())
	}

	opts.HandlerOptions = &slog.HandlerOptions{
		AddSource: opts.AddSource,
		Level:     ParseLevel(opts.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// skip standart fields
			if map[string]bool{
				slog.LevelKey:   true,
				slog.MessageKey: true,
				slog.TimeKey:    true,
			}[a.Key] {
				return slog.Attr{}
			}

			key := strings.Split(a.Key, ";")

			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)

				dir, file := filepath.Split(s.File)

				a.Value = slog.StringValue(fmt.Sprintf("%s:%d",
					filepath.Join(filepath.Base(dir), file),
					s.Line,
				))
			} else if key[0] == "raw" {
				a.Key = strings.Join(key[1:], ";")
				a.Value = slog.StringValue(fmt.Sprintf("%#v", a.Value.Any()))
			}

			return a
		},
	}

	handler := NewHandler(os.Stdout, &opts)

	return slog.New(handler.WithAttrs(opts.Attr))
}

// SetGlobalLogger creates a new logger with the specified options and sets it as the default global logger.
// This affects all subsequent calls to slog.Info(), slog.Debug(), etc. throughout the application.
func SetGlobalLogger(opts Options) {
	logger := NewLogger(opts)

	slog.SetDefault(logger)
}

// ParseLevel converts a string representation of log level to slog.Level.
// Valid inputs (case-insensitive): "debug", "info", "warn", "error".
// Returns slog.LevelInfo for any unrecognized input as a safe default.
func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// ParseColor returns a colorized string representation of the log level.
// Colors are applied using fatih/color package: white (debug), green (info),
// yellow (warn), red (error). Input is case-insensitive.
func ParseColor(level string) string {
	switch strings.ToLower(level) {
	case "debug":
		return color.WhiteString(level)
	case "error":
		return color.RedString(level)
	case "warn":
		return color.YellowString(level)
	default:
		return color.GreenString(level)
	}
}
