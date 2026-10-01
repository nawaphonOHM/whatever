package mongodb

import (
	"github.com/testcontainers/testcontainers-go"
)

// Option defines a functional option for configuring a MongoDB container.
type Option func(*Options)

// Options holds configuration settings for running a MongoDB container.
type Options struct {
	Image            string
	Username         string
	Password         string
	Database         string
	ReplicaSet       string
	Env              map[string]string
	ContainerOptions []testcontainers.ContainerCustomizer
}

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return &Options{
		Image: DefaultImage,
		Env:   make(map[string]string),
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
