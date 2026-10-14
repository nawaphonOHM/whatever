package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type schemaTestCase struct {
	name    string
	content string
	wantErr string
}

func runSchemaTestCases(t *testing.T, cases []schemaTestCase) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempConfigFile(t, "cors.json", tc.content)
			cfg, err := LoadAndValidateCORS(path)
			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), schemaValidationErr)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func runYAMLSchemaTestCases(t *testing.T, cases []schemaTestCase) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempConfigFile(t, "cors.yaml", tc.content)
			cfg, err := LoadAndValidateCORS(path)
			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), schemaValidationErr)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
