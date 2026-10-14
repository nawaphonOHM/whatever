package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAndValidateCORS_NonExistentFile(t *testing.T) {
	cfg, err := LoadAndValidateCORS("/path/to/nonexistent/file.json")
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read CORS configuration file")
}

func TestLoadAndValidateCORS_UnsupportedExtension(t *testing.T) {
	path := writeTempConfigFile(t, "cors.txt", "some content")
	cfg, err := LoadAndValidateCORS(path)
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "unsupported CORS configuration file format")
}

func TestLoadAndValidateCORS_InvalidJSONSyntax(t *testing.T) {
	t.Run("malformed json syntax", func(t *testing.T) {
		path := writeTempConfigFile(t, "cors.json", "{invalid-json: true,}")
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse JSON CORS configuration")
	})

	t.Run("empty json file", func(t *testing.T) {
		path := writeTempConfigFile(t, "empty.json", "")
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse JSON CORS configuration")
	})

	t.Run("trailing data after json value", func(t *testing.T) {
		path := writeTempConfigFile(t, "trailing.json", `{"allow-origins": ["https://example.com"]} extra`)
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "multiple JSON values or trailing content detected")
	})
}

func TestLoadAndValidateCORS_InvalidYAMLSyntax(t *testing.T) {
	t.Run("malformed yaml syntax", func(t *testing.T) {
		path := writeTempConfigFile(t, "cors.yaml", "allow-origins: [unclosed list")
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse YAML CORS configuration")
	})

	t.Run("empty yaml file", func(t *testing.T) {
		path := writeTempConfigFile(t, "empty.yaml", "")
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to parse YAML CORS configuration")
	})

	t.Run("multi-document yaml rejected", func(t *testing.T) {
		multiDoc := "allow-origins:\n  - https://example.com\n---\nallow-origins:\n  - https://attacker.com\n"
		path := writeTempConfigFile(t, "multidoc.yaml", multiDoc)
		cfg, err := LoadAndValidateCORS(path)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "multiple documents or trailing content detected")
	})
}
