package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Handler is a custom slog.Handler that formats log records with support for JSON and text output.
// It wraps the standard slog.Handler and provides additional formatting capabilities including
// colored text output and pretty-printed JSON.
type Handler struct {
	slog.Handler

	format string        // format specifies output format: "json" or "text"
	pretty bool          // pretty enables JSON indentation
	w      io.Writer     // w is the output destination
	b      *bytes.Buffer // b is an internal buffer for processing log records
	m      *sync.Mutex   // m protects concurrent access to the buffer
}

// Handle processes a log record and writes it to the output writer.
// For JSON format, it creates a structured record with level, message, time, and attributes.
// For text format, it creates a human-readable colored output.
// This method is thread-safe and handles concurrent logging calls.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	h.m.Lock()

	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	var (
		fields = make(map[string]interface{}, r.NumAttrs())
		out    []byte
	)

	if h.format == "json" {
		fields["level"] = strings.ToLower(r.Level.String())
		fields["msg"] = r.Message
		fields["time"] = r.Time.Format(time.DateTime)
	} else {
		out = []byte(fmt.Sprintf("%s %s %s ",
			r.Time.Format(time.DateTime),
			ParseColor(r.Level.String()),
			color.CyanString(r.Message),
		))
	}

	if err := h.Handler.Handle(ctx, r); err != nil {
		return err
	}

	attrs := map[string]any{}
	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return err
	}

	for k, v := range attrs {
		fields[k] = v
	}

	if h.pretty {
		if b, err := json.MarshalIndent(fields, "", "  "); err != nil {
			return err
		} else {
			out = append(out, b...)
		}
	} else {
		if b, err := json.Marshal(fields); err != nil {
			return err
		} else {
			out = append(out, b...)
		}
	}

	h.w.Write(append(out, "\n"...))

	return nil
}

// WithAttrs returns a new Handler with the specified attributes added to all log records.
// If no attributes are provided, returns the same handler.
// This method creates a shallow copy of the handler with updated attributes.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	h2 := *h
	h2.Handler = h.Handler.WithAttrs(attrs)

	return &h2
}

// WithGroup returns a new Handler with the specified group name applied to all subsequent attributes.
// Groups allow hierarchical organization of log attributes in the output.
// This method creates a shallow copy of the handler with the group applied.
func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.Handler = h.Handler.WithGroup(name)

	return &h2
}

// NewHandler creates and initializes a new Handler with the specified output writer and options.
// If the format option is not "json" or "text", it defaults to "json".
// The handler uses an internal JSON handler for processing attributes and a buffer for intermediate storage.
func NewHandler(out io.Writer, opts *Options) Handler {
	b := new(bytes.Buffer)

	if !map[string]bool{
		"json": true,
		"text": true,
	}[opts.Format] {
		opts.Format = "json"
	}

	return Handler{
		Handler: slog.NewJSONHandler(b, opts.HandlerOptions),
		format:  opts.Format,
		pretty:  opts.Pretty,
		b:       b,
		m:       &sync.Mutex{},
		w:       out,
	}
}
