package client_test

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	maxExpiry     = 7 * 24 * time.Hour
	defaultExpiry = 15 * time.Minute
	signedHeaders = "X-Amz-SignedHeaders"
	presignKey    = "k"
)

func presignedQuery(t *testing.T, rawURL string) url.Values {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u.Query()
}

func TestPresign_ExpiryBounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		expiry  time.Duration
		want    time.Duration
		wantErr bool
	}{
		{name: "default when zero", expiry: 0, want: defaultExpiry},
		{name: "exactly seven days", expiry: maxExpiry, want: maxExpiry},
		{name: "above seven days", expiry: maxExpiry + time.Second, wantErr: true},
		{name: "negative", expiry: -time.Second, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newTestPresignClient()
			getIn := contracts.PresignGetObjectInput{Key: presignKey, ExpiresIn: tc.expiry}
			putIn := contracts.PresignPutObjectInput{Key: presignKey, ExpiresIn: tc.expiry}
			get, getErr := c.PresignGetObject(t.Context(), getIn)
			put, putErr := c.PresignPutObject(t.Context(), putIn)
			if tc.wantErr {
				require.ErrorIs(t, getErr, contracts.ErrInvalidExpiry)
				require.ErrorIs(t, putErr, contracts.ErrInvalidExpiry)
				return
			}
			require.NoError(t, getErr)
			require.NoError(t, putErr)
			seconds := strconv.Itoa(int(tc.want.Seconds()))
			assert.Equal(t, seconds, presignedQuery(t, get.URL).Get("X-Amz-Expires"))
			assert.Equal(t, seconds, presignedQuery(t, put.URL).Get("X-Amz-Expires"))
			assert.Equal(t, tc.want, get.ExpiresIn)
		})
	}
}

func TestPresign_OversizedDefaultIsRejected(t *testing.T) {
	t.Parallel()

	c := newTestPresignClient()
	c.SetPresignExpiry(maxExpiry + time.Hour)
	_, err := c.PresignGetObject(t.Context(), contracts.PresignGetObjectInput{Key: presignKey})
	assert.ErrorIs(t, err, contracts.ErrInvalidExpiry)
}

func TestPresignPut_ContentTypeIsSignedAndReturned(t *testing.T) {
	t.Parallel()

	input := contracts.PresignPutObjectInput{Key: presignKey, ContentType: "image/png"}
	signed, err := newTestPresignClient().PresignPutObject(t.Context(), input)
	require.NoError(t, err)
	assert.Equal(t, []string{"image/png"}, signed.Headers["Content-Type"])
	assert.Contains(t, presignedQuery(t, signed.URL).Get(signedHeaders), "content-type")
	assert.NotContains(t, signed.Headers, "Authorization")
}

func TestPresignPut_NoContentTypeSignsHostOnly(t *testing.T) {
	t.Parallel()

	input := contracts.PresignPutObjectInput{Key: presignKey}
	plain, err := newTestPresignClient().PresignPutObject(t.Context(), input)
	require.NoError(t, err)
	assert.NotContains(t, presignedQuery(t, plain.URL).Get(signedHeaders), "content-type")
	assert.NotContains(t, plain.Headers, "Content-Type")
}
