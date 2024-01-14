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

type Handler struct {
	slog.Handler

	format string
	pretty bool
	w      io.Writer
	b      *bytes.Buffer
	m      *sync.Mutex
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	h.m.Lock()

	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	var out []byte

	if h.format == "json" {
		r.Add(slog.String("level", strings.ToLower(r.Level.String())))
		r.Add(slog.String("msg", r.Message))
		r.Add(slog.String("time", r.Time.Format(time.DateTime)))
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
		r.Add(slog.Any(k, v))
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	if h.pretty {
		if r.Message == "raw" {
			for k, v := range fields {
				out = append(out, []byte(fmt.Sprintf("\n\t%s: %#v", k, v))...)
			}
		} else {
			if b, err := json.MarshalIndent(fields, "", "  "); err != nil {
				return err
			} else {
				out = append(out, b...)
			}
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

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	h2 := *h
	h2.Handler = h.Handler.WithAttrs(attrs)

	return &h2
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.Handler = h.Handler.WithGroup(name)

	return &h2
}

func NewHandler(out io.Writer, opts *Options) Handler {
	b := new(bytes.Buffer)

	return Handler{
		Handler: slog.NewJSONHandler(b, opts.HandlerOptions),
		format:  opts.Format,
		pretty:  opts.Pretty,
		b:       b,
		m:       &sync.Mutex{},
		w:       out,
	}
}