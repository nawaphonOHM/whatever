package server

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/internal/rest/middleware"
)

// mapHTTPMethods converts HTTP methods into uppercase strings.
func mapHTTPMethods(methods []contracts.HTTPMethod) []string {
	out := make([]string, 0, len(methods))
	for _, method := range methods {
		if method.IsValid() {
			out = append(out, method.GetName())
		}
	}
	return out
}

func corsFromMeta(meta *contracts.Meta) *contracts.CorsSetting {
	if meta == nil {
		return nil
	}
	return meta.Cors()
}

func applyCORSOrigins(
	cfg *middleware.CORSConfig,
	cors *contracts.CorsSetting,
) {
	origins := cors.AllowOrigin()
	if len(origins) > 0 {
		cfg.AllowOrigins = origins
	}
}

func applyCORSMethods(
	cfg *middleware.CORSConfig,
	cors *contracts.CorsSetting,
) {
	methods := mapHTTPMethods(cors.AllowHTTPMethods())
	if len(methods) > 0 {
		cfg.AllowMethods = methods
	}
}

// buildCORSConfig derives CORSConfig from Meta.Cors() or returns default.
func buildCORSConfig(meta *contracts.Meta) middleware.CORSConfig {
	cfg := middleware.DefaultCORSConfig()
	cors := corsFromMeta(meta)
	if cors == nil {
		return cfg
	}
	applyCORSOrigins(&cfg, cors)
	applyCORSMethods(&cfg, cors)
	return cfg
}
