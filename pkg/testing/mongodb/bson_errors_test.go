package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

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
