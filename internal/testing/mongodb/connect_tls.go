package mongodb

import (
	"strings"
)

var tlsErrorPatterns = []string{
	"server requires tls",
	"server requires ssl",
	"ssl handshake",
	"tls handshake",
	"first record does not look like a tls handshake",
	"command requires ssl",
	"command requires tls",
	"requires ssl",
	"requires tls",
	"ssl required",
	"tls required",
	"connection closed",
	"connection reset by peer",
	"incomplete read of full message",
	"broken pipe",
	"server selection error",
}

// matchTLSPattern matches error message substrings against known TLS error indicators.
func matchTLSPattern(msg string) bool {
	for _, pattern := range tlsErrorPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}

// isTLSError checks whether an error indicates that the server requires a TLS connection.
func isTLSError(err error) bool {
	if err == nil {
		return false
	}
	return matchTLSPattern(strings.ToLower(err.Error()))
}
