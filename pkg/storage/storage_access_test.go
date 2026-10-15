package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestPublicSentinelErrors_AccessAndExpiry(t *testing.T) {
	t.Parallel()

	assert.ErrorIs(t, storage.ErrAccessDenied, contracts.ErrAccessDenied)
	assert.ErrorIs(t, storage.ErrInvalidExpiry, contracts.ErrInvalidExpiry)
}
