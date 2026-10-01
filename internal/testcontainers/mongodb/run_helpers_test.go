package mongodb

import (
	"testing"
)

const expectedCustomizerCount = 4

func TestBuildEnvMap_Empty(t *testing.T) {
	o := &Options{}
	env := buildEnvMap(o)
	if len(env) != 0 {
		t.Errorf("expected empty env map, got %v", env)
	}
}

func TestBuildEnvMap_WithDbAndCustom(t *testing.T) {
	o := &Options{
		Database: "my_db",
		Env:      map[string]string{"CUSTOM_KEY": "CUSTOM_VAL"},
	}
	env := buildEnvMap(o)
	if env["MONGO_INITDB_DATABASE"] != "my_db" || env["CUSTOM_KEY"] != "CUSTOM_VAL" {
		t.Errorf("unexpected env map: %v", env)
	}
}

func TestAuthCustomizers(t *testing.T) {
	oEmpty := &Options{}
	if len(authCustomizers(oEmpty)) != 0 {
		t.Error("expected 0 auth customizers for empty options")
	}
	oAuth := &Options{Username: "user", Password: "pwd"}
	if len(authCustomizers(oAuth)) != 2 {
		t.Errorf("expected 2 auth customizers, got %d", len(authCustomizers(oAuth)))
	}
}

func TestReplicaSetCustomizer(t *testing.T) {
	if len(replicaSetCustomizers(&Options{})) != 0 {
		t.Error("expected 0 replica set customizers for empty options")
	}
	if len(replicaSetCustomizers(&Options{ReplicaSet: "rs0"})) != 1 {
		t.Error("expected 1 replica set customizer")
	}
}

func TestEnvCustomizers(t *testing.T) {
	if len(envCustomizers(&Options{})) != 0 {
		t.Error("expected 0 env customizers for empty options")
	}
	if len(envCustomizers(&Options{Database: "db"})) != 1 {
		t.Error("expected 1 env customizer for configured db")
	}
}

func TestBuildCustomizers(t *testing.T) {
	o := &Options{
		Username:   "admin",
		Password:   "secret",
		ReplicaSet: "rs0",
		Database:   "db",
	}
	customizers := buildCustomizers(o)
	if len(customizers) != expectedCustomizerCount {
		t.Errorf("expected %d customizers, got %d", expectedCustomizerCount, len(customizers))
	}
}
