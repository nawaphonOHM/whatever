package mongodb

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

const (
	testHost     = "localhost"
	testPort     = 27017
	testDatabase = "app"
)

func TestNewContainer_Defaults(t *testing.T) {
	c := NewContainer(nil, nil)
	testNilGettersStrings(t, c)
	if c.RawContainer() != nil {
		t.Error("expected nil raw container")
	}
}

func verifyContainerOptions(t *testing.T, c *Container) {
	if c.Username() != "testuser" || c.Password() != "testpass" {
		t.Error("credentials mismatch")
	}
}

func verifyContainerSettings(t *testing.T, c *Container) {
	if c.Database() != "testdb" || c.ReplicaSet() != "rs0" {
		t.Error("database or replica set mismatch")
	}
}

func TestNewContainer_WithOptions(t *testing.T) {
	opts := &Options{
		Username:   "testuser",
		Password:   "testpass",
		Database:   "testdb",
		ReplicaSet: "rs0",
	}
	c := NewContainer(nil, opts)
	verifyContainerOptions(t, c)
	verifyContainerSettings(t, c)
}

func verifyConfigHostPort(t *testing.T, cfg *config.Config) {
	if cfg.Host != testHost {
		t.Errorf("host mismatch: got %s", cfg.Host)
	}
	if cfg.Port != testPort {
		t.Errorf("port mismatch: got %d", cfg.Port)
	}
}

func verifyConfigCredentials(t *testing.T, cfg *config.Config, user, pass string) {
	if cfg.Username != user || cfg.Password != pass {
		t.Errorf("user/pass mismatch: got user=%s pass=%s", cfg.Username, cfg.Password)
	}
}

func TestPopulateConfig_WithAuth(t *testing.T) {
	c := &Container{username: "admin", password: "secret", defaultDatabase: testDatabase}
	cfg := c.populateConfig(testHost, testPort)
	verifyConfigHostPort(t, cfg)
	verifyConfigCredentials(t, cfg, "admin", "secret")
	if cfg.AuthSource != "admin" || cfg.Database != testDatabase {
		t.Errorf("auth/database mismatch: auth=%s db=%s", cfg.AuthSource, cfg.Database)
	}
}

func TestPopulateConfig_WithoutAuth(t *testing.T) {
	c := &Container{defaultDatabase: testDatabase}
	cfg := c.populateConfig(testHost, testPort)
	verifyConfigHostPort(t, cfg)
	verifyConfigCredentials(t, cfg, "", "")
	if cfg.AuthSource != "" || cfg.Database != testDatabase {
		t.Errorf("auth/database mismatch: auth=%s db=%s", cfg.AuthSource, cfg.Database)
	}
}
