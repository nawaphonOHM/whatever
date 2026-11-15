//go:build testcontainers

package storage_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const listingMaxKeys = 3

func TestLocalStackListObjects(t *testing.T) {
	requireSharedLocalStack(t)
	ctx := t.Context()
	prefix := localStackObjectKey(t, "listing/")
	keys := localStackListingKeys(prefix)
	putLocalStackListingObjects(t, ctx, keys)
	assertLocalStackPrefixListing(t, ctx, prefix, keys)
	assertLocalStackNestedListing(t, ctx, prefix)
	assertLocalStackDelimitedListing(t, ctx, prefix)
}

func localStackListingKeys(prefix string) []string {
	return []string{prefix + "root-a.txt", prefix + "root-b.txt", prefix + "nested/a.txt", prefix + "nested/b.txt"}
}

func putLocalStackListingObjects(t *testing.T, ctx context.Context, keys []string) {
	for _, key := range keys {
		_, err := sharedStorage.PutObject(ctx, storage.PutObjectInput{
			Key: key, Body: strings.NewReader("listing payload"),
		})
		require.NoError(t, err)
		cleanupLocalStackObject(t, key)
	}
}

func assertLocalStackPrefixListing(t *testing.T, ctx context.Context, prefix string, keys []string) {
	listed, err := sharedStorage.ListObjects(ctx, storage.ListObjectsInput{Prefix: prefix})
	require.NoError(t, err)
	require.Len(t, listed.Objects, len(keys))
	listedKeys := make(map[string]struct{}, len(listed.Objects))
	for _, object := range listed.Objects {
		listedKeys[object.Key] = struct{}{}
	}
	for _, key := range keys {
		require.Contains(t, listedKeys, key)
	}
}

func assertLocalStackNestedListing(t *testing.T, ctx context.Context, prefix string) {
	nested, err := sharedStorage.ListObjects(ctx, storage.ListObjectsInput{Prefix: prefix + "nested/"})
	require.NoError(t, err)
	require.Len(t, nested.Objects, 2)
	require.Empty(t, nested.Prefixes)
}

func assertLocalStackDelimitedListing(t *testing.T, ctx context.Context, prefix string) {
	delimited, err := sharedStorage.ListObjects(ctx, storage.ListObjectsInput{
		Prefix: prefix, Delimiter: "/", MaxKeys: listingMaxKeys,
	})
	require.NoError(t, err)
	require.Equal(t, listingMaxKeys, len(delimited.Objects)+len(delimited.Prefixes))
	require.Contains(t, delimited.Prefixes, prefix+"nested/")
}
