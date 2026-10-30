package core

import (
	"encoding/json"
	"strings"
)

const minCandidateLen = 2

func isObjectOrArray(first, last byte) bool {
	if first == '{' {
		return last == '}'
	}
	if first == '[' {
		return last == ']'
	}
	return false
}

func isJSONCandidate(s string) bool {
	if len(s) < minCandidateLen {
		return false
	}
	return isObjectOrArray(s[0], s[len(s)-1])
}

// ParseGCPMessage detects valid JSON object/array messages or returns the plain string.
func ParseGCPMessage(msg string) any {
	trimmed := strings.TrimSpace(msg)
	if isJSONCandidate(trimmed) && json.Valid([]byte(trimmed)) {
		return json.RawMessage(trimmed)
	}
	return msg
}
