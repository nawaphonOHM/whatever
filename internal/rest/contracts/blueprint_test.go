package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	prefixV1 = "/v1"
	prefixV2 = "/v2"
)

func TestMetaNewAndGetters(t *testing.T) {
	meta := NewMeta()
	assert.NotNil(t, meta)
	assert.Empty(t, meta.CorsPath())

	meta.WithCorsPath("/path/to/cors.yaml")
	assert.Equal(t, "/path/to/cors.yaml", meta.CorsPath())
}

func TestMetaNilReceiver(t *testing.T) {
	var meta *Meta
	assert.Empty(t, meta.CorsPath())
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
	meta := NewMeta().WithCorsPath("/configs/cors.yaml")
	api := &RRestAPIRegistration{Prefix: "/api"}

	bp := NewBluePrint().WithMeta(meta).WithAPIs(api)
	assert.Equal(t, "/configs/cors.yaml", bp.Meta().CorsPath())
	assert.Len(t, bp.Apis(), 1)
}
