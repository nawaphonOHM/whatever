package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/nawaphonOHM/whatever/internal/rest/middleware"
	"github.com/nawaphonOHM/whatever/pkg/logger"
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

const (
	codeNotFound         = "NOT_FOUND"
	msgNotFound          = "Route not found"
	codeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	msgMethodNotAllowed  = "Method not allowed"
)

// defaultSkipPaths returns health paths skipped by request logger.
func defaultSkipPaths(extra ...string) []string {
	skip := []string{ReservedHealthPath, ReservedReadyPath}
	if len(extra) > 0 {
		skip = append(skip, extra...)
	}
	return skip
}

// attachErrorHandlers registers RFC 9457 404/405 handlers.
func attachErrorHandlers(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		resp := rest.NotFound(codeNotFound, msgNotFound)
		resp.Write(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		resp := rest.Error(
			http.StatusMethodNotAllowed,
			codeMethodNotAllowed,
			msgMethodNotAllowed,
		)
		resp.Write(c)
	})
}

// SetupDefaultMiddlewares attaches recommended middlewares.
func (s *Server) SetupDefaultMiddlewares(
	loggerSkipPaths ...string,
) {
	s.Engine.HandleMethodNotAllowed = true
	skipPaths := defaultSkipPaths(loggerSkipPaths...)
	s.Engine.Use(
		middleware.RequestID(),
		logger.WithConfig(logger.Config{SkipPaths: skipPaths}),
		middleware.Recovery(),
		middleware.CORS(middleware.DefaultCORSConfig()),
	)
	attachErrorHandlers(s.Engine)
}

// NewFromRegistrations loads config and prepares a Server.
func NewFromRegistrations(
	registrations []*rest.RRestAPIRegistration,
) (*Server, error) {
	cfg, err := intcfg.Load[Config]()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load server config: %w",
			err,
		)
	}
	srv := New(cfg)
	srv.SetupDefaultMiddlewares()
	err = srv.RegisterRoutesWithVersion(
		registrations,
		cfg.AppVersion,
	)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

// RegisterRoutes mounts registrations without version metadata.
func (s *Server) RegisterRoutes(
	registrations []*rest.RRestAPIRegistration,
) error {
	return RegisterRoutesWithVersion(s.Engine, registrations, "")
}

// RegisterRoutesWithVersion mounts registrations with version.
func (s *Server) RegisterRoutesWithVersion(
	registrations []*rest.RRestAPIRegistration,
	version string,
) error {
	return RegisterRoutesWithVersion(
		s.Engine,
		registrations,
		version,
	)
}
