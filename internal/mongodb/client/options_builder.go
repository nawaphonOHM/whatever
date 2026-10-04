package client

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
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

// BuildClientOptionsWithTLS creates official mongo-driver ClientOptions from
// Config with explicit TLS setting and merges any additional driver options.
func BuildClientOptionsWithTLS(
	cfg *config.Config,
	enableTLS bool,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	c := cfg
	if c == nil {
		c = config.DefaultConfig()
	}

	opts := options.Client().
		ApplyURI(c.BuildURI(enableTLS)).
		SetConnectTimeout(c.ConnectTimeout).
		SetServerSelectionTimeout(c.ServerSelectionTimeout).
		SetTimeout(c.SocketTimeout).
		SetMaxConnIdleTime(c.MaxConnIdleTime).
		SetMaxPoolSize(c.MaxPoolSize).
		SetMinPoolSize(c.MinPoolSize)

	if c.AppName != "" {
		opts.SetAppName(c.AppName)
	}

	return applyExtraOptions(opts, extraOpts)
}

// BuildClientOptions creates official mongo-driver ClientOptions from Config
// with TLS disabled by default and merges any additional driver options.
func BuildClientOptions(
	cfg *config.Config,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	return BuildClientOptionsWithTLS(cfg, false, extraOpts...)
}
