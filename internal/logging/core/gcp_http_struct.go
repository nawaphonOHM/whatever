package core

import "log/slog"

const maxHTTPFieldCount = 8

func appendStrAttr(attrs *[]slog.Attr, key, val string) {
	if val != "" {
		*attrs = append(*attrs, slog.String(key, val))
	}
}

func appendStructStrAttrs(attrs *[]slog.Attr, s *HttpRequest) {
	appendStrAttr(attrs, "requestMethod", s.RequestMethod)
	appendStrAttr(attrs, "requestUrl", s.RequestURL)
	appendStrAttr(attrs, "latency", s.Latency)
	appendStrAttr(attrs, "userAgent", s.UserAgent)
	appendStrAttr(attrs, "remoteIp", s.RemoteIP)
	appendStrAttr(attrs, "protocol", s.Protocol)
}

func appendStructNumAttrs(attrs *[]slog.Attr, s *HttpRequest) {
	if s.ResponseSize != 0 {
		*attrs = append(*attrs, slog.Int64("responseSize", s.ResponseSize))
	}
	if s.Status != 0 {
		*attrs = append(*attrs, slog.Int("status", s.Status))
	}
}

func structToAttrs(s *HttpRequest) []slog.Attr {
	attrs := make([]slog.Attr, 0, maxHTTPFieldCount)
	appendStructStrAttrs(&attrs, s)
	appendStructNumAttrs(&attrs, s)
	return attrs
}

func explicitPtrToAttrs(v any) ([]slog.Attr, bool) {
	if s, ok := v.(*HttpRequest); ok && s != nil {
		return structToAttrs(s), true
	}
	return nil, false
}

func explicitStructToAttrs(v any) ([]slog.Attr, bool) {
	if s, ok := v.(HttpRequest); ok {
		return structToAttrs(&s), true
	}
	return explicitPtrToAttrs(v)
}

func explicitValToAttrs(v any) ([]slog.Attr, bool) {
	if attrs, ok := explicitStructToAttrs(v); ok {
		return attrs, true
	}
	if m, ok := v.(map[string]any); ok {
		return mapToAttrs(m), true
	}
	return nil, false
}
