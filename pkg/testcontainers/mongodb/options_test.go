package mongodb_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/pkg/testcontainers/mongodb"
)

func verifyOptionCredentials(t *testing.T, o *mongodb.Options) {
	if o.Username != "admin" || o.Password != "secret" {
		t.Error("username or password mismatch")
	}
}

func verifyOptionImageAndDb(t *testing.T, o *mongodb.Options) {
	if o.Image != "mongo:7" || o.Database != "app" {
		t.Error("image or database mismatch")
	}
}

func verifyOptionReplicaAndEnv(t *testing.T, o *mongodb.Options) {
	if o.ReplicaSet != "rs0" || o.Env["K1"] != "V1" {
		t.Error("replica set or env mismatch")
	}
}

func TestOptionsBuilders(t *testing.T) {
	opts := mongodb.NewOptions(
		nil,
		mongodb.WithImage("mongo:7"),
		mongodb.WithUsername("admin"),
		mongodb.WithPassword("secret"),
		mongodb.WithDatabase("app"),
		mongodb.WithReplicaSet("rs0"),
		mongodb.WithEnv("K1", "V1"),
		mongodb.WithContainerOptions(nil),
	)
	verifyOptionCredentials(t, opts)
	verifyOptionImageAndDb(t, opts)
	verifyOptionReplicaAndEnv(t, opts)
}

func TestDefaultOptions(t *testing.T) {
	def := mongodb.DefaultOptions()
	if def.Image != mongodb.DefaultImage {
		t.Errorf("unexpected image: %s", def.Image)
	}
	if def.Env == nil {
		t.Error("expected non-nil Env map")
	}
}
