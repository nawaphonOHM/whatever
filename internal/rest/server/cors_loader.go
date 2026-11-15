package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/middleware"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

const (
	corsYAMLEnvName      = "OHM9996_SERVER_CORS_YML_YAML"
	corsJSONEnvName      = "OHM9996_SERVER_CORS_JSON"
	msgCORSNotConfigured = "no CORS environment configuration found; CORS is disabled"
	msgCORSInvalid       = "invalid CORS environment configuration; CORS is disabled"
)

func readAndParseCORS(path string) ([]byte, *CORSFileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CORS configuration file %q: %w", path, err)
	}
	return parseCORSFile(path, data)
}

func loadValidCORS(path string) (*middleware.CORSConfig, error) {
	_, fileCfg, err := readAndParseCORS(path)
	if err != nil {
		return nil, err
	}
	return buildCORSConfig(fileCfg), nil
}

// LoadAndValidateCORS reads a CORS configuration file (.json, .yaml, or .yml),
// validates it against the embedded CORS schema, and builds a CORSConfig.
func LoadAndValidateCORS(path string) (*middleware.CORSConfig, error) {
	if path == "" {
		return middleware.DefaultCORSConfig(), nil
	}
	return loadValidCORS(path)
}

type corsEnvironmentPolicy struct {
	parser func(string, []byte) ([]byte, *CORSFileConfig, error)
	name   string
	value  string
}

func loadCORSConfig(cfg *Config) *middleware.CORSConfig {
	for _, policy := range corsEnvironmentPolicies(cfg) {
		if corsConfig, configured := loadEnvironmentCORS(policy); configured {
			return corsConfig
		}
	}

	logging.Info(msgCORSNotConfigured)
	return middleware.DefaultCORSConfig()
}

func corsEnvironmentPolicies(cfg *Config) []corsEnvironmentPolicy {
	if cfg == nil {
		return nil
	}
	return []corsEnvironmentPolicy{
		{name: corsJSONEnvName, value: cfg.CORSJSON, parser: parseJSONCORS},
		{name: corsYAMLEnvName, value: cfg.CORSYAML, parser: parseYAMLCORS},
	}
}

func loadEnvironmentCORS(policy corsEnvironmentPolicy) (*middleware.CORSConfig, bool) {
	if strings.TrimSpace(policy.value) == "" {
		return nil, false
	}

	_, fileCfg, err := policy.parser(policy.name, []byte(policy.value))
	if err != nil {
		logging.Warn(msgCORSInvalid, "variable", policy.name)
		return middleware.DefaultCORSConfig(), true
	}
	return buildCORSConfig(fileCfg), true
}
