package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

// D represents an ordered BSON document.
type D = bson.D

// E represents a single BSON element in a D document.
type E = bson.E

// M represents an unordered BSON document.
type M = bson.M

// A represents a BSON array.
type A = bson.A

// NewD constructs an ordered BSON document from elements. With no elements, it
// returns a non-nil empty document.
func NewD(elements ...E) D {
	if len(elements) == 0 {
		return D{}
	}

	return D(elements)
}

// NewE constructs a single BSON element with the provided key and value.
func NewE(key string, value any) E {
	return E{Key: key, Value: value}
}

// NewM constructs an empty unordered BSON map.
func NewM() M {
	return make(M)
}

// NewA constructs a BSON array from values. With no values, it returns a
// non-nil empty array.
func NewA(values ...any) A {
	if len(values) == 0 {
		return A{}
	}

	return A(values)
}

// Marshal returns the BSON encoding of val.
func Marshal(val any) ([]byte, error) {
	return bson.Marshal(val)
}

// Unmarshal parses the BSON-encoded data and stores the result in the value pointed to by val.
func Unmarshal(data []byte, val any) error {
	return bson.Unmarshal(data, val)
}
