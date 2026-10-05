package client

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetMockPing overrides pingClient for testing and returns a restore function.
func SetMockPing(fn func(context.Context, *mongo.Client) error) func() {
	prev := pingClient
	if fn != nil {
		pingClient = fn
	} else {
		pingClient = defaultPingClient
	}
	return func() {
		pingClient = prev
	}
}

// SetMockPingFirestore overrides pingFirestoreClient for testing and returns a restore function.
func SetMockPingFirestore(fn func(context.Context, *mongo.Client, string) error) func() {
	prev := pingFirestoreClient
	if fn != nil {
		pingFirestoreClient = fn
	} else {
		pingFirestoreClient = defaultPingFirestoreClient
	}
	return func() {
		pingFirestoreClient = prev
	}
}

// SetMockProbe overrides probeClient for testing and returns a restore function.
func SetMockProbe(fn func(context.Context, *mongo.Client, string) error) func() {
	prev := probeClient
	if fn != nil {
		probeClient = fn
	} else {
		probeClient = defaultProbeClient
	}
	return func() {
		probeClient = prev
	}
}
