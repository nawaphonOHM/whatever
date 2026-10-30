package core

import (
	"math"
	"time"
)

const (
	maxMillisInt64 = int64(math.MaxInt64 / int64(time.Millisecond))
	maxMillisFloat = float64(maxMillisInt64)
)

func isMsKey(k string) bool {
	return k == "latency_ms" || k == "latencyMs" || k == "latency_millis"
}

func isLatencyKey(k string) bool {
	return k == "latency" || k == "Latency" || isMsKey(k)
}

func formatDurationVal(d time.Duration) (string, bool) {
	if d >= 0 {
		return formatDuration(d), true
	}
	return "", false
}

func formatDurationMillis(n int64) (string, bool) {
	if n >= 0 && n <= maxMillisInt64 {
		return formatDuration(time.Duration(n) * time.Millisecond), true
	}
	return "", false
}

func formatUintMillis(u uint64) (string, bool) {
	if u <= uint64(maxMillisInt64) {
		return formatDurationMillis(int64(u))
	}
	return "", false
}

// formatGCPLatencyAttr formats a latency attribute value into a standard GCP duration string.
// Key-awareness & unit policy:
//   - time.Duration values are formatted directly to standard GCP duration strings (<=9 fractional digits).
//   - Duration strings with units (e.g., "15ms", "0.045s", "15s", "2m3s") are parsed via time.ParseDuration.
//   - Unitless numeric values (int, uint, float) and unitless numeric strings (e.g., 15, 15.5, "15", "15.5")
//     are consistently interpreted as milliseconds matching Gin and standard logging convention.
//   - Non-finite floats (NaN, ±Inf), negative durations/numbers, and out-of-range values are rejected
//     and preserved in the custom payload.
func formatGCPLatencyAttr(key string, v any) (string, bool) {
	if !isLatencyKey(key) {
		return "", false
	}
	return formatMsLatency(v)
}
