package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// applyMaxPoolSize configures maximum connection pool capacity.
func applyMaxPoolSize(opts *options.ClientOptions, maxPool uint64) {
	if maxPool > 0 {
		opts.SetMaxPoolSize(maxPool)
	}
}

// applyMinPoolSize configures minimum connection pool capacity.
func applyMinPoolSize(opts *options.ClientOptions, minPool uint64) {
	if minPool > 0 {
		opts.SetMinPoolSize(minPool)
	}
}

// applyIdleTime configures maximum connection idle duration.
func applyIdleTime(opts *options.ClientOptions, idle time.Duration) {
	if idle > 0 {
		opts.SetMaxConnIdleTime(idle)
	}
}

// applyDirectConnection configures direct connection if enabled on Options.
func applyDirectConnection(opts *options.ClientOptions, o *Options) {
	if o != nil && o.DirectConnection {
		opts.SetDirect(true)
	}
}

// applyPoolAndDirect configures connection pool limits and direct connection.
func applyPoolAndDirect(opts *options.ClientOptions, o *Options) {
	if o == nil {
		return
	}
	applyMaxPoolSize(opts, o.MaxPoolSize)
	applyMinPoolSize(opts, o.MinPoolSize)
	applyIdleTime(opts, o.MaxConnIdleTime)
	applyDirectConnection(opts, o)
}
