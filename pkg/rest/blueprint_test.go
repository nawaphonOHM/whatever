package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMeta_WithCorsPath(t *testing.T) {
	meta := NewMeta().WithCorsPath("configs/cors.json")
	assert.Equal(t, "configs/cors.json", meta.CorsPath())

	var nilMeta *Meta
	assert.Empty(t, nilMeta.CorsPath())
}

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
		WithMeta(NewMeta().WithCorsPath("cors.yaml")).
		WithAPIs(reg1).
		AddAPIs(reg2)

	assert.NotNil(t, bp.Meta())
	assert.Equal(t, "cors.yaml", bp.Meta().CorsPath())
	assert.Len(t, bp.Apis(), 2)
}

func TestNewBluePrint_NilSafety(t *testing.T) {
	var nilBP *BluePrint
	assert.Nil(t, nilBP.Meta())
	assert.Nil(t, nilBP.Apis())
}
