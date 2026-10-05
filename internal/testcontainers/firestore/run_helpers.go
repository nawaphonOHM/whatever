package firestore

import (
	"maps"

	"github.com/testcontainers/testcontainers-go"
	tcfirestore "github.com/testcontainers/testcontainers-go/modules/gcloud/firestore"
)

// projectIDCustomizers returns a project ID customizer when configured.
func projectIDCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	if o.ProjectID == "" {
		return nil
	}
	return []testcontainers.ContainerCustomizer{tcfirestore.WithProjectID(o.ProjectID)}
}

// datastoreModeCustomizers returns a datastore mode customizer when configured.
func datastoreModeCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	if !o.DatastoreMode {
		return nil
	}
	return []testcontainers.ContainerCustomizer{tcfirestore.WithDatastoreMode()}
}

// buildEnvMap creates the merged environment map.
func buildEnvMap(o *Options) map[string]string {
	if len(o.Env) == 0 {
		return nil
	}
	env := make(map[string]string, len(o.Env))
	maps.Copy(env, o.Env)
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

// buildCustomizers converts Options into a slice of testcontainers customizers.
func buildCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	var list []testcontainers.ContainerCustomizer
	list = append(list, projectIDCustomizers(o)...)
	list = append(list, datastoreModeCustomizers(o)...)
	list = append(list, envCustomizers(o)...)
	return append(list, o.ContainerOptions...)
}
