package logger

import (
	"context"
	"log/slog"
)

type ComponentLevelHandler struct {
	next         slog.Handler
	defaultLevel slog.Level
	componentMin map[string]slog.Level
	staticAttrs  []slog.Attr
	groups       []string
}

func NewComponentLevelHandler(next slog.Handler, defaultLevel slog.Level, componentMin map[string]slog.Level) *ComponentLevelHandler {
	return &ComponentLevelHandler{
		next:         next,
		defaultLevel: defaultLevel,
		componentMin: componentMin,
	}
}

func (h *ComponentLevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *ComponentLevelHandler) Handle(ctx context.Context, r slog.Record) error {
	component := h.findComponent(r)

	min := h.defaultLevel
	if component != "" {
		if v, ok := h.componentMin[component]; ok {
			min = v
		}
	}

	if r.Level < min {
		return nil
	}
	if len(h.staticAttrs) > 0 {
		r = r.Clone()
		for _, a := range h.staticAttrs {
			r.AddAttrs(a)
		}
	}

	return h.next.Handle(ctx, r)
}

func (h *ComponentLevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nh := *h
	nh.staticAttrs = append(append([]slog.Attr(nil), h.staticAttrs...), attrs...)
	nh.next = h.next.WithAttrs(attrs)
	return &nh
}

func (h *ComponentLevelHandler) WithGroup(name string) slog.Handler {
	nh := *h
	nh.groups = append(append([]string(nil), h.groups...), name)
	nh.next = h.next.WithGroup(name)
	return &nh
}

func (h *ComponentLevelHandler) findComponent(r slog.Record) string {
	if c := attrString(r, "component"); c != "" {
		return c
	}
	if c := attrsString(h.staticAttrs, "component"); c != "" {
		return c
	}
	return ""
}

func attrString(r slog.Record, key string) string {
	var out string
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == key {
			if v, ok := a.Value.Any().(string); ok {
				out = v
				return false
			}
			if a.Value.Kind() == slog.KindString {
				out = a.Value.String()
				return false
			}
		}
		return true
	})
	return out
}

func attrsString(attrs []slog.Attr, key string) string {
	for _, a := range attrs {
		if a.Key == key {
			if v, ok := a.Value.Any().(string); ok {
				return v
			}
			if a.Value.Kind() == slog.KindString {
				return a.Value.String()
			}
		}
	}
	return ""
}
