package client

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

// applyAuth sets authentication credentials on driver options if configured.
func applyAuth(opts *options.ClientOptions, c *config.Config) {
	if c.Username != "" || c.Password != "" {
		opts.SetAuth(options.Credential{
			Username: c.Username,
			Password: c.Password,
		})
	}
}

// applyConnectTimeouts sets connect and server selection timeouts.
func applyConnectTimeouts(opts *options.ClientOptions, c *config.Config) {
	if c.ConnectTimeout > 0 {
		opts.SetConnectTimeout(c.ConnectTimeout)
	}
	if c.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(c.ServerSelectionTimeout)
	}
}

// applySocketTimeout sets socket timeout on driver options.
func applySocketTimeout(opts *options.ClientOptions, c *config.Config) {
	if c.SocketTimeout > 0 {
		opts.SetTimeout(c.SocketTimeout)
	}
}

// applyPoolSizes sets max and min pool sizes on driver options.
func applyPoolSizes(opts *options.ClientOptions, c *config.Config) {
	if c.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(c.MaxPoolSize)
	}
	if c.MinPoolSize > 0 {
		opts.SetMinPoolSize(c.MinPoolSize)
	}
}

// applyPoolMetadata sets idle timeout and app name metadata.
func applyPoolMetadata(opts *options.ClientOptions, c *config.Config) {
	if c.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(c.MaxConnIdleTime)
	}
	if c.AppName != "" {
		opts.SetAppName(c.AppName)
	}
}

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

// BuildClientOptions creates official mongo-driver ClientOptions from Config
// and merges any additional driver options.
func BuildClientOptions(
	cfg *config.Config,
	extraOpts ...*options.ClientOptions,
) *options.ClientOptions {
	c := cfg
	if c == nil {
		c = config.DefaultConfig()
	}

	opts := options.Client().ApplyURI(c.URI)
	applyAuth(opts, c)
	applyConnectTimeouts(opts, c)
	applySocketTimeout(opts, c)
	applyPoolSizes(opts, c)
	applyPoolMetadata(opts, c)
	return applyExtraOptions(opts, extraOpts)
}
