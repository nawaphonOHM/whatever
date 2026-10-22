package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// serverAPIGroup returns the configured API group of the server.
func (s *Server) serverAPIGroup() (contracts.APIGroup, error) {
	cfg := resolveConfig(s.Config)
	return resolveConfigAPIGroup(cfg.APIGroup)
}

// errOverrideGroupMismatch formats an error when the override group conflicts with the server group.
func errOverrideGroupMismatch(overrideGrp, serverGrp contracts.APIGroup) error {
	return fmt.Errorf(
		"%w: config override group %q conflicts with server group %q",
		ErrAPIGroupMismatch,
		overrideGrp,
		serverGrp,
	)
}

// matchOverrideGroup verifies resolved override group against server group.
func matchOverrideGroup(grp, serverGrp contracts.APIGroup) error {
	if grp != serverGrp {
		return errOverrideGroupMismatch(grp, serverGrp)
	}
	return nil
}

// checkOverrideGroup validates that the override API group matches the server's API group.
func checkOverrideGroup(overrideGrp string, serverGrp contracts.APIGroup) error {
	if overrideGrp == "" {
		return nil
	}
	grp, err := resolveConfigAPIGroup(overrideGrp)
	if err != nil {
		return err
	}
	return matchOverrideGroup(grp, serverGrp)
}

// applyOverrideConfig validates override group and returns the effective configuration.
func applyOverrideConfig(cfg *Config, serverGrp contracts.APIGroup) (*Config, error) {
	if err := checkOverrideGroup(cfg.APIGroup, serverGrp); err != nil {
		return nil, err
	}
	c := *cfg
	c.APIGroup = string(serverGrp)
	return &c, nil
}

// resolveRegistrationConfig validates override group and returns the effective server configuration.
func (s *Server) resolveRegistrationConfig(cfg *Config) (*Config, error) {
	serverGrp, err := s.serverAPIGroup()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		c := *resolveConfig(s.Config)
		c.APIGroup = string(serverGrp)
		return &c, nil
	}
	return applyOverrideConfig(cfg, serverGrp)
}

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
	effectiveCfg, err := s.resolveRegistrationConfig(cfg)
	if err != nil {
		return err
	}
	return RegisterRoutesWithConfig(
		s.Engine,
		registrations,
		effectiveCfg,
	)
}
