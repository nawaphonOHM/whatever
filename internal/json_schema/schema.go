// Package jsonschema embeds JSON schemas used for REST configuration validation.
package jsonschema

import (
	_ "embed"
)

// CORSSchema contains the embedded JSON schema for CORS configuration.
//
//go:embed cors.schema.json
var CORSSchema []byte
