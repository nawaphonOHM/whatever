package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	errPingFormat         = "failed to ping mongodb: %w"
	errCreateClientFormat = "failed to create mongodb client: %w"
)

var pingClient = defaultPingClient

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, readpref.Primary())
}

// verifyPingAndDisconnect verifies ping and disconnects if ping fails.
func verifyPingAndDisconnect(
	ctx context.Context,
	rawClient *mongo.Client,
) error {
	if err := pingClient(ctx, rawClient); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(
				fmt.Errorf(errPingFormat, err),
				disconnectErr,
			)
		}
		return fmt.Errorf(errPingFormat, err)
	}
	return nil
}

// verifyClientPing checks ping connectivity if EnablePing is true.
func verifyClientPing(ctx context.Context, o *Options, rawClient *mongo.Client) error {
	if !o.EnablePing {
		return nil
	}
	return verifyPingAndDisconnect(ctx, rawClient)
}
