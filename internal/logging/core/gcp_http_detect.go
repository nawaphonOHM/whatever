package core

import (
	"log/slog"
	"strings"
)

func isHTTPMessage(msg string) bool {
	return strings.EqualFold(msg, "HTTP Request") ||
		strings.EqualFold(msg, "http_request") ||
		strings.EqualFold(msg, "gin")
}

func isExplicitHTTPKey(k string) bool {
	return k == "httpRequest" || k == "http_request"
}

func hasExplicitHTTPAttr(attrs []slog.Attr) bool {
	for _, a := range attrs {
		if isExplicitHTTPKey(a.Key) {
			return true
		}
	}
	return false
}
