package storage_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

type exitTracker struct {
	code   int
	called bool
}

func setupExitTracker() (*exitTracker, func()) {
	tracker := &exitTracker{}
	prev := storage.SetExitFunc(func(code int) {
		tracker.code = code
		tracker.called = true
	})
	cleanup := func() {
		_ = storage.SetExitFunc(prev)
	}
	return tracker, cleanup
}

func serverErrorHTTPClient() *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     http.StatusText(http.StatusInternalServerError),
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}, nil
		}),
	}
}

func probeFailureOpts() []storage.Option {
	return []storage.Option{
		storage.WithHTTPClient(serverErrorHTTPClient()),
		storage.WithDriverOptions(func(opts *s3.Options) { opts.RetryMaxAttempts = 1 }),
		storage.WithSkipProbe(false),
	}
}

func TestPublicStorage_SetExitFunc_ProbeFailure(t *testing.T) {
	tracker, cleanup := setupExitTracker()
	defer cleanup()

	c, err := storage.ConnectWithConfig(
		context.Background(),
		testStaticConfig(),
		probeFailureOpts()...,
	)

	assert.Error(t, err)
	assert.Nil(t, c)
	assert.True(t, tracker.called)
	assert.NotZero(t, tracker.code)
}
