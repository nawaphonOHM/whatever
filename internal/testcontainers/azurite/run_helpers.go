package azurite

import (
	"maps"

	"github.com/testcontainers/testcontainers-go"
)

func buildCustomizers(o *Options) []testcontainers.ContainerCustomizer {
	var customizers []testcontainers.ContainerCustomizer
	if len(o.Env) > 0 {
		env := make(map[string]string, len(o.Env))
		maps.Copy(env, o.Env)
		customizers = append(customizers, testcontainers.WithEnv(env))
	}
	return append(customizers, o.ContainerOptions...)
}
