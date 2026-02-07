package logger

import (
	"context"
	"log/slog"
)

type ComponentLevelHandler struct {
	next         slog.Handler
	defaultLevel slog.Level
	levels       map[string]slog.Level
	attrs        []slog.Attr
}

func NewComponentLevelHandler(
	next slog.Handler,
	defaultLevel slog.Level,
	levels map[string]slog.Level,
) *ComponentLevelHandler {
	return &ComponentLevelHandler{
		next:         next,
		defaultLevel: defaultLevel,
		levels:       levels,
	}
}

func (h *ComponentLevelHandler) Enabled(
	ctx context.Context,
	level slog.Level,
) bool {

	// If the level is below default AND we can't
	// yet prove a component wants DEBUG,
	// allow it — we must inspect the record later.
	if level < h.defaultLevel {
		return true
	}

	return h.next.Enabled(ctx, level)
}

func (h *ComponentLevelHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {

	component := h.componentFromRecord(r)

	min := h.defaultLevel
	if lvl, ok := h.levels[component]; ok {
		min = lvl
	}

	// Drop if below allowed level
	if r.Level < min {
		return nil
	}

	return h.next.Handle(ctx, r)
}

func (h *ComponentLevelHandler) componentFromRecord(
	r slog.Record,
) string {

	// Check record attrs first
	var component string

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "component" {
			component = a.Value.String()
			return false
		}
		return true
	})

	if component != "" {
		return component
	}

	// Check static attrs attached via WithAttrs
	for _, a := range h.attrs {
		if a.Key == "component" {
			return a.Value.String()
		}
	}

	return ""
}

func (h *ComponentLevelHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {

	return &ComponentLevelHandler{
		next:         h.next.WithAttrs(attrs),
		defaultLevel: h.defaultLevel,
		levels:       h.levels,
		attrs:        append(h.attrs, attrs...),
	}
}

func (h *ComponentLevelHandler) WithGroup(
	name string,
) slog.Handler {

	return &ComponentLevelHandler{
		next:         h.next.WithGroup(name),
		defaultLevel: h.defaultLevel,
		levels:       h.levels,
		attrs:        h.attrs,
	}
}
