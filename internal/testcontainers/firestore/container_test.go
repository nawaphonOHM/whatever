package firestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainer_Defaults(t *testing.T) {
	c := NewContainer(nil, nil)
	require.NotNil(t, c)
	assert.Equal(t, DefaultProjectID, c.ProjectID())
	assert.False(t, c.DatastoreMode())
	assert.Nil(t, c.RawContainer())
}

func TestNewContainer_WithOptions(t *testing.T) {
	opts := &Options{
		ProjectID:     "custom-project-id",
		DatastoreMode: true,
	}
	c := NewContainer(nil, opts)
	require.NotNil(t, c)
	assert.Equal(t, "custom-project-id", c.ProjectID())
	assert.True(t, c.DatastoreMode())
	assert.Nil(t, c.RawContainer())
}
