package core

import (
	"log/slog"
	"slices"
)

func mapToAttrs(m map[string]any) []slog.Attr {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	attrs := make([]slog.Attr, 0, len(keys))
	for _, k := range keys {
		attrs = append(attrs, slog.Any(k, m[k]))
	}
	return attrs
}

func appendLeftoverStr(remaining *[]slog.Attr, key, val string) {
	if val != "" {
		*remaining = append(*remaining, slog.String(key, val))
	}
}

func finalizeExplicitURL(req *HttpRequest, pathVal, queryVal string, leftovers *[]slog.Attr) {
	if pathVal != "" && req.RequestURL == "" {
		req.RequestURL = composeURL(pathVal, queryVal)
		return
	}
	appendLeftoverStr(leftovers, "path", pathVal)
	appendLeftoverStr(leftovers, "query", queryVal)
}

func processExplicitSubAttrs(req *HttpRequest, attrs []slog.Attr) []slog.Attr {
	var explicitPath, explicitQuery string
	var explicitMs bool
	ctx := httpAttrContext{pathVal: &explicitPath, queryVal: &explicitQuery, latencyFromMs: &explicitMs}
	leftovers := make([]slog.Attr, 0, len(attrs))
	for _, sub := range attrs {
		if !applyHTTPAttr(req, sub, &ctx) {
			leftovers = append(leftovers, sub)
		}
	}
	finalizeExplicitURL(req, explicitPath, explicitQuery, &leftovers)
	return leftovers
}

func leftoversToAny(attrs []slog.Attr) []any {
	res := make([]any, len(attrs))
	for i, a := range attrs {
		res[i] = a
	}
	return res
}

func resolveExplicitAttrs(val slog.Value) ([]slog.Attr, bool) {
	if val.Kind() == slog.KindGroup {
		return val.Group(), true
	}
	return explicitValToAttrs(val.Any())
}

func parseExplicitHTTPAttr(req *HttpRequest, a slog.Attr) (bool, *slog.Attr) {
	subAttrs, ok := resolveExplicitAttrs(a.Value.Resolve())
	if !ok {
		return false, nil
	}
	leftovers := processExplicitSubAttrs(req, subAttrs)
	if len(leftovers) == 0 {
		return true, nil
	}
	leftoverGroup := slog.Group(a.Key, leftoversToAny(leftovers)...)
	return true, &leftoverGroup
}
