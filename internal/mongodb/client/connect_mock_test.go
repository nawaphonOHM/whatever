package client

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupMockPingAndProbe(probeErr error) func() {
	cleanupPing := SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	})
	cleanupProbe := SetMockProbe(func(context.Context, *mongo.Client, string) error {
		return probeErr
	})
	return func() {
		cleanupProbe()
		cleanupPing()
	}
}
