package client

import (
	"fmt"
	"os"
	"strings"
)

const connErrorLogFormat = "failed to connect to mongodb: %v; exiting peacefully\n"

// exitFunc is a package-level hook for os.Exit, allowing tests to intercept process termination.
var exitFunc = os.Exit

// SetExitFunc overrides the package-level exitFunc hook for testing and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	prev := exitFunc
	exitFunc = fn
	return prev
}

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

// logConnectionError outputs diagnostic message for connection failures.
func logConnectionError(err error) {
	if _, printErr := fmt.Fprintf(os.Stderr, connErrorLogFormat, err); printErr != nil {
		return
	}
}

// handleConnectionError logs the connection failure, triggers graceful termination, and returns the error.
func handleConnectionError(err error) error {
	if err == nil {
		return nil
	}
	logConnectionError(err)
	exitFunc(0)
	return err
}
