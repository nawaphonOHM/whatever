package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
)

func TestBuildGRPCOptions_Variants(t *testing.T) {
	cfgWithScheme := &config.Config{
		Endpoint: "dns:///localhost:4317",
		Insecure: true,
	}
	opts1 := buildGRPCOptions(cfgWithScheme)
	assert.NotEmpty(t, opts1)

	cfgWithoutSchemeSecure := &config.Config{
		Endpoint: "localhost:4317",
		Insecure: false,
	}
	opts2 := buildGRPCOptions(cfgWithoutSchemeSecure)
	assert.NotEmpty(t, opts2)
}

func TestBuildHTTPOptions_Variants(t *testing.T) {
	cfgWithSchemeSecure := &config.Config{
		Endpoint: "https://localhost:4318/v1/traces",
		Protocol: config.ProtocolHTTPJSON,
		Insecure: false,
	}
	opts1 := buildHTTPOptions(cfgWithSchemeSecure)
	assert.NotEmpty(t, opts1)

	cfgWithoutSchemeInsecure := &config.Config{
		Endpoint: "localhost:4318",
		Protocol: config.ProtocolHTTPProtobuf,
		Insecure: true,
	}
	opts2 := buildHTTPOptions(cfgWithoutSchemeInsecure)
	assert.NotEmpty(t, opts2)
}
