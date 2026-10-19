package contracts_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestSentinelErrors_Messages(t *testing.T) {
	t.Parallel()

	require.Equal(t, "storage object not found", contracts.ErrObjectNotFound.Error())
	require.Equal(t, "storage bucket not found", contracts.ErrBucketNotFound.Error())
	require.Equal(t, "storage client is not initialized", contracts.ErrNilClient.Error())
	require.Equal(t, "storage config cannot be nil", contracts.ErrNilConfig.Error())
}

func TestSentinelErrors_AccessAndExpiry(t *testing.T) {
	t.Parallel()

	require.Equal(t, "storage access denied", contracts.ErrAccessDenied.Error())
	require.Equal(t, "storage presign expiry is invalid", contracts.ErrInvalidExpiry.Error())
	assert.False(t, errors.Is(contracts.ErrAccessDenied, contracts.ErrObjectNotFound))
	assert.False(t, errors.Is(contracts.ErrInvalidExpiry, contracts.ErrAccessDenied))
}

func TestSentinelErrors_Distinct(t *testing.T) {
	t.Parallel()

	pairs := []struct {
		first, second error
	}{
		{contracts.ErrObjectNotFound, contracts.ErrBucketNotFound},
		{contracts.ErrObjectNotFound, contracts.ErrNilClient},
		{contracts.ErrObjectNotFound, contracts.ErrNilConfig},
		{contracts.ErrBucketNotFound, contracts.ErrNilClient},
		{contracts.ErrBucketNotFound, contracts.ErrNilConfig},
		{contracts.ErrNilClient, contracts.ErrNilConfig},
	}

	for _, p := range pairs {
		assert.False(t, errors.Is(p.first, p.second))
		assert.False(t, errors.Is(p.second, p.first))
	}
}
