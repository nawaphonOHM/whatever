package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/internal/rest/health"
)

// wrapMiddleware adapts a framework Middleware to gin.HandlerFunc.
func wrapMiddleware(mw contracts.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := newContext(ginCtx)
		mw(c)
	}
}

// invokeHandler runs the handler unless the context was aborted.
func invokeHandler(ginCtx *gin.Context, h contracts.Handler) {
	if ginCtx.IsAborted() {
		return
	}
	c := newContext(ginCtx)
	if resp := h(c); resp != nil {
		resp.Write(ginCtx)
	}
}

// wrapHandler adapts a framework Handler to gin.HandlerFunc.
func wrapHandler(h contracts.Handler) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		invokeHandler(ginCtx, h)
	}
}

// buildHandlers composes middlewares and handler into gin handlers.
func buildHandlers(
	middlewares []contracts.Middleware,
	handler contracts.Handler,
) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0, len(middlewares)+1)
	for _, mw := range middlewares {
		handlers = append(handlers, wrapMiddleware(mw))
	}
	return append(handlers, wrapHandler(handler))
}

// mountvalidatedRoutes registers validated routes onto the engine.
func mountvalidatedRoutes(
	engine *gin.Engine,
	routes []*validatedRoute,
) {
	for _, route := range routes {
		handlers := buildHandlers(route.Middlewares, route.Handler)
		engine.Handle(route.Method, route.Path, handlers...)
	}
}

// mountHealthEndpoints registers reserved health probes.
func mountHealthEndpoints(engine *gin.Engine, version string) {
	hh := health.New(version)
	engine.GET(ReservedHealthPath, hh.Health)
	engine.GET(ReservedReadyPath, hh.Ready)
}

// RegisterRoutesWithVersion validates and mounts API registrations.
func RegisterRoutesWithVersion(
	engine *gin.Engine,
	registrations []*contracts.RRestAPIRegistration,
	version string,
) error {
	validated, err := validateRegistrations(registrations)
	if err != nil {
		return err
	}
	mountHealthEndpoints(engine, version)
	mountvalidatedRoutes(engine, validated)
	return nil
}
