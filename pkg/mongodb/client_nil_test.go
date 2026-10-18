package mongodb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

// Constants for nil client test cases.
const (
	testNilDBParam       = "test"
	testNilCollUsers     = "users"
	expectedNilClientErr = "mongodb client is not initialized"
	expectedNilConfigErr = "mongodb config cannot be nil"
)

// TestNilClient verifies safe method execution against a nil client pointer.
func TestNilClient(t *testing.T) {
	var client *mongodb.Client
	ctx := context.Background()

	assert.Nil(t, client.Database())
	assert.Nil(t, client.Database(testNilDBParam))
	assert.Nil(t, client.Collection(testNilCollUsers))
	assert.Nil(t, client.Collection(testNilCollUsers, testNilDBParam))
	assert.Nil(t, client.RawClient())
	assert.False(t, client.IsFirestore())
	assert.ErrorIs(t, client.Ping(ctx), mongodb.ErrNilClient)
	assert.ErrorIs(t, client.Disconnect(ctx), mongodb.ErrNilClient)
}

// TestSentinelErrors verifies exported sentinel error messages.
func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, mongodb.ErrNilClient)
	assert.Equal(t, expectedNilClientErr, mongodb.ErrNilClient.Error())
	assert.NotNil(t, mongodb.ErrNilConfig)
	assert.Equal(t, expectedNilConfigErr, mongodb.ErrNilConfig.Error())
	assert.NotNil(t, mongodb.ErrNoDocuments)
	assert.True(t, errors.Is(mongodb.ErrNoDocuments, mongo.ErrNoDocuments))
}

// TestClientCollectionReturnType verifies the public collection return type.
func TestClientCollectionReturnType(t *testing.T) {
	var client *mongodb.Client
	var coll *mongodb.Collection
	assert.Nil(t, client)
	coll = client.Collection(testNilCollUsers)
	assert.Nil(t, coll)
}
