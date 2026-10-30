package core

import "log/slog"

func isUserAgentKey(k string) bool {
	return k == "user_agent" || k == "userAgent"
}

func isRemoteIPKey(k string) bool {
	switch k {
	case "client_ip", "remoteIp", "remote_ip", "clientIp":
		return true
	}
	return false
}

func resolveAttrString(a slog.Attr) (string, bool) {
	resolved := a.Value.Resolve()
	if resolved.Kind() != slog.KindString {
		return "", false
	}
	return resolved.String(), true
}

// setStrField stores the attr string value into dst when dst is unset and consumes empty strings.
func setStrField(dst *string, a slog.Attr) bool {
	val, ok := resolveAttrString(a)
	if !ok {
		return false
	}
	if *dst == "" {
		*dst = val
		return true
	}
	return val == ""
}

func applyUserAgent(req *HttpRequest, a slog.Attr) bool {
	if !isUserAgentKey(a.Key) {
		return false
	}
	return setStrField(&req.UserAgent, a)
}

func applyRemoteIP(req *HttpRequest, a slog.Attr) bool {
	if !isRemoteIPKey(a.Key) {
		return false
	}
	return setStrField(&req.RemoteIP, a)
}

func applyUserAgentOrIP(req *HttpRequest, a slog.Attr) bool {
	if applyUserAgent(req, a) {
		return true
	}
	return applyRemoteIP(req, a)
}

func applyProtocol(req *HttpRequest, a slog.Attr) bool {
	if a.Key != "protocol" {
		return false
	}
	return setStrField(&req.Protocol, a)
}

func applySizeOrProtocol(req *HttpRequest, a slog.Attr) bool {
	if applyResponseSize(req, a) {
		return true
	}
	return applyProtocol(req, a)
}

func applyClientHTTPAttr(req *HttpRequest, a slog.Attr) bool {
	if applyUserAgentOrIP(req, a) {
		return true
	}
	return applySizeOrProtocol(req, a)
}
