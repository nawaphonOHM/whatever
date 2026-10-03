package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Option defines a functional option for configuring MongoDB test connections.
type Option func(*Options)

// Options holds configuration settings for test MongoDB connections.
type Options struct {
	URI                    string
	Host                   string
	Protocol               string
	Database               string
	Username               string
	Password               string
	AuthSource             string
	AppName                string
	UUIDRepresentation     string
	DriverOptions          []*options.ClientOptions
	ConnectTimeout         time.Duration
	ServerSelectionTimeout time.Duration
	SocketTimeout          time.Duration
	MaxConnIdleTime        time.Duration
	MaxPoolSize            uint64
	MinPoolSize            uint64
	Port                   int
	EnablePing             bool
	EnableTLS              bool
	DirectConnection       bool
}

// DefaultOptions returns an Options struct populated with recommended test defaults.
func DefaultOptions() *Options {
	return &Options{
		Host:                   DefaultHost,
		Port:                   DefaultPort,
		Protocol:               DefaultProtocol,
		UUIDRepresentation:     DefaultUUIDRepresentation,
		ConnectTimeout:         DefaultConnectTimeout,
		ServerSelectionTimeout: DefaultServerSelectionTimeout,
		SocketTimeout:          DefaultSocketTimeout,
		MaxConnIdleTime:        DefaultMaxConnIdleTime,
		MaxPoolSize:            DefaultMaxPoolSize,
		MinPoolSize:            DefaultMinPoolSize,
		EnablePing:             true,
		EnableTLS:              false,
		DirectConnection:       false,
	}
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	o := DefaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}
