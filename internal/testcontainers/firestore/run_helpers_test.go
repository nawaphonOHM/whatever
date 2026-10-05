package firestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

const expectedCustomizers = 4

func TestBuildEnvMap_Empty(t *testing.T) {
	opts := &Options{
		Env: make(map[string]string),
	}
	assert.Nil(t, buildEnvMap(opts))
	assert.Nil(t, envCustomizers(opts))
}

func TestBuildEnvMap_WithValues(t *testing.T) {
	opts := &Options{
		Env: map[string]string{
			"K1": "V1",
			"K2": "V2",
		},
	}
	env := buildEnvMap(opts)
	require.NotNil(t, env)
	assert.Equal(t, "V1", env["K1"])
	assert.Equal(t, "V2", env["K2"])

	cust := envCustomizers(opts)
	require.Len(t, cust, 1)
}

func TestProjectIDCustomizers(t *testing.T) {
	optsEmpty := &Options{ProjectID: ""}
	assert.Nil(t, projectIDCustomizers(optsEmpty))

	optsSet := &Options{ProjectID: "proj-123"}
	assert.Len(t, projectIDCustomizers(optsSet), 1)
}

func TestDatastoreModeCustomizers(t *testing.T) {
	optsFalse := &Options{DatastoreMode: false}
	assert.Nil(t, datastoreModeCustomizers(optsFalse))

	optsTrue := &Options{DatastoreMode: true}
	assert.Len(t, datastoreModeCustomizers(optsTrue), 1)
}

func TestBuildCustomizers(t *testing.T) {
	c1 := &dummyCustomizer{}
	opts := &Options{
		ProjectID:        "proj-abc",
		DatastoreMode:    true,
		Env:              map[string]string{"ENV_A": "VAL_A"},
		ContainerOptions: []testcontainers.ContainerCustomizer{c1},
	}
	customizers := buildCustomizers(opts)
	assert.Len(t, customizers, expectedCustomizers) // projectID, datastoreMode, env, and 1 containerOption
}
