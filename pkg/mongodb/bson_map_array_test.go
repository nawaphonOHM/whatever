package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

func verifyDriverMap(m bson.M) bson.M {
	return m
}

func verifyPkgMap(m mongodb.M) mongodb.M {
	return m
}

func verifyDriverArr(a bson.A) bson.A {
	return a
}

func verifyPkgArr(a mongodb.A) mongodb.A {
	return a
}

// TestBSONMapArrayAliasesInteroperability verifies assignment compatibility for M and A.
func TestBSONMapArrayAliasesInteroperability(t *testing.T) {
	driverMap := bson.M{"hello": "world"}
	pkgMap := verifyPkgMap(driverMap)
	assert.Equal(t, driverMap, pkgMap)

	driverMap2 := verifyDriverMap(pkgMap)
	assert.Equal(t, pkgMap, driverMap2)

	driverArr := bson.A{"item1", 1}
	pkgArr := verifyPkgArr(driverArr)
	assert.Equal(t, driverArr, pkgArr)

	driverArr2 := verifyDriverArr(pkgArr)
	assert.Equal(t, pkgArr, driverArr2)
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
		arr := mongodb.NewA("first", 1, subDoc)
		require.NotNil(t, arr)
		assert.Equal(t, "first", arr[0])
		assert.Equal(t, 1, arr[1])
		assert.Equal(t, subDoc, arr[2])
	})
}
