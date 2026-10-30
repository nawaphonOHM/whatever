package core

import (
	"math"
	"strconv"
	"strings"
)

const (
	decimalBase = 10
	int64Bits   = 64
)

func parseSignedSize(v any) (int64, bool) {
	if n, ok := v.(int); ok {
		return int64(n), n >= 0
	}
	if n, ok := v.(int64); ok {
		return n, n >= 0
	}
	return 0, false
}

func parseUintVal(u uint64) (int64, bool) {
	if u <= uint64(math.MaxInt64) {
		return int64(u), true
	}
	return 0, false
}

func parseUintSize(v any) (int64, bool) {
	if u, ok := v.(uint); ok {
		return parseUintVal(uint64(u))
	}
	if u, ok := v.(uint64); ok {
		return parseUintVal(u)
	}
	return 0, false
}

func float64ToSize(f float64) (int64, bool) {
	if f < 0 || f >= float64(math.MaxInt64) {
		return 0, false
	}
	n := int64(f)
	return n, float64(n) == f
}

func parseFloatSize(v any) (int64, bool) {
	if f, ok := v.(float64); ok {
		return float64ToSize(f)
	}
	if f, ok := v.(float32); ok {
		return float64ToSize(float64(f))
	}
	return 0, false
}

func parseUnsignedOrFloatSize(v any) (int64, bool) {
	if n, ok := parseUintSize(v); ok {
		return n, true
	}
	return parseFloatSize(v)
}

func parseNumericSize(v any) (int64, bool) {
	if n, ok := parseSignedSize(v); ok {
		return n, true
	}
	return parseUnsignedOrFloatSize(v)
}

func parseStringSize(s string) (int64, bool) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), decimalBase, int64Bits)
	if err == nil && n >= 0 {
		return n, true
	}
	return 0, false
}

func isResponseSizeKey(k string) bool {
	switch k {
	case "bytes_out", "responseSize", "response_size", "bytesOut":
		return true
	}
	return false
}

func parseResponseSize(v any) (int64, bool) {
	if n, ok := parseNumericSize(v); ok {
		return n, true
	}
	if s, ok := v.(string); ok {
		return parseStringSize(s)
	}
	return 0, false
}
