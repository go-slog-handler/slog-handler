package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Options struct {
	*slog.HandlerOptions

	AddSource bool
	Attr      []slog.Attr
	Format    string
	Level     string
	Pretty    bool
}

func NewLogger(opts Options) *slog.Logger {
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

func SetGlobalLogger(opts Options) {
	logger := NewLogger(opts)

	slog.SetDefault(logger)
}

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
