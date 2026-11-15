package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBluePrint_BuildersAndAccessors(t *testing.T) {
	reg1 := &RRestAPIRegistration{
		Prefix: "/v1",
		Apis:   []*ExportableAPI{{Path: "/users", Method: GET}},
	}
	reg2 := &RRestAPIRegistration{
		Prefix: "/v2",
		Apis:   []*ExportableAPI{{Path: "/items", Method: POST}},
	}

	bp := NewBluePrint().
		WithMeta(NewMeta()).
		WithAPIs(reg1).
		AddAPIs(reg2)

	assert.NotNil(t, bp.Meta())
	assert.Len(t, bp.Apis(), 2)
}

func TestNewBluePrint_NilSafety(t *testing.T) {
	var nilBP *BluePrint
	assert.Nil(t, nilBP.Meta())
	assert.Nil(t, nilBP.Apis())
}
