package server

import (
	"net/http"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/middleware"
)

// CORSFileConfig represents the deserialized cors.schema.json structure.
type CORSFileConfig struct {
	AllowMethods   []string `json:"allow-methods" yaml:"allow-methods"`
	RequestMethods []string `json:"request-methods" yaml:"request-methods"`
	AllowOrigins   []string `json:"allow-origins" yaml:"allow-origins"`
	AllowHeaders   []string `json:"allow-headers" yaml:"allow-headers"`
	RequestHeaders []string `json:"request-headers" yaml:"request-headers"`
}

var validMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
	http.MethodConnect: {},
	http.MethodTrace:   {},
}

func filterValidMethods(methods []string) []string {
	out := make([]string, 0, len(methods))
	for _, m := range methods {
		trimmed := strings.ToUpper(strings.TrimSpace(m))
		if _, ok := validMethods[trimmed]; ok {
			out = append(out, trimmed)
		}
	}
	return out
}

func applyAllowFields(cfg *middleware.CORSConfig, fileCfg *CORSFileConfig) {
	if len(fileCfg.AllowOrigins) > 0 {
		cfg.AllowOrigins = fileCfg.AllowOrigins
	}
	if len(fileCfg.AllowHeaders) > 0 {
		cfg.AllowHeaders = fileCfg.AllowHeaders
	}
}

func buildCORSConfig(fileCfg *CORSFileConfig) *middleware.CORSConfig {
	cfg := middleware.DefaultCORSConfig()
	if fileCfg == nil {
		return cfg
	}
	applyAllowFields(cfg, fileCfg)
	if methods := filterValidMethods(fileCfg.AllowMethods); len(methods) > 0 {
		cfg.AllowMethods = methods
	}
	return cfg
}
