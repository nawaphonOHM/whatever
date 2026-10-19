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

func TestMarshalErrors(t *testing.T) {
	// Channels cannot be marshaled to BSON
	ch := make(chan int)
	data, err := mongodb.Marshal(ch)
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestUnmarshalErrors(t *testing.T) {
	t.Run("invalid corrupted bytes", func(t *testing.T) {
		var person testPerson
		err := mongodb.Unmarshal([]byte("invalid-corrupted-bytes"), &person)
		assert.Error(t, err)
	})

	t.Run("empty bytes", func(t *testing.T) {
		var person testPerson
		err := mongodb.Unmarshal(make([]byte, 0), &person)
		assert.Error(t, err)
	})

	t.Run("nil bytes", func(t *testing.T) {
		var person testPerson
		err := mongodb.Unmarshal(nil, &person)
		assert.Error(t, err)
	})

	t.Run("non pointer destination", func(t *testing.T) {
		data, err := mongodb.Marshal(testPerson{Name: "Bob"})
		require.NoError(t, err)

		var person testPerson
		err = mongodb.Unmarshal(data, person)
		assert.Error(t, err)
	})
}
