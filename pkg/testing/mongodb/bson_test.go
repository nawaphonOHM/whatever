package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testPersonAge = 30
	testRoleAdmin = "admin"
)

type testPerson struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Age   int    `bson:"age"`
}

func TestMarshalUnmarshalStruct(t *testing.T) {
	orig := testPerson{Name: "Alice", Email: "alice@example.com", Age: testPersonAge}

	data, err := mongodb.Marshal(orig)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var decoded testPerson
	err = mongodb.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, orig, decoded)
}

func TestMarshalUnmarshalMap(t *testing.T) {
	orig := map[string]any{"role": testRoleAdmin, "active": true}

	data, err := mongodb.Marshal(orig)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var decoded map[string]any
	err = mongodb.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, testRoleAdmin, decoded["role"])
	assert.Equal(t, true, decoded["active"])
}

func TestMarshalUnmarshalDriverInterop(t *testing.T) {
	doc := bson.D{{Key: "key", Value: "val"}}

	t.Run("pkg marshal to driver unmarshal", func(t *testing.T) {
		data, err := mongodb.Marshal(doc)
		require.NoError(t, err)

		var m bson.M
		require.NoError(t, bson.Unmarshal(data, &m))
		assert.Equal(t, "val", m["key"])
	})

	t.Run("driver marshal to pkg unmarshal", func(t *testing.T) {
		data, err := bson.Marshal(doc)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, mongodb.Unmarshal(data, &m))
		assert.Equal(t, "val", m["key"])
	})
}
