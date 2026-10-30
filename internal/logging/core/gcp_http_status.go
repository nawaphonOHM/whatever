package core

import (
	"strconv"
	"strings"
)

const (
	minHTTPStatus = 100
	maxHTTPStatus = 599
)

func isValidStatusRange(n int) bool {
	return n >= minHTTPStatus && n <= maxHTTPStatus
}

func parseSignedStatus(v any) (int, bool) {
	if n, ok := v.(int); ok {
		return n, isValidStatusRange(n)
	}
	if n, ok := v.(int64); ok {
		return int(n), isValidStatusRange(int(n))
	}
	return 0, false
}

func parseUintStatus(v any) (int, bool) {
	if u, ok := v.(uint); ok {
		return int(u), isValidStatusRange(int(u))
	}
	if u, ok := v.(uint64); ok {
		return int(u), isValidStatusRange(int(u))
	}
	return 0, false
}

func isStatusInRange(f float64) bool {
	return f >= minHTTPStatus && f <= maxHTTPStatus
}

func parseFloatStatus(v any) (int, bool) {
	if f, ok := v.(float64); ok && isStatusInRange(f) {
		n := int(f)
		return n, float64(n) == f
	}
	return 0, false
}

func parseNumericStatus(v any) (int, bool) {
	if n, ok := parseSignedStatus(v); ok {
		return n, true
	}
	if n, ok := parseUintStatus(v); ok {
		return n, true
	}
	return parseFloatStatus(v)
}

func parseStringStatus(s string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err == nil && isValidStatusRange(n) {
		return n, true
	}
	return 0, false
}

func isStatusKey(k string) bool {
	return k == "status" || k == "status_code" || k == "statusCode"
}

func parseHTTPStatus(v any) (int, bool) {
	if n, ok := parseNumericStatus(v); ok {
		return n, true
	}
	if s, ok := v.(string); ok {
		return parseStringStatus(s)
	}
	return 0, false
}
