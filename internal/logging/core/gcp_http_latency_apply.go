package core

import "log/slog"

func applyMsLatency(req *HttpRequest, lat string, latencyFromMs *bool) bool {
	if req.Latency == "" {
		req.Latency = lat
		if latencyFromMs != nil {
			*latencyFromMs = true
		}
	}
	return true
}

func tryUpgradeLatency(req *HttpRequest, lat string, latencyFromMs *bool) bool {
	if latencyFromMs == nil || !*latencyFromMs {
		return false
	}
	req.Latency = lat
	*latencyFromMs = false
	return true
}

func applyDurationLatency(req *HttpRequest, lat string, latencyFromMs *bool) bool {
	if req.Latency == "" {
		req.Latency = lat
		return true
	}
	return tryUpgradeLatency(req, lat, latencyFromMs)
}

func setLatency(req *HttpRequest, a slog.Attr, latencyFromMs *bool) bool {
	lat, ok := formatGCPLatencyAttr(a.Key, a.Value.Resolve().Any())
	if !ok {
		return false
	}
	if isMsKey(a.Key) {
		return applyMsLatency(req, lat, latencyFromMs)
	}
	return applyDurationLatency(req, lat, latencyFromMs)
}

func applyLatencyAttr(req *HttpRequest, a slog.Attr, latencyFromMs *bool) bool {
	if !isLatencyKey(a.Key) {
		return false
	}
	return setLatency(req, a, latencyFromMs)
}
