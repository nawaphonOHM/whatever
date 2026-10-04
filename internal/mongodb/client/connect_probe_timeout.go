package client

import (
	"context"
	"time"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

const defaultProbeTimeout = 10 * time.Second

type probeTimeoutKey struct{}

func withProbeTimeout(ctx context.Context, timeout time.Duration) context.Context {
	if timeout <= 0 {
		return ctx
	}
	return context.WithValue(ctx, probeTimeoutKey{}, timeout)
}

func probeTimeoutFromContext(ctx context.Context) time.Duration {
	if v, ok := ctx.Value(probeTimeoutKey{}).(time.Duration); ok && v > 0 {
		return v
	}
	return defaultProbeTimeout
}

func resolveProbeTimeout(cfg *config.Config) time.Duration {
	if cfg == nil || cfg.SocketTimeout <= 0 {
		return defaultProbeTimeout
	}
	return cfg.SocketTimeout
}
