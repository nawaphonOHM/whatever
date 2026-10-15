package client_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	concurrentGoroutines = 8
	concurrentIterations = 20
)

func stepPresign(t *testing.T, c *client.Client, i int) {
	t.Helper()
	c.SetPresignExpiry(time.Duration(i+1) * time.Minute)
	_, getErr := c.PresignGetObject(t.Context(), contracts.PresignGetObjectInput{Key: presignKey})
	_, putErr := c.PresignPutObject(t.Context(), contracts.PresignPutObjectInput{Key: presignKey})
	require.NoError(t, getErr)
	require.NoError(t, putErr)
}

func runPresignRacer(t *testing.T, c *client.Client, wg *sync.WaitGroup) {
	t.Helper()
	defer wg.Done()
	for i := range concurrentIterations {
		stepPresign(t, c, i)
	}
}

func TestPresign_ConcurrentSetExpiryAndSign(t *testing.T) {
	t.Parallel()

	c := newTestPresignClient()
	var wg sync.WaitGroup
	wg.Add(concurrentGoroutines)
	for range concurrentGoroutines {
		go runPresignRacer(t, c, &wg)
	}
	wg.Wait()
}
