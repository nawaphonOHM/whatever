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

// MixedMessage represents a structured log message containing both a string header and parsed JSON content.
type MixedMessage struct {
	Header  string          `json:"header"`
	Content json.RawMessage `json:"content"`
}

func parsePureJSON(s string) (json.RawMessage, bool) {
	if isJSONCandidate(s) && json.Valid([]byte(s)) {
		return json.RawMessage(s), true
	}
	return nil, false
}

// ParseGCPMessage detects valid JSON object/array messages or mixed string+JSON messages,
// or returns the plain string.
func ParseGCPMessage(msg string) any {
	if raw, ok := parsePureJSON(strings.TrimSpace(msg)); ok {
		return raw
	}
	if mixed, ok := parseMixedMessage(msg); ok {
		return mixed
	}
	return msg
}
