package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/health"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

// mountvalidatedRoutes registers validated routes onto the engine.
func mountvalidatedRoutes(
	engine *gin.Engine,
	routes []*validatedRoute,
) {
	for _, route := range routes {
		handlers := buildHandlers(route.Middlewares, route.Handler)
		engine.Handle(route.Method, route.Path, handlers...)
	}
}

// mountHealthEndpoints registers reserved health probes.
func mountHealthEndpoints(engine *gin.Engine, version string) {
	hh := health.New(version)
	engine.GET(ReservedHealthPath, hh.Health)
	engine.GET(ReservedReadyPath, hh.Ready)
}

// RegisterRoutesWithConfig validates registrations against cfg and mounts routes and health endpoints.
func RegisterRoutesWithConfig(
	engine *gin.Engine,
	registrations []*contracts.RRestAPIRegistration,
	cfg *Config,
) error {
	decorated := callstack.DecorateErr("server.RegisterRoutesWithConfig", func() error {
		c := resolveConfig(cfg)
		validated, err := validateRegistrationsWithConfig(registrations, c)
		if err != nil {
			return err
		}
		logging.Info("registering REST routes",
			"route_count", len(validated),
			"version", c.AppVersion,
			"api_group", c.APIGroup,
		)
		mountHealthEndpoints(engine, c.AppVersion)
		logging.Info("registered reserved health endpoints",
			"paths", []string{ReservedHealthPath, ReservedReadyPath},
		)
		mountvalidatedRoutes(engine, validated)
		logging.Info("route registration complete", "total_routes", len(validated)+2)
		return nil
	})
	return decorated()
}

// RegisterRoutesWithVersion validates and mounts API registrations with standard API group.
func RegisterRoutesWithVersion(
	engine *gin.Engine,
	registrations []*contracts.RRestAPIRegistration,
	version string,
) error {
	cfg := DefaultConfig()
	cfg.AppVersion = version
	return RegisterRoutesWithConfig(engine, registrations, cfg)
}
