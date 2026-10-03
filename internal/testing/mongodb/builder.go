package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// applyExtraOptions merges additional driver options if provided.
func applyExtraOptions(
	opts *options.ClientOptions,
	extraOpts []*options.ClientOptions,
) *options.ClientOptions {
	if len(extraOpts) == 0 {
		return opts
	}
	allOpts := append([]*options.ClientOptions{opts}, extraOpts...)
	return options.MergeClientOptions(allOpts...)
}

// applyConnectTimeout sets connection timeout.
func applyConnectTimeout(opts *options.ClientOptions, timeout time.Duration) {
	if timeout > 0 {
		opts.SetConnectTimeout(timeout)
	}
}

// applyServerSelectionTimeout sets server selection timeout.
func applyServerSelectionTimeout(opts *options.ClientOptions, timeout time.Duration) {
	if timeout > 0 {
		opts.SetServerSelectionTimeout(timeout)
	}
}

// applySocketTimeout sets socket operation timeout.
func applySocketTimeout(opts *options.ClientOptions, timeout time.Duration) {
	if timeout > 0 {
		opts.SetTimeout(timeout)
	}
}

// applyTimeouts configures connection and socket timeouts.
func applyTimeouts(opts *options.ClientOptions, o *Options) {
	if o == nil {
		return
	}
	applyConnectTimeout(opts, o.ConnectTimeout)
	applyServerSelectionTimeout(opts, o.ServerSelectionTimeout)
	applySocketTimeout(opts, o.SocketTimeout)
}

// hasURI reports whether Options contains a non-empty URI.
func hasURI(o *Options) bool {
	return o != nil && o.URI != ""
}

// collectDriverOptions gathers custom driver options without modifying parameter.
func collectDriverOptions(o *Options, extraOpts []*options.ClientOptions) []*options.ClientOptions {
	if o == nil || len(o.DriverOptions) == 0 {
		return extraOpts
	}
	return append(extraOpts, o.DriverOptions...)
}

// applyOptionsConfig configures timeouts, pool limits, and app name.
func applyOptionsConfig(opts *options.ClientOptions, o *Options) {
	applyTimeouts(opts, o)
	applyPoolAndDirect(opts, o)
	if o != nil && o.AppName != "" {
		opts.SetAppName(o.AppName)
	}
}

// BuildClientOptions creates official mongo-driver ClientOptions from Options
// and merges any additional driver options.
func BuildClientOptions(
	o *Options,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	opts := options.Client()
	opts.ApplyURI(resolveURI(o))
	applyOptionsConfig(opts, o)
	combined := collectDriverOptions(o, extraOpts)
	return applyExtraOptions(opts, combined)
}

// BuildClientOptionsWithTLS creates official mongo-driver ClientOptions from
// Options with explicit TLS setting and merges any additional driver options.
func BuildClientOptionsWithTLS(
	o *Options,
	enableTLS bool,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	optsCopy := DefaultOptions()
	if o != nil {
		*optsCopy = *o
	}
	optsCopy.EnableTLS = enableTLS
	return BuildClientOptions(optsCopy, extraOpts...)
}
