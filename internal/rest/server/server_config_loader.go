package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	intcfg "github.com/nawaphonOHM/whatever/v3/internal/rest/config"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

func logConfigFallback(filenames []string) {
	if len(filenames) == 0 {
		logging.Info(
			"falling back to default env configuration files",
			"files",
			intcfg.DefaultEnvFiles(),
		)
	}
}

// LoadConfig loads server configuration with structured logging.
func LoadConfig(filenames ...string) (*Config, error) {
	decorated := callstack.DecorateFuncErr("server.LoadConfig", func() (*Config, error) {
		logging.Info("loading REST configuration")
		logConfigFallback(filenames)
		cfg, err := intcfg.Load[Config](filenames...)
		if err != nil {
			logging.Error("failed to load REST server config", "error", err)
			return nil, fmt.Errorf("failed to load server config: %w", err)
		}
		logging.Info("REST configuration loaded successfully")
		return cfg, nil
	})
	return decorated()
}
