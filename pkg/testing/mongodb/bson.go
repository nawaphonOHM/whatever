package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

// Marshal returns the BSON encoding of val.
func Marshal(val any) ([]byte, error) {
	return bson.Marshal(val)
}

// Unmarshal parses the BSON-encoded data and stores the result in the value pointed to by val.
func Unmarshal(data []byte, val any) error {
	return bson.Unmarshal(data, val)
}
