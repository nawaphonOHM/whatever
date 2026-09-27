package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig_BuildURI_Metadata tests URI construction with authSource and appName.
func TestConfig_BuildURI_Metadata(t *testing.T) {
	t.Run("authSource and appName included in query string", func(t *testing.T) {
		cfg := &Config{
			Host:               "localhost",
			Port:               testStandardURIPort,
			Username:           "user",
			Password:           "pass",
			AuthSource:         "admin",
			AppName:            "my-service",
			Protocol:           ProtocolMongoDB,
			UUIDRepresentation: UUIDRepresentationStandard,
		}
		expected := "mongodb://user:pass@localhost:27017/" +
			"?uuidRepresentation=standard&tls=false&authSource=admin&appName=my-service"
		assert.Equal(t, expected, cfg.BuildURI(false))
	})

	t.Run("mongodb+srv with authSource, appName, and TLS", func(t *testing.T) {
		cfg := &Config{
			Host:               "atlas.example.com",
			Username:           "clouduser",
			Password:           "cloudpass",
			AuthSource:         "admin",
			AppName:            "atlas-app",
			Protocol:           ProtocolMongoDBSrv,
			UUIDRepresentation: UUIDRepresentationStandard,
		}
		expected := "mongodb+srv://clouduser:cloudpass@atlas.example.com/" +
			"?uuidRepresentation=standard&tls=true&authSource=admin&appName=atlas-app"
		assert.Equal(t, expected, cfg.BuildURI(true))
	})
}
