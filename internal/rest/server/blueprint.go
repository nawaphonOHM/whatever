package server

import (
	"fmt"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/pkg/logging"
)

func loadConfig() (*Config, error) {
	cfg, err := intcfg.Load[Config]()
	if err != nil {
		logging.Error("failed to load REST server config", "error", err.Error())
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
		logging.Error("failed to create REST server", "error", err.Error())
		return nil, err
	}
	if err := registerBlueprint(srv, cfg, bp); err != nil {
		logging.Error("failed to register REST routes", "error", err.Error())
		return nil, err
	}
	return srv, nil
}

// NewFromBluePrint loads config and prepares a Server from a BluePrint.
func NewFromBluePrint(
	bluePrint *contracts.BluePrint,
) (*Server, error) {
	if bluePrint == nil {
		logging.Error("failed to initialize REST server: blueprint is nil")
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
