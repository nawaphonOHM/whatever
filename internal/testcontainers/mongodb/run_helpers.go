package mongodb

import (
	"maps"

	"github.com/testcontainers/testcontainers-go"
	tcmongodb "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// authCustomizers returns username and password customizers when configured.
func authCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	var list []testcontainers.ContainerCustomizer
	if o.Username != "" {
		list = append(list, tcmongodb.WithUsername(o.Username))
	}
	if o.Password != "" {
		list = append(list, tcmongodb.WithPassword(o.Password))
	}
	return list
}

// buildEnvMap creates the merged environment map including database initialization.
func buildEnvMap(o *Options) map[string]string {
	env := make(map[string]string)
	maps.Copy(env, o.Env)
	if o.Database != "" {
		env["MONGO_INITDB_DATABASE"] = o.Database
	}
	return env
}

// envCustomizers returns an environment customizer if environment variables exist.
func envCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	env := buildEnvMap(o)
	if len(env) == 0 {
		return nil
	}
	return []testcontainers.ContainerCustomizer{testcontainers.WithEnv(env)}
}

// replicaSetCustomizers returns a replica set customizer if replica set is configured.
func replicaSetCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	if o.ReplicaSet == "" {
		return nil
	}
	return []testcontainers.ContainerCustomizer{tcmongodb.WithReplicaSet(o.ReplicaSet)}
}

// buildCustomizers converts Options into a slice of testcontainers customizers.
func buildCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	var list []testcontainers.ContainerCustomizer
	list = append(list, authCustomizers(o)...)
	list = append(list, replicaSetCustomizers(o)...)
	list = append(list, envCustomizers(o)...)
	return append(list, o.ContainerOptions...)
}
