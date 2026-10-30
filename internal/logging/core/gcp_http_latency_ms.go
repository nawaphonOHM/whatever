package core

import (
	"math"
	"strconv"
	"strings"
	"time"
)

func formatMsSigned(v any) (string, bool) {
	if n, ok := v.(int); ok {
		return formatDurationMillis(int64(n))
	}
	if n, ok := v.(int64); ok {
		return formatDurationMillis(n)
	}
	return "", false
}

func formatMsUint(v any) (string, bool) {
	if u, ok := v.(uint); ok {
		return formatUintMillis(uint64(u))
	}
	if u, ok := v.(uint64); ok {
		return formatUintMillis(u)
	}
	return "", false
}

func formatFloatMillis(ms float64) (string, bool) {
	if !isFiniteNonNegativeFloat(ms, maxMillisFloat) {
		return "", false
	}
	ns := math.Round(ms * nanosPerMilli)
	return formatDuration(time.Duration(ns)), true
}

func formatMsFloat(v any) (string, bool) {
	if f, ok := v.(float64); ok {
		return formatFloatMillis(f)
	}
	if f, ok := v.(float32); ok {
		return formatFloatMillis(float64(f))
	}
	return "", false
}

func parseMsFloatString(s string) (string, bool) {
	f, err := strconv.ParseFloat(s, floatBitSize)
	if err == nil {
		return formatFloatMillis(f)
	}
	return "", false
}

func parseMsStringLatency(s string) (string, bool) {
	trimmed := strings.TrimSpace(s)
	if res, ok := parseDurationString(trimmed); ok {
		return res, true
	}
	return parseMsFloatString(trimmed)
}

func formatMsDuration(v any) (string, bool) {
	if d, ok := v.(time.Duration); ok {
		return formatDurationVal(d)
	}
	if s, ok := v.(string); ok {
		return parseMsStringLatency(s)
	}
	return "", false
}

func formatMsNumeric(v any) (string, bool) {
	if res, ok := formatMsSigned(v); ok {
		return res, true
	}
	if res, ok := formatMsUint(v); ok {
		return res, true
	}
	return formatMsFloat(v)
}

func formatMsLatency(v any) (string, bool) {
	if res, ok := formatMsDuration(v); ok {
		return res, true
	}
	return formatMsNumeric(v)
}
