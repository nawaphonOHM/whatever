package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFirestoreDefaultPort = 27017
	testFirestorePort443     = 443
	testFirestorePort8443    = 8443
	testFirestoreHost        = "project.firestore.goog"
	testFirestoreUser        = "user"
	testFirestorePassword    = "pwd"
	testFirestoreAuthMech    = "&authMechanism=SCRAM-SHA-256"
)

// TestConfig_BuildURI_Firestore tests URI generation with Firestore defaults.
func TestConfig_BuildURI_Firestore(t *testing.T) {
	t.Run("default port 27017 mapped to port 443 with all firestore query params", func(t *testing.T) {
		cfg := &Config{
			Host:               "fc28fdab-9482-4359-9d76-0bb5f9fd1c4d.asia-southeast1.firestore.goog",
			Port:               testFirestoreDefaultPort,
			Username:           "testuser",
			Password:           "secretpass",
			UUIDRepresentation: UUIDRepresentationStandard,
		}
		expected := "mongodb://testuser:secretpass@" +
			"fc28fdab-9482-4359-9d76-0bb5f9fd1c4d.asia-southeast1.firestore.goog:443/" +
			"?uuidRepresentation=standard&tls=true&loadBalanced=true&retryWrites=false" +
			testFirestoreAuthMech

		assert.Equal(t, expected, cfg.BuildURI(false))
		assert.Equal(t, expected, cfg.BuildURI(true))
		assert.Equal(t, expected, cfg.buildURIFromFields(false))
	})

	t.Run("port 0 defaults to 443 for firestore target", func(t *testing.T) {
		cfg := &Config{
			Host:     testFirestoreHost,
			Port:     0,
			Username: testFirestoreUser,
			Password: testFirestorePassword,
		}
		expected := "mongodb://user:pwd@project.firestore.goog:443/" +
			"?uuidRepresentation=unspecified&tls=true&loadBalanced=true&retryWrites=false" +
			testFirestoreAuthMech

		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("explicit port 443 preserved", func(t *testing.T) {
		cfg := &Config{
			Host:     testFirestoreHost,
			Port:     testFirestorePort443,
			Username: testFirestoreUser,
			Password: testFirestorePassword,
		}
		expected := "mongodb://user:pwd@project.firestore.goog:443/" +
			"?uuidRepresentation=unspecified&tls=true&loadBalanced=true&retryWrites=false" +
			testFirestoreAuthMech

		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("custom port preserved if explicitly specified", func(t *testing.T) {
		cfg := &Config{
			Host:     testFirestoreHost,
			Port:     testFirestorePort8443,
			Username: testFirestoreUser,
			Password: testFirestorePassword,
		}
		expected := "mongodb://user:pwd@project.firestore.goog:8443/" +
			"?uuidRepresentation=unspecified&tls=true&loadBalanced=true&retryWrites=false" +
			testFirestoreAuthMech

		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("includes authSource and appName when provided", func(t *testing.T) {
		cfg := &Config{
			Host:               testFirestoreHost,
			Port:               testFirestorePort443,
			Username:           testFirestoreUser,
			Password:           testFirestorePassword,
			AuthSource:         "admin",
			AppName:            "firestore-service",
			UUIDRepresentation: UUIDRepresentationStandard,
		}
		expected := "mongodb://user:pwd@project.firestore.goog:443/" +
			"?uuidRepresentation=standard&tls=true&loadBalanced=true&retryWrites=false" +
			testFirestoreAuthMech + "&authSource=admin&appName=firestore-service"

		assert.Equal(t, expected, cfg.BuildURI(false))
	})
}
