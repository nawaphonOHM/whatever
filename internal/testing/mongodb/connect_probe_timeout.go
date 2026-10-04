package mongodb

import (
	"context"
	"time"
)

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
	return DefaultSocketTimeout
}

func resolveProbeTimeout(o *Options) time.Duration {
	if o == nil || o.SocketTimeout <= 0 {
		return DefaultSocketTimeout
	}
	return o.SocketTimeout
}
