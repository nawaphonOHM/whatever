package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testStandardURIPort = 27017
	testRemoteURIPort   = 27018
)

// TestConfig_BuildURI tests MongoDB URI construction for various configs and TLS settings.
func TestConfig_BuildURI(t *testing.T) {
	t.Run("standard mongodb scheme with TLS false", func(t *testing.T) {
		cfg := &Config{
			Host:               "localhost",
			Port:               testStandardURIPort,
			Username:           "admin",
			Password:           "secret",
			Protocol:           ProtocolMongoDB,
			UUIDRepresentation: UUIDRepresentationStandard,
		}
		expected := "mongodb://admin:secret@localhost:27017/?uuidRepresentation=standard&tls=false"
		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("standard mongodb scheme with TLS true", func(t *testing.T) {
		cfg := &Config{
			Host:               "remote-db",
			Port:               testRemoteURIPort,
			Username:           "admin",
			Password:           "secret",
			Protocol:           ProtocolMongoDB,
			UUIDRepresentation: UUIDRepresentationUnspecified,
		}
		expected := "mongodb://admin:secret@remote-db:27018/?uuidRepresentation=unspecified&tls=true"
		assert.Equal(t, expected, cfg.BuildURI(true))
	})

	t.Run("mongodb+srv scheme omits port", func(t *testing.T) {
		cfg := &Config{
			Host:               "cluster0.example.mongodb.net",
			Port:               testStandardURIPort,
			Username:           "srvuser",
			Password:           "srvpass",
			Protocol:           ProtocolMongoDBSrv,
			UUIDRepresentation: UUIDRepresentationCSharpLegacy,
		}
		expected := "mongodb+srv://srvuser:srvpass@cluster0.example.mongodb.net/" +
			"?uuidRepresentation=csharpLegacy&tls=true"
		assert.Equal(t, expected, cfg.BuildURI(true))
	})

	t.Run("special characters in credentials are URL escaped", func(t *testing.T) {
		cfg := &Config{
			Host:               "localhost",
			Port:               testStandardURIPort,
			Username:           "user@domain.com",
			Password:           "p@ss:word#1",
			Protocol:           ProtocolMongoDB,
			UUIDRepresentation: UUIDRepresentationPythonLegacy,
		}
		expected := "mongodb://user%40domain.com:p%40ss%3Aword%231@localhost:27017/" +
			"?uuidRepresentation=pythonLegacy&tls=false"
		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("defaults applied when protocol and uuid representation empty", func(t *testing.T) {
		cfg := &Config{
			Host:     "localhost",
			Port:     testStandardURIPort,
			Username: "user",
			Password: "pass",
		}
		expected := "mongodb://user:pass@localhost:27017/?uuidRepresentation=unspecified&tls=false"
		assert.Equal(t, expected, cfg.BuildURI(false))
	})
}
