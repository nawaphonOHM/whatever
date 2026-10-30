package core

import "log/slog"

func isMethodKey(k string) bool {
	switch k {
	case "method", "requestMethod", "request_method", "http_method":
		return true
	}
	return false
}

func isPathOrURLKey(k string) bool {
	switch k {
	case "requestUrl", "request_url", "url", "uri":
		return true
	}
	return false
}

func applyMethod(req *HttpRequest, a slog.Attr) bool {
	if !isMethodKey(a.Key) {
		return false
	}
	return setStrField(&req.RequestMethod, a)
}

func applyPath(req *HttpRequest, a slog.Attr, pathVal *string) bool {
	if a.Key != "path" || req.RequestURL != "" {
		return false
	}
	return setStrField(pathVal, a)
}

func applyMethodOrPath(req *HttpRequest, a slog.Attr, pathVal *string) bool {
	if applyMethod(req, a) {
		return true
	}
	return applyPath(req, a, pathVal)
}

func applyQuery(req *HttpRequest, a slog.Attr, queryVal *string) bool {
	if a.Key != "query" || req.RequestURL != "" {
		return false
	}
	return setStrField(queryVal, a)
}

func applyURL(req *HttpRequest, a slog.Attr) bool {
	if !isPathOrURLKey(a.Key) {
		return false
	}
	return setStrField(&req.RequestURL, a)
}

func applyQueryOrURL(req *HttpRequest, a slog.Attr, queryVal *string) bool {
	if applyQuery(req, a, queryVal) {
		return true
	}
	return applyURL(req, a)
}
