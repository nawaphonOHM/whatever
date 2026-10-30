package core

import "log/slog"

func flushTopAttrs(payload map[string]any, topAttrs []slog.Attr, groups []string) {
	if len(topAttrs) == 0 || len(groups) > 0 {
		return
	}
	addAttrsToMap(payload, nil, topAttrs, nil)
}

// WithGroup returns a new GCPHandler with the specified group added.
func (h *GCPHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newPayload := cloneMap(h.payload)
	flushTopAttrs(newPayload, h.topAttrs, h.groups)
	return &GCPHandler{
		out:     h.out,
		mu:      h.mu,
		opts:    h.opts,
		groups:  cloneGroups(h.groups, name),
		payload: newPayload,
	}
}

func (h *GCPHandler) withUngroupedAttrs(attrs []slog.Attr) slog.Handler {
	replaced := replaceAttrsList(nil, attrs, h.opts.ReplaceAttr)
	newTop := make([]slog.Attr, len(h.topAttrs)+len(replaced))
	copy(newTop, h.topAttrs)
	copy(newTop[len(h.topAttrs):], replaced)
	return &GCPHandler{
		out:      h.out,
		mu:       h.mu,
		opts:     h.opts,
		groups:   h.groups,
		payload:  cloneMap(h.payload),
		topAttrs: newTop,
	}
}

// WithAttrs returns a new GCPHandler with the specified attributes added.
func (h *GCPHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	if len(h.groups) == 0 {
		return h.withUngroupedAttrs(attrs)
	}
	newPayload := cloneMap(h.payload)
	addAttrsToMap(newPayload, h.groups, attrs, h.opts.ReplaceAttr)
	return &GCPHandler{
		out:      h.out,
		mu:       h.mu,
		opts:     h.opts,
		groups:   h.groups,
		payload:  newPayload,
		topAttrs: h.topAttrs,
	}
}
