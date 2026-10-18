package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const testMagicInt = 123

func verifyDriverDoc(d bson.D) bson.D {
	return d
}

func verifyTestingDoc(d mongodb.D) mongodb.D {
	return d
}

func verifyDriverElem(e bson.E) bson.E {
	return e
}

func verifyTestingElem(e mongodb.E) mongodb.E {
	return e
}

func verifyDriverMap(m bson.M) bson.M {
	return m
}

func verifyTestingMap(m mongodb.M) mongodb.M {
	return m
}

func verifyDriverArr(a bson.A) bson.A {
	return a
}

func verifyTestingArr(a mongodb.A) mongodb.A {
	return a
}

// TestBSONTypeAliasesInteroperability verifies compile-time and runtime
// bi-directional assignment compatibility between pkg/testing/mongodb aliases and mongo-driver bson types.
func TestBSONTypeAliasesInteroperability(t *testing.T) {
	t.Run("element interoperability", func(t *testing.T) {
		driverElem := bson.E{Key: "k1", Value: "v1"}
		testingElem := verifyTestingElem(driverElem)
		assert.Equal(t, driverElem, testingElem)

		driverElem2 := verifyDriverElem(testingElem)
		assert.Equal(t, testingElem, driverElem2)
	})

	t.Run("document interoperability", func(t *testing.T) {
		driverDoc := bson.D{{Key: "k1", Value: "v1"}}
		testingDoc := verifyTestingDoc(driverDoc)
		assert.Equal(t, driverDoc, testingDoc)

		driverDoc2 := verifyDriverDoc(testingDoc)
		assert.Equal(t, testingDoc, driverDoc2)
	})

	t.Run("map interoperability", func(t *testing.T) {
		driverMap := bson.M{"key": "value"}
		testingMap := verifyTestingMap(driverMap)
		assert.Equal(t, driverMap, testingMap)

		driverMap2 := verifyDriverMap(testingMap)
		assert.Equal(t, testingMap, driverMap2)
	})

	t.Run("array interoperability", func(t *testing.T) {
		driverArr := bson.A{"val1", testMagicInt}
		testingArr := verifyTestingArr(driverArr)
		assert.Equal(t, driverArr, testingArr)

		driverArr2 := verifyDriverArr(testingArr)
		assert.Equal(t, testingArr, driverArr2)
	})
}
