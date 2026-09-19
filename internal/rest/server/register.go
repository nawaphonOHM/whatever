package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/health"
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

// wrapMiddleware adapts a framework Middleware to gin.HandlerFunc.
func wrapMiddleware(mw rest.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := rest.NewContext(ginCtx)
		mw(c)
	}
}

// writeIfPresent writes resp when non-nil.
func writeIfPresent(c *rest.Context, resp rest.Response) {
	if resp == nil {
		return
	}
	resp.Write(c.GinContext())
}

// invokeHandler runs the handler unless the context was aborted.
func invokeHandler(c *rest.Context, h rest.Handler) {
	if c.GinContext().IsAborted() {
		return
	}
	writeIfPresent(c, h(c))
}

// wrapHandler adapts a framework Handler to gin.HandlerFunc.
func wrapHandler(h rest.Handler) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		invokeHandler(rest.NewContext(ginCtx), h)
	}
}

// buildHandlers composes middlewares and handler into gin handlers.
func buildHandlers(
	middlewares []rest.Middleware,
	handler rest.Handler,
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
	registrations []*rest.RestAPIRegistration,
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
