package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testFieldKey = "field"
	testScoreVal = int32(100)
)

// TestNestedBSONConstructionAndSerialization verifies serialization of nested empty documents,
// arrays, nil values, and round-trip conversion with Marshal and Unmarshal.
func TestNestedBSONConstructionAndSerialization(t *testing.T) {
	t.Run("nested empty document", func(t *testing.T) {
		doc := mongodb.NewD(mongodb.NewE("nested", mongodb.NewD()))
		encoded, err := mongodb.Marshal(doc)
		require.NoError(t, err)

		driverExpected, err := bson.Marshal(bson.D{{Key: "nested", Value: bson.D{}}})
		require.NoError(t, err)
		assert.Equal(t, driverExpected, encoded)

		var decoded mongodb.M
		err = mongodb.Unmarshal(encoded, &decoded)
		require.NoError(t, err)
		assert.Contains(t, decoded, "nested")
	})

	t.Run("nested empty array", func(t *testing.T) {
		doc := mongodb.NewD(mongodb.NewE("items", mongodb.NewA()))
		encoded, err := mongodb.Marshal(doc)
		require.NoError(t, err)

		driverExpected, err := bson.Marshal(bson.D{{Key: "items", Value: bson.A{}}})
		require.NoError(t, err)
		assert.Equal(t, driverExpected, encoded)

		var decoded mongodb.M
		err = mongodb.Unmarshal(encoded, &decoded)
		require.NoError(t, err)
		assert.Contains(t, decoded, "items")
	})

	t.Run("nil element value", func(t *testing.T) {
		doc := mongodb.NewD(mongodb.NewE(testFieldKey, nil))
		encoded, err := mongodb.Marshal(doc)
		require.NoError(t, err)

		driverExpected, err := bson.Marshal(bson.D{{Key: testFieldKey, Value: nil}})
		require.NoError(t, err)
		assert.Equal(t, driverExpected, encoded)

		var decoded mongodb.M
		err = mongodb.Unmarshal(encoded, &decoded)
		require.NoError(t, err)
		assert.Contains(t, decoded, testFieldKey)
		assert.Nil(t, decoded[testFieldKey])
	})

	t.Run("complex construction roundtrip with testing serialization", func(t *testing.T) {
		complexDoc := mongodb.NewD(
			mongodb.NewE("name", "test-item"),
			mongodb.NewE("tags", mongodb.NewA("tag1", "tag2")),
			mongodb.NewE("metadata", mongodb.NewD(mongodb.NewE("score", testScoreVal))),
			mongodb.NewE("subdoc", mongodb.NewD(mongodb.NewE("active", true))),
		)

		encoded, err := mongodb.Marshal(complexDoc)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)

		var decodedDoc mongodb.D
		err = mongodb.Unmarshal(encoded, &decodedDoc)
		require.NoError(t, err)
		assert.Equal(t, complexDoc, decodedDoc)

		var driverDoc bson.D
		err = bson.Unmarshal(encoded, &driverDoc)
		require.NoError(t, err)
		assert.Equal(t, complexDoc, driverDoc)

		var decodedMap mongodb.M
		err = mongodb.Unmarshal(encoded, &decodedMap)
		require.NoError(t, err)
		assert.Equal(t, "test-item", decodedMap["name"])
	})
}
