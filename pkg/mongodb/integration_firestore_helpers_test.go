package mongodb_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"

	intclient "github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
)

const (
	defaultFirestorePortNum = 443
)

func getLiveFirestoreURI() string {
	if uri := os.Getenv(envFirestoreURI); uri != "" {
		return uri
	}
	return os.Getenv(envTestFirestoreURI)
}

func parsePort(portStr string) int {
	if portStr == "" {
		return defaultFirestorePortNum
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return defaultFirestorePortNum
	}
	return p
}

func applyURIAuthEnv(t *testing.T, user *url.Userinfo) {
	t.Helper()
	if user == nil {
		return
	}
	t.Setenv(envUsername, user.Username())
	if pass, ok := user.Password(); ok {
		t.Setenv(envPassword, pass)
	}
}

func applyParsedURIEnv(t *testing.T, u *url.URL) {
	t.Helper()
	port := parsePort(u.Port())
	t.Setenv(envHost, u.Hostname())
	t.Setenv(envPort, strconv.Itoa(port))
	applyURIAuthEnv(t, u.User)
	if db := strings.TrimPrefix(u.Path, "/"); db != "" {
		t.Setenv("OHM9996_MONGODB_DATABASE", db)
	}
}

func setupLiveFirestoreEnv(t *testing.T, rawURI string) {
	t.Helper()
	u, err := url.Parse(rawURI)
	if err != nil {
		return
	}
	applyParsedURIEnv(t, u)
}

func setupFirestoreDemoEnv(t *testing.T) {
	t.Helper()
	t.Setenv(envHost, testFirestoreDemoHost)
	t.Setenv(envPort, "443")
	t.Setenv(envUsername, testUser)
	t.Setenv(envPassword, testPassword)
	t.Setenv("OHM9996_MONGODB_DATABASE", testFirestoreDemoDB)
}

func setupMockFirestorePingProbe(onPing func()) func() {
	restorePing := intclient.SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		if onPing != nil {
			onPing()
		}
		return nil
	})
	restoreProbe := intclient.SetMockProbe(func(context.Context, *intclient.Client, string) error {
		return nil
	})
	return func() {
		restorePing()
		restoreProbe()
	}
}

func setupResetFallbackMocks(attempts *int) func() {
	restorePing := intclient.SetMockPing(func(context.Context, *mongo.Client) error {
		*attempts++
		if *attempts == 1 {
			return errors.New("incomplete read of full message: read tcp: connection reset by peer")
		}
		return nil
	})
	restoreProbe := intclient.SetMockProbe(func(context.Context, *intclient.Client, string) error {
		return nil
	})
	return func() {
		restorePing()
		restoreProbe()
	}
}
