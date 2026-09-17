package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Option defines a functional option for customizing internal MongoDB client options.
type Option func(*Options)

// Options holds internal driver configuration options.
type Options struct {
	DriverOptions []*options.ClientOptions
}

// NewOptions creates an Options struct by evaluating internal functional options.
func NewOptions(opts ...Option) *Options {
	o := &Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}

// WithDriverOptions appends raw mongo-driver ClientOptions internally.
func WithDriverOptions(opts ...*options.ClientOptions) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.DriverOptions = append(o.DriverOptions, opt)
			}
		}
	}
}

// BuildClientOptions creates official mongo-driver ClientOptions from internal Config
// and merges any additional internal driver options.
func BuildClientOptions(cfg *Config, extraOpts ...*options.ClientOptions) *options.ClientOptions {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	opts := options.Client().ApplyURI(cfg.URI)

	if cfg.ConnectTimeout > 0 {
		opts.SetConnectTimeout(cfg.ConnectTimeout)
	}
	if cfg.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
	}
	if cfg.SocketTimeout > 0 {
		opts.SetTimeout(cfg.SocketTimeout)
	}
	if cfg.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		opts.SetMinPoolSize(cfg.MinPoolSize)
	}
	if cfg.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	}
	if cfg.AppName != "" {
		opts.SetAppName(cfg.AppName)
	}

	if len(extraOpts) > 0 {
		allOpts := append([]*options.ClientOptions{opts}, extraOpts...)
		opts = options.MergeClientOptions(allOpts...)
	}

	return opts
}
