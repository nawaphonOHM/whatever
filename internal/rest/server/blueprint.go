package server

import (
	"fmt"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

func loadConfig() (*Config, error) {
	cfg, err := intcfg.Load[Config]()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load server config: %w",
			err,
		)
	}
	return cfg, nil
}

func registerBlueprint(
	srv *Server,
	cfg *Config,
	bp *contracts.BluePrint,
) error {
	srv.SetupMiddlewares(buildCORSConfig(bp.Meta()))
	return srv.RegisterRoutesWithVersion(bp.Apis(), cfg.AppVersion)
}

func newServerFromConfig(
	cfg *Config,
	bp *contracts.BluePrint,
) (*Server, error) {
	srv, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if err := registerBlueprint(srv, cfg, bp); err != nil {
		return nil, err
	}
	return srv, nil
}

// NewFromBluePrint loads config and prepares a Server from a BluePrint.
func NewFromBluePrint(
	bluePrint *contracts.BluePrint,
) (*Server, error) {
	if bluePrint == nil {
		return nil, ErrNilBluePrint
	}
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return newServerFromConfig(cfg, bluePrint)
}

// NewFromRegistrations loads config and prepares a Server.
func NewFromRegistrations(
	registrations []*contracts.RRestAPIRegistration,
) (*Server, error) {
	bp := contracts.NewBluePrint().WithAPIs(registrations...)
	return NewFromBluePrint(bp)
}
