package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	otelcfg "github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	otelmw "github.com/nawaphonOHM/whatever/internal/opentelemetry/middleware"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/internal/rest/middleware"
	intprob "github.com/nawaphonOHM/whatever/internal/rest/problem"
	"github.com/nawaphonOHM/whatever/pkg/logger"
)

const (
	codeNotFound         = "NOT_FOUND"
	msgNotFound          = "Route not found"
	codeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	msgMethodNotAllowed  = "Method not allowed"
)

// defaultSkipPaths returns health paths skipped by request logger and tracer.
func defaultSkipPaths(extra ...string) []string {
	skip := []string{ReservedHealthPath, ReservedReadyPath}
	if len(extra) > 0 {
		skip = append(skip, extra...)
	}
	return skip
}

// resolveOTelConfig prepares OpenTelemetry config with default health skip paths.
func (s *Server) resolveOTelConfig(extraSkipPaths ...string) *otelcfg.Config {
	var cfg otelcfg.Config
	if s.Config != nil && s.Config.OTel != nil {
		cfg = *s.Config.OTel
	} else {
		cfg = *otelcfg.DefaultConfig()
	}
	cfg.SkipPaths = defaultSkipPaths(append(cfg.SkipPaths, extraSkipPaths...)...)
	return &cfg
}

// attachErrorHandlers registers RFC 9457 404/405 handlers.
func attachErrorHandlers(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		resp := intprob.New(http.StatusNotFound, codeNotFound, msgNotFound)
		resp.Write(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		resp := intprob.New(
			http.StatusMethodNotAllowed,
			codeMethodNotAllowed,
			msgMethodNotAllowed,
		)
		resp.Write(c)
	})
}

// SetupMiddlewares attaches middlewares with custom CORS configuration.
func (s *Server) SetupMiddlewares(
	corsCfg *middleware.CORSConfig,
	loggerSkipPaths ...string,
) {
	s.Engine.Use(otelmw.Middleware(s.resolveOTelConfig(loggerSkipPaths...)))
	s.Engine.Use(middleware.RequestID())
	if s.Config == nil || s.Config.EnableAccessLog {
		skipPaths := defaultSkipPaths(loggerSkipPaths...)
		s.Engine.Use(logger.WithConfig(&logger.Config{SkipPaths: skipPaths}))
	}
	s.Engine.Use(middleware.Recovery(), middleware.CORS(corsCfg))
	attachErrorHandlers(s.Engine)
}

// SetupDefaultMiddlewares attaches recommended middlewares with default CORS.
func (s *Server) SetupDefaultMiddlewares(
	loggerSkipPaths ...string,
) {
	s.SetupMiddlewares(middleware.DefaultCORSConfig(), loggerSkipPaths...)
}

// RegisterRoutes mounts registrations without version metadata.
func (s *Server) RegisterRoutes(
	registrations []*contracts.RRestAPIRegistration,
) error {
	return RegisterRoutesWithVersion(s.Engine, registrations, "")
}

// RegisterRoutesWithVersion mounts registrations with version.
func (s *Server) RegisterRoutesWithVersion(
	registrations []*contracts.RRestAPIRegistration,
	version string,
) error {
	return RegisterRoutesWithVersion(
		s.Engine,
		registrations,
		version,
	)
}
