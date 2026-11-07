package mongodb_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

const (
	testNameKey   = "name"
	testNameVal   = "Gopher"
	testStatusKey = "status"
	testActiveVal = "active"
)

type objectIDDocument struct {
	ID mongodb.ObjectID `bson:"_id" json:"_id"`
}

// TestBSONSerializationCompatibility verifies round-trip BSON marshaling.
func TestBSONSerializationCompatibility(t *testing.T) {
	doc := mongodb.NewD(
		mongodb.NewE(testNameKey, testNameVal),
		mongodb.NewE(testStatusKey, testActiveVal),
	)

	bytes, err := bson.Marshal(doc)
	require.NoError(t, err)

	var unmarshaledMap mongodb.M
	err = bson.Unmarshal(bytes, &unmarshaledMap)
	require.NoError(t, err)

	assert.Equal(t, testNameVal, unmarshaledMap[testNameKey])
	assert.Equal(t, testActiveVal, unmarshaledMap[testStatusKey])
}

func TestEmptyNewDMarshalsAsEmptyDocument(t *testing.T) {
	doc := mongodb.NewD()
	require.NotNil(t, doc)
	encoded, err := bson.Marshal(doc)
	require.NoError(t, err)
	expected, err := bson.Marshal(bson.D{})
	require.NoError(t, err)
	assert.Equal(t, expected, encoded)
}

func TestNestedEmptyNewDMarshalsAsDocument(t *testing.T) {
	doc := mongodb.NewD(mongodb.NewE("nested", mongodb.NewD()))
	encoded, err := bson.Marshal(doc)
	require.NoError(t, err)
	expected, err := bson.Marshal(bson.D{{Key: "nested", Value: bson.D{}}})
	require.NoError(t, err)
	assert.Equal(t, expected, encoded)
}

func TestObjectIDSerializationRoundTrips(t *testing.T) {
	id := mongodb.NewObjectID()
	assertObjectIDBSONRoundTrip(t, id)
	assertObjectIDJSONRoundTrip(t, id)
}

func assertObjectIDBSONRoundTrip(t *testing.T, id mongodb.ObjectID) {
	encoded, err := bson.Marshal(mongodb.NewD(mongodb.NewE("_id", id)))
	require.NoError(t, err)

	var result objectIDDocument
	require.NoError(t, bson.Unmarshal(encoded, &result))
	assert.Equal(t, id, result.ID)
}

func assertObjectIDJSONRoundTrip(t *testing.T, id mongodb.ObjectID) {
	encoded, err := json.Marshal(objectIDDocument{ID: id})
	require.NoError(t, err)

	var result objectIDDocument
	require.NoError(t, json.Unmarshal(encoded, &result))
	assert.Equal(t, id, result.ID)
}

func TestEmptyNewAMarshalsAsArray(t *testing.T) {
	doc := mongodb.NewD(mongodb.NewE("items", mongodb.NewA()))
	require.NotNil(t, mongodb.NewA())
	encoded, err := bson.Marshal(doc)
	require.NoError(t, err)
	expected, err := bson.Marshal(bson.D{{Key: "items", Value: bson.A{}}})
	require.NoError(t, err)
	assert.Equal(t, expected, encoded)
}

// TestClientUsagePatterns demonstrates realistic client queries and pipelines.
func TestClientUsagePatterns(t *testing.T) {
	t.Run("filter serialization", func(t *testing.T) {
		filter := mongodb.NewD(
			mongodb.NewE(testStatusKey, testActiveVal),
			mongodb.NewE("tags", mongodb.NewA("news", "tech")),
		)
		filterBytes, err := bson.Marshal(filter)
		require.NoError(t, err)
		require.NotEmpty(t, filterBytes)
	})

	t.Run("pipeline serialization", func(t *testing.T) {
		pipeline := mongodb.NewA(
			mongodb.NewD(mongodb.NewE("$match", mongodb.NewD(mongodb.NewE("valid", true)))),
		)
		cmdDoc := mongodb.NewD(
			mongodb.NewE("aggregate", "items"),
			mongodb.NewE("pipeline", pipeline),
		)
		cmdBytes, err := bson.Marshal(cmdDoc)
		require.NoError(t, err)
		require.NotEmpty(t, cmdBytes)
	})
}
