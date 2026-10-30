package core

import "log/slog"

// httpAttrContext holds temporary pointers for URL and latency parsing state.
type httpAttrContext struct {
	pathVal       *string
	queryVal      *string
	latencyFromMs *bool
}

func applyBasicHTTPAttr(req *HttpRequest, a slog.Attr, ctx *httpAttrContext) bool {
	if applyMethodOrPath(req, a, ctx.pathVal) {
		return true
	}
	return applyQueryOrURL(req, a, ctx.queryVal)
}

func applyExtendedHTTPAttr(req *HttpRequest, a slog.Attr, latencyFromMs *bool) bool {
	if applyStatusOrLatency(req, a, latencyFromMs) {
		return true
	}
	return applyClientHTTPAttr(req, a)
}

func applyHTTPAttr(req *HttpRequest, a slog.Attr, ctx *httpAttrContext) bool {
	if applyBasicHTTPAttr(req, a, ctx) {
		return true
	}
	return applyExtendedHTTPAttr(req, a, ctx.latencyFromMs)
}
