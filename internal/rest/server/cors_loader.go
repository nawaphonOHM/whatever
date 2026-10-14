package server

import (
	"fmt"
	"os"

	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/v2/internal/rest/middleware"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
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

func blueprintCorsPath(bp *contracts.BluePrint) string {
	if bp == nil || bp.Meta() == nil {
		return ""
	}
	return bp.Meta().CorsPath()
}

func serverConfigCorsPath(cfg *Config) string {
	if cfg == nil {
		return ""
	}
	return cfg.CORSPath
}

func resolveCORSPath(cfg *Config, bp *contracts.BluePrint) string {
	if path := blueprintCorsPath(bp); path != "" {
		return path
	}
	return serverConfigCorsPath(cfg)
}

func loadCORSConfig(cfg *Config, bp *contracts.BluePrint) (*middleware.CORSConfig, error) {
	path := resolveCORSPath(cfg, bp)
	corsCfg, err := LoadAndValidateCORS(path)
	if err != nil {
		logging.Error("failed to load CORS configuration", "error", err.Error())
		return nil, err
	}
	return corsCfg, nil
}
