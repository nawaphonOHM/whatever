package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartREST_NilBluePrint(t *testing.T) {
	err := StartREST(nil)
	assert.ErrorIs(t, err, ErrNilBluePrint)
}

func TestStartREST_NilRegistration(t *testing.T) {
	bp := NewBluePrint().WithAPIs(nil)
	err := StartREST(bp)
	assert.ErrorIs(t, err, ErrNilRegistration)
}

func TestStartREST_ReservedPathConflict(t *testing.T) {
	api := &ExportableAPI{
		Path:    "/health",
		Method:  GET,
		Handler: func(Context) Response { return OK("ok") },
	}
	reg := &RRestAPIRegistration{
		Apis: []*ExportableAPI{api},
	}
	bp := NewBluePrint().WithAPIs(reg)

	err := StartREST(bp)
	assert.ErrorIs(t, err, ErrReservedPath)
}

func TestStartREST_DuplicateRouteConflict(t *testing.T) {
	api1 := &ExportableAPI{
		Path:    "/ping",
		Method:  GET,
		Handler: func(Context) Response { return OK("pong") },
	}
	api2 := &ExportableAPI{
		Path:    "/ping",
		Method:  GET,
		Handler: func(Context) Response { return OK("pong") },
	}
	reg := &RRestAPIRegistration{
		Apis: []*ExportableAPI{api1, api2},
	}
	bp := NewBluePrint().WithAPIs(reg)

	err := StartREST(bp)
	assert.ErrorIs(t, err, ErrDuplicateRoute)
}
