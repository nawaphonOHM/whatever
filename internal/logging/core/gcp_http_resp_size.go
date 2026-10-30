package core

import "log/slog"

func setResponseSize(req *HttpRequest, a slog.Attr) bool {
	if req.ResponseSize != 0 {
		return false
	}
	sz, ok := parseResponseSize(a.Value.Resolve().Any())
	if !ok {
		return false
	}
	req.ResponseSize = sz
	return true
}

func applyResponseSize(req *HttpRequest, a slog.Attr) bool {
	if !isResponseSizeKey(a.Key) {
		return false
	}
	return setResponseSize(req, a)
}
