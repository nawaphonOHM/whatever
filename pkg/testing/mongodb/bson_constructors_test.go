package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testFirstElem     = "first"
	expectedElemCount = 3
)

// TestNewE verifies construction of BSON elements with various key and value types.
func TestNewE(t *testing.T) {
	t.Run("string and nil values", func(t *testing.T) {
		elem1 := mongodb.NewE("username", "john_doe")
		assert.Equal(t, "username", elem1.Key)
		assert.Equal(t, "john_doe", elem1.Value)

		elem2 := mongodb.NewE("deleted_at", nil)
		assert.Equal(t, "deleted_at", elem2.Key)
		assert.Nil(t, elem2.Value)
	})

	t.Run("empty key and nested doc", func(t *testing.T) {
		elem3 := mongodb.NewE("", "empty-key-val")
		assert.Equal(t, "", elem3.Key)
		assert.Equal(t, "empty-key-val", elem3.Value)

		nestedDoc := mongodb.NewD(mongodb.NewE("status", "active"))
		elem4 := mongodb.NewE("filter", nestedDoc)
		assert.Equal(t, "filter", elem4.Key)
		assert.Equal(t, nestedDoc, elem4.Value)
	})
}

// TestNewD verifies construction of ordered BSON documents and key ordering.
func TestNewD(t *testing.T) {
	t.Run("empty elements", func(t *testing.T) {
		doc := mongodb.NewD()
		assert.NotNil(t, doc)
		assert.Empty(t, doc)
		assert.Len(t, doc, 0)
	})

	t.Run("multiple elements preserve order", func(t *testing.T) {
		e0 := mongodb.NewE(testFirstElem, 0)
		e1 := mongodb.NewE("second", 1)
		e2 := mongodb.NewE("third", 2)

		doc := mongodb.NewD(e0, e1, e2)
		require.NotNil(t, doc)
		require.Len(t, doc, expectedElemCount)
		assert.Equal(t, testFirstElem, doc[0].Key)
		assert.Equal(t, 0, doc[0].Value)
		assert.Equal(t, "second", doc[1].Key)
		assert.Equal(t, 1, doc[1].Value)
		assert.Equal(t, "third", doc[2].Key)
		assert.Equal(t, 2, doc[2].Value)
	})
}

// TestNewM verifies construction and mutability of unordered BSON maps.
func TestNewM(t *testing.T) {
	m := mongodb.NewM()
	require.NotNil(t, m)
	assert.Empty(t, m)

	m["key1"] = "val1"
	m["key2"] = 2
	assert.Equal(t, "val1", m["key1"])
	assert.Equal(t, 2, m["key2"])
}

// TestNewA verifies construction of BSON arrays and element ordering.
func TestNewA(t *testing.T) {
	t.Run("empty values", func(t *testing.T) {
		arr := mongodb.NewA()
		assert.NotNil(t, arr)
		assert.Empty(t, arr)
		assert.Len(t, arr, 0)
	})

	t.Run("elements preserve order", func(t *testing.T) {
		subDoc := mongodb.NewD(mongodb.NewE("sub", "val"))
		arr := mongodb.NewA(testFirstElem, 1, subDoc)
		require.NotNil(t, arr)
		require.Len(t, arr, expectedElemCount)
		assert.Equal(t, testFirstElem, arr[0])
		assert.Equal(t, 1, arr[1])
		assert.Equal(t, subDoc, arr[2])
	})
}
