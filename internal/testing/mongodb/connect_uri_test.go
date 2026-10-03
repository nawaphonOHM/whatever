package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectURI_EmptyURI(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectURI(ctx, "")
	require.ErrorIs(t, err, ErrEmptyURI)
	assert.Nil(t, client)

	client, err = ConnectURI(ctx, "   ")
	require.ErrorIs(t, err, ErrEmptyURI)
	assert.Nil(t, client)
}

func TestConnectURI_ExtractDatabase(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectURI(ctx, "mongodb://localhost:27017/extracted_db", WithPing(false))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "extracted_db", client.Database().Name())
}

func TestConnectURI_OverrideDatabase(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectURI(
		ctx,
		"mongodb://localhost:27017/extracted_db",
		WithPing(false),
		WithDatabase("overridden_db"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "overridden_db", client.Database().Name())
}
