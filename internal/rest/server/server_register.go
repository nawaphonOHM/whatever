package server

import (
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// RegisterRoutes mounts registrations without version metadata using server config.
func (s *Server) RegisterRoutes(
	registrations []*contracts.RRestAPIRegistration,
) error {
	return s.RegisterRoutesWithConfig(registrations, s.Config)
}

// RegisterRoutesWithVersion mounts registrations with version using server config.
func (s *Server) RegisterRoutesWithVersion(
	registrations []*contracts.RRestAPIRegistration,
	version string,
) error {
	cfg := resolveConfig(s.Config)
	cfgCopy := *cfg
	cfgCopy.AppVersion = version
	return s.RegisterRoutesWithConfig(
		registrations,
		&cfgCopy,
	)
}

// RegisterRoutesWithConfig mounts registrations using the provided configuration.
func (s *Server) RegisterRoutesWithConfig(
	registrations []*contracts.RRestAPIRegistration,
	cfg *Config,
) error {
	effectiveCfg := cfg
	if effectiveCfg == nil {
		effectiveCfg = s.Config
	}
	return RegisterRoutesWithConfig(
		s.Engine,
		registrations,
		effectiveCfg,
	)
}
