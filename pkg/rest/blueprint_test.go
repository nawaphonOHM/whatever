package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCorsSetting_DefaultsAndBuilders(t *testing.T) {
	cors := NewCorsSetting().
		WithAllowOrigin("https://a.com", "https://b.com").
		WithAllowHTTPMethods(GET, POST)

	assert.Equal(
		t,
		[]string{"https://a.com", "https://b.com"},
		cors.AllowOrigin(),
	)
	assert.Equal(
		t,
		[]HTTPMethod{GET, POST},
		cors.AllowHTTPMethods(),
	)
}

func TestCorsSetting_NilSafety(t *testing.T) {
	var cors *CorsSetting
	assert.Nil(t, cors.AllowOrigin())
	assert.Nil(t, cors.AllowHTTPMethods())
}

func TestNewMeta_WithCors(t *testing.T) {
	cors := NewCorsSetting().WithAllowOrigin("https://example.com")
	meta := NewMeta().WithCors(cors)

	require.NotNil(t, meta.Cors())
	assert.Equal(
		t,
		[]string{"https://example.com"},
		meta.Cors().AllowOrigin(),
	)

	var nilMeta *Meta
	assert.Nil(t, nilMeta.Cors())
}

func TestNewBluePrint_BuildersAndAccessors(t *testing.T) {
	api1 := &ExportableAPI{Path: "/users", Method: GET}
	api2 := &ExportableAPI{Path: "/items", Method: POST}
	reg1 := &RRestAPIRegistration{
		Prefix: "/v1",
		Apis:   []*ExportableAPI{api1},
	}
	reg2 := &RRestAPIRegistration{
		Prefix: "/v2",
		Apis:   []*ExportableAPI{api2},
	}

	bp := NewBluePrint().
		WithMeta(NewMeta()).
		WithAPIs(reg1).
		AddAPIs(reg2)

	assert.NotNil(t, bp.Meta())
	assert.Len(t, bp.Apis(), 2)

	var nilBP *BluePrint
	assert.Nil(t, nilBP.Meta())
	assert.Nil(t, nilBP.Apis())
}
