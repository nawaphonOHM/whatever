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

func applyFirestoreOptions(opts *options.ClientOptions, c *config.Config) {
	if c.IsFirestore() {
		opts.SetLoadBalanced(true).
			SetRetryWrites(false)
	}
}

func applyAppOptions(opts *options.ClientOptions, c *config.Config) {
	if c.AppName != "" {
		opts.SetAppName(c.AppName)
	}
}

func resolveConfig(cfg *config.Config) *config.Config {
	if cfg == nil {
		return config.DefaultConfig()
	}
	return cfg
}

// BuildClientOptionsWithTLS creates official mongo-driver ClientOptions from
// Config with explicit TLS setting and merges any additional driver options.
func BuildClientOptionsWithTLS(
	cfg *config.Config,
	enableTLS bool,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	c := resolveConfig(cfg)
	opts := options.Client().
		ApplyURI(c.BuildURI(enableTLS)).
		SetConnectTimeout(c.ConnectTimeout).
		SetServerSelectionTimeout(c.ServerSelectionTimeout).
		SetTimeout(c.SocketTimeout).
		SetMaxConnIdleTime(c.MaxConnIdleTime).
		SetMaxPoolSize(c.MaxPoolSize).
		SetMinPoolSize(c.MinPoolSize)

	applyFirestoreOptions(opts, c)
	applyAppOptions(opts, c)

	return applyExtraOptions(opts, extraOpts)
}

func resolveDefaultTLS(cfg *config.Config) bool {
	return cfg != nil && cfg.IsFirestore()
}

// BuildClientOptions creates official mongo-driver ClientOptions from Config
// with TLS disabled by default (or enabled if Firestore) and merges any additional driver options.
func BuildClientOptions(
	cfg *config.Config,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	return BuildClientOptionsWithTLS(cfg, resolveDefaultTLS(cfg), extraOpts...)
}
