package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

func verifyDriverDoc(d bson.D) bson.D {
	return d
}

func verifyPkgDoc(d mongodb.D) mongodb.D {
	return d
}

func verifyDriverElem(e bson.E) bson.E {
	return e
}

func verifyPkgElem(e mongodb.E) mongodb.E {
	return e
}

// TestBSONTypeAliasesInteroperability verifies compile-time and runtime
// assignment compatibility between pkg/mongodb aliases and mongo-driver bson types.
func TestBSONTypeAliasesInteroperability(t *testing.T) {
	driverElem := bson.E{Key: "k1", Value: "v1"}
	pkgElem := verifyPkgElem(driverElem)
	assert.Equal(t, driverElem, pkgElem)

	driverElem2 := verifyDriverElem(pkgElem)
	assert.Equal(t, pkgElem, driverElem2)

	driverDoc := bson.D{{Key: "k1", Value: "v1"}}
	pkgDoc := verifyPkgDoc(driverDoc)
	assert.Equal(t, driverDoc, pkgDoc)

	driverDoc2 := verifyDriverDoc(pkgDoc)
	assert.Equal(t, pkgDoc, driverDoc2)
}

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
		e0 := mongodb.NewE("first", 0)
		e1 := mongodb.NewE("second", 1)
		e2 := mongodb.NewE("third", 2)

		doc := mongodb.NewD(e0, e1, e2)
		require.NotNil(t, doc)
		assert.Equal(t, "first", doc[0].Key)
		assert.Equal(t, 0, doc[0].Value)
		assert.Equal(t, "second", doc[1].Key)
		assert.Equal(t, 1, doc[1].Value)
		assert.Equal(t, "third", doc[2].Key)
		assert.Equal(t, 2, doc[2].Value)
	})
}
