package core

import (
	"encoding/json"
	"strings"
)

func parseCandidateJSON(headerPart string, jsonPart []byte) (MixedMessage, bool) {
	if !validJSONCandidate(jsonPart) {
		return MixedMessage{}, false
	}
	return MixedMessage{
		Header:  headerPart,
		Content: json.RawMessage(jsonPart),
	}, true
}

func validJSONCandidate(jsonPart []byte) bool {
	if len(jsonPart) < minCandidateLen {
		return false
	}
	if !isObjectOrArray(jsonPart[0], jsonPart[len(jsonPart)-1]) {
		return false
	}
	return json.Valid(jsonPart)
}

func isCandidateBoundary(b byte) bool {
	return b == '{' || b == '['
}

func tryCandidateAt(msg string, msgBytes []byte, idx int) (MixedMessage, bool) {
	if !isCandidateBoundary(msgBytes[idx]) {
		return MixedMessage{}, false
	}
	headerPart := strings.TrimSpace(msg[:idx])
	if headerPart == "" {
		return MixedMessage{}, false
	}
	return parseCandidateJSON(headerPart, msgBytes[idx:])
}

func parseMixedMessage(msg string) (MixedMessage, bool) {
	trimmed := strings.TrimSpace(msg)
	if !hasMixedJSONBoundary(trimmed) {
		return MixedMessage{}, false
	}
	return scanMixedMessage(trimmed)
}

func hasMixedJSONBoundary(msg string) bool {
	if len(msg) < minCandidateLen {
		return false
	}
	return isClosingBoundary(msg[len(msg)-1])
}

func scanMixedMessage(msg string) (MixedMessage, bool) {
	msgBytes := []byte(msg)
	for i := range msgBytes {
		if mixed, ok := tryCandidateAt(msg, msgBytes, i); ok {
			return mixed, true
		}
	}
	return MixedMessage{}, false
}

func isClosingBoundary(b byte) bool {
	return b == '}' || b == ']'
}
