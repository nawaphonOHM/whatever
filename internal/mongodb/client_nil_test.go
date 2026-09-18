package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// verifyNilClientMethods verifies that methods return nil or ErrNilClient
// safely.
func verifyNilClientMethods(t *testing.T, c *Client) {
	ctx := context.Background()

	assert.Nil(t, c.Database())
	assert.Nil(t, c.Database("test"))
	assert.Nil(t, c.Collection(testCollUsers))
	assert.Nil(t, c.Collection(testCollUsers, "test"))
	assert.Nil(t, c.RawClient())
	assert.ErrorIs(t, c.Ping(ctx), ErrNilClient)
	assert.ErrorIs(t, c.Disconnect(ctx), ErrNilClient)
}

// TestClient_NilSafety tests method calls on nil and uninitialized clients.
func TestClient_NilSafety(t *testing.T) {
	var nilClient *Client
	verifyNilClientMethods(t, nilClient)

	uninitializedClient := &Client{
		rawClient:       nil,
		defaultDatabase: testDefaultDB,
	}
	verifyNilClientMethods(t, uninitializedClient)
}
