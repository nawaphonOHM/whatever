package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/nawaphonOHM/whatever/internal/logging/central"
	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

// SetExitFunc overrides the process exit hook used by the central logger worker
// and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	prev := central.DefaultWorker().ExitFunc()
	central.DefaultWorker().SetExitFunc(fn)
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

// handleConnectionError logs the connection failure, triggers abnormal termination via Central Log,
// and returns the error.
func handleConnectionError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf("failed to connect to mongodb: %v; exiting abnormally", err)
	central.ExitWithAbnormal(ctx, msg, err, nil)
	return err
}

// fallbackTLSAttempt attempts connection with TLS enabled upon TLS requirement error.
func fallbackTLSAttempt(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	central.DefaultWorker().Logger().InfoContext(ctx, "server requires TLS; attempting connection with TLS enabled",
		"host", cfg.Host,
		"port", cfg.Port,
	)
	tlsClient, tlsErr := attemptConnection(ctx, cfg, true, opts...)
	if tlsErr != nil {
		return nil, handleConnectionError(ctx, tlsErr)
	}
	return tlsClient, nil
}
