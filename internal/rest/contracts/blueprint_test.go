package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	prefixV1 = "/v1"
	prefixV2 = "/v2"
)

func TestCorsSettingNewAndGetters(t *testing.T) {
	cors := NewCorsSetting()
	assert.NotNil(t, cors)
	assert.Nil(t, cors.AllowOrigin())
	assert.Nil(t, cors.AllowHTTPMethods())

	cors.WithAllowOrigin("http://localhost", "https://example.com")
	cors.WithAllowHTTPMethods(GET, POST)
	assert.Equal(t, []string{"http://localhost", "https://example.com"},
		cors.AllowOrigin())
	assert.Equal(t, []HTTPMethod{GET, POST}, cors.AllowHTTPMethods())
}

func TestCorsSettingImmutability(t *testing.T) {
	origins := []string{"https://example.com"}
	methods := []HTTPMethod{GET}
	cors := NewCorsSetting().WithAllowOrigin(origins...).
		WithAllowHTTPMethods(methods...)

	origins[0] = "mutated"
	methods[0] = DELETE
	assert.Equal(t, []string{"https://example.com"}, cors.AllowOrigin())
	assert.Equal(t, []HTTPMethod{GET}, cors.AllowHTTPMethods())

	gotOrigins := cors.AllowOrigin()
	gotOrigins[0] = "mutated2"
	assert.Equal(t, []string{"https://example.com"}, cors.AllowOrigin())
}

func TestCorsSettingNilReceiver(t *testing.T) {
	var cors *CorsSetting
	assert.Nil(t, cors.AllowOrigin())
	assert.Nil(t, cors.AllowHTTPMethods())
}

func TestMetaNewAndGetters(t *testing.T) {
	meta := NewMeta()
	assert.NotNil(t, meta)
	assert.Nil(t, meta.Cors())

	cors := NewCorsSetting().WithAllowOrigin("*")
	meta.WithCors(cors)
	assert.Same(t, cors, meta.Cors())
}

func TestMetaNilReceiver(t *testing.T) {
	var meta *Meta
	assert.Nil(t, meta.Cors())
}

func TestBluePrintNewAndGetters(t *testing.T) {
	bp := NewBluePrint()
	assert.NotNil(t, bp)
	assert.Nil(t, bp.Meta())
	assert.Nil(t, bp.Apis())

	meta := NewMeta()
	api := &RRestAPIRegistration{Prefix: "/test"}
	bp.WithMeta(meta).WithAPIs(api)

	assert.Same(t, meta, bp.Meta())
	assert.Equal(t, []*RRestAPIRegistration{api}, bp.Apis())
}

func TestBluePrintWithAPIsReplace(t *testing.T) {
	api1 := &RRestAPIRegistration{Prefix: prefixV1}
	api2 := &RRestAPIRegistration{Prefix: prefixV2}
	bp := NewBluePrint().WithAPIs(api1)
	assert.Len(t, bp.Apis(), 1)

	bp.WithAPIs(api2)
	assert.Equal(t, []*RRestAPIRegistration{api2}, bp.Apis())
}

func TestBluePrintAddAPIs(t *testing.T) {
	api1 := &RRestAPIRegistration{Prefix: prefixV1}
	api2 := &RRestAPIRegistration{Prefix: prefixV2}
	bp := NewBluePrint().WithAPIs(api1).AddAPIs(api2)

	assert.Equal(t, []*RRestAPIRegistration{api1, api2}, bp.Apis())
}

func TestBluePrintImmutability(t *testing.T) {
	api1 := &RRestAPIRegistration{Prefix: prefixV1}
	apis := []*RRestAPIRegistration{api1}
	bp := NewBluePrint().WithAPIs(apis...)

	apis[0] = &RRestAPIRegistration{Prefix: "/mutated"}
	assert.Equal(t, prefixV1, string(bp.Apis()[0].Prefix))

	got := bp.Apis()
	got[0] = &RRestAPIRegistration{Prefix: "/mutated2"}
	assert.Equal(t, prefixV1, string(bp.Apis()[0].Prefix))
}

func TestBluePrintNilReceiver(t *testing.T) {
	var bp *BluePrint
	assert.Nil(t, bp.Meta())
	assert.Nil(t, bp.Apis())
}

func TestBluePrintFluentChaining(t *testing.T) {
	cors := NewCorsSetting().WithAllowOrigin("*").
		WithAllowHTTPMethods(GET, POST, OPTIONS)
	meta := NewMeta().WithCors(cors)
	api := &RRestAPIRegistration{Prefix: "/api"}

	bp := NewBluePrint().WithMeta(meta).WithAPIs(api)
	assert.Equal(t, []string{"*"}, bp.Meta().Cors().AllowOrigin())
	assert.Len(t, bp.Apis(), 1)
}
