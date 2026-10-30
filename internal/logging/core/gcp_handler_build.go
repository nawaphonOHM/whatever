package core

import (
	"log/slog"
	"time"
)

func resolveRecordTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func (h *GCPHandler) enrichSource(payload map[string]any, pc uintptr) {
	if h.opts.AddSource && pc != 0 {
		addSourceLocation(payload, pc)
	}
}

func replaceGroupAttr(groups []string, a slog.Attr, replace ReplaceAttrFunc) []slog.Attr {
	val := a.Value.Resolve()
	if a.Key == "" {
		return replaceAttrsList(groups, val.Group(), replace)
	}
	subGroups := cloneGroups(groups, a.Key)
	subReplaced := replaceAttrsList(subGroups, val.Group(), replace)
	if len(subReplaced) > 0 {
		return []slog.Attr{slog.Group(a.Key, leftoversToAny(subReplaced)...)}
	}
	return nil
}

func appendReplacedAttr(res []slog.Attr, groups []string, a slog.Attr, replace ReplaceAttrFunc) []slog.Attr {
	val := a.Value.Resolve()
	if val.Kind() == slog.KindGroup {
		return append(res, replaceGroupAttr(groups, a, replace)...)
	}
	if attr, ok := applyReplace(groups, a, replace); ok {
		return append(res, attr)
	}
	return res
}

func replaceAttrsList(groups []string, attrs []slog.Attr, replace ReplaceAttrFunc) []slog.Attr {
	res := make([]slog.Attr, 0, len(attrs))
	for _, a := range attrs {
		res = appendReplacedAttr(res, groups, a, replace)
	}
	return res
}

func mergeTopAttrs(stored []slog.Attr, record slog.Record, replace ReplaceAttrFunc) []slog.Attr {
	rawAttrs := collectRecordAttrs(record)
	recordReplaced := replaceAttrsList(nil, rawAttrs, replace)
	allTop := make([]slog.Attr, len(stored)+len(recordReplaced))
	copy(allTop, stored)
	copy(allTop[len(stored):], recordReplaced)
	return allTop
}

func (h *GCPHandler) processUngroupedAttrs(record slog.Record, payload map[string]any) *HttpRequest {
	allTop := mergeTopAttrs(h.topAttrs, record, h.opts.ReplaceAttr)
	isHTTP := isHTTPMessage(record.Message)
	if isHTTP || hasExplicitHTTPAttr(allTop) {
		httpReq, remaining := extractHTTPRequest(allTop, isHTTP)
		addAttrsToMap(payload, nil, remaining, nil)
		return httpReq
	}
	addAttrsToMap(payload, nil, allTop, nil)
	return nil
}

func (h *GCPHandler) processRecordAttrs(record slog.Record, payload map[string]any) *HttpRequest {
	var httpReq *HttpRequest
	if len(h.groups) == 0 {
		httpReq = h.processUngroupedAttrs(record, payload)
	} else {
		rawAttrs := collectRecordAttrs(record)
		addAttrsToMap(payload, h.groups, rawAttrs, h.opts.ReplaceAttr)
	}
	h.enrichSource(payload, record.PC)
	return httpReq
}

func (h *GCPHandler) buildEntry(record slog.Record) GCPLogEntry {
	recordPayload := cloneMap(h.payload)
	httpReq := h.processRecordAttrs(record, recordPayload)
	return GCPLogEntry{
		HTTPRequest: httpReq,
		Payload:     recordPayload,
		Message:     ParseGCPMessage(record.Message),
		Timestamp:   resolveRecordTime(record.Time),
		Severity:    SeverityFromLevel(record.Level),
	}
}
