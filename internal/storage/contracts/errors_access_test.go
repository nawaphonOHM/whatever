package contracts_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestSentinelErrors_AccessAndExpiry(t *testing.T) {
	t.Parallel()

	require.Equal(t, "storage access denied", contracts.ErrAccessDenied.Error())
	require.Equal(t, "storage presign expiry is invalid", contracts.ErrInvalidExpiry.Error())
	assert.False(t, errors.Is(contracts.ErrAccessDenied, contracts.ErrObjectNotFound))
	assert.False(t, errors.Is(contracts.ErrInvalidExpiry, contracts.ErrAccessDenied))
}
