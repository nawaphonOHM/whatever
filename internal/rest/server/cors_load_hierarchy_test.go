package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCORSConfig_MetaOverridesServerConfig(t *testing.T) {
	metaPath := writeTempConfigFile(t, "meta_cors.json", `{"allow-origins": ["https://meta.example.com"]}`)
	envPath := writeTempConfigFile(t, "env_cors.json", `{"allow-origins": ["https://env.example.com"]}`)

	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta().WithCorsPath(metaPath))
	serverCfg := &Config{}
	serverCfg.CORSPath = envPath

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://meta.example.com"}, cfg.AllowOrigins)
}

func TestLoadCORSConfig_FallbackToServerConfig(t *testing.T) {
	envPath := writeTempConfigFile(t, "env_cors.json", `{"allow-origins": ["https://env.example.com"]}`)

	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta())
	serverCfg := &Config{}
	serverCfg.CORSPath = envPath

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://env.example.com"}, cfg.AllowOrigins)
}

func TestLoadCORSConfig_DefaultDisabledWhenNeither(t *testing.T) {
	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta())
	serverCfg := &Config{}

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Nil(t, cfg.AllowOrigins)
}
