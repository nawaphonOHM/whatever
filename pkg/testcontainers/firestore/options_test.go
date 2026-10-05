package firestore_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
	"github.com/stretchr/testify/assert"
)

func verifyOptionImageAndProject(t *testing.T, o *firestore.Options) {
	if o.Image != "custom-image:1.0" || o.ProjectID != "test-proj" {
		t.Error("image or projectID mismatch")
	}
}

func verifyOptionDatastoreAndEnv(t *testing.T, o *firestore.Options) {
	if !o.DatastoreMode || o.Env["K1"] != "V1" {
		t.Error("datastore mode or env mismatch")
	}
}

func TestOptionsBuilders(t *testing.T) {
	opts := firestore.NewOptions(
		nil,
		firestore.WithImage("custom-image:1.0"),
		firestore.WithProjectID("test-proj"),
		firestore.WithDatastoreMode(true),
		firestore.WithEnv("K1", "V1"),
		firestore.WithContainerOptions(nil),
	)
	verifyOptionImageAndProject(t, opts)
	verifyOptionDatastoreAndEnv(t, opts)
}

func TestDefaultOptions(t *testing.T) {
	def := firestore.DefaultOptions()
	assert.Equal(t, firestore.DefaultImage, def.Image)
	assert.Equal(t, firestore.DefaultProjectID, def.ProjectID)
	assert.False(t, def.DatastoreMode)
	assert.NotNil(t, def.Env)
}
