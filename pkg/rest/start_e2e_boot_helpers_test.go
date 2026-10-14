package rest

import "testing"

func startCORSE2E(t *testing.T, corsPath string) (string, <-chan error) {
	t.Helper()
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newCORSBlueprintWithFixture(corsPath))
	return pingURL(port), errCh
}
