package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v2/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/v2/internal/rest/health"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
)

// wrapMiddleware adapts a framework Middleware to gin.HandlerFunc.
func wrapMiddleware(mw contracts.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		reqCtx, _ := callstack.EnsureContext(ginCtx.Request.Context())
		ginCtx.Request = ginCtx.Request.WithContext(reqCtx)
		c := newContext(ginCtx)
		decorated := callstack.DecorateContext("rest.Middleware", func(context.Context) error {
			mw(c)
			return nil
		})
		if err := decorated(reqCtx); err != nil {
			logging.ErrorContext(reqCtx, "middleware execution failed", "err", err)
		}
	}
}

// writeResponse writes the response if non-nil.
func writeResponse(ginCtx *gin.Context, resp contracts.Response) {
	if resp != nil {
		resp.Write(ginCtx)
	}
}

// invokeHandler runs the handler unless the context was aborted.
func invokeHandler(ginCtx *gin.Context, h contracts.Handler) {
	if ginCtx.IsAborted() {
		return
	}
	reqCtx, _ := callstack.EnsureContext(ginCtx.Request.Context())
	ginCtx.Request = ginCtx.Request.WithContext(reqCtx)
	c := newContext(ginCtx)
	var resp contracts.Response
	decorated := callstack.DecorateContext("rest.Handler", func(context.Context) error {
		resp = h(c)
		return nil
	})
	if err := decorated(reqCtx); err != nil {
		logging.ErrorContext(reqCtx, "handler execution failed", "err", err)
	}
	writeResponse(ginCtx, resp)
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
	decorated := callstack.DecorateErr("server.RegisterRoutesWithVersion", func() error {
		validated, err := validateRegistrations(registrations)
		if err != nil {
			return err
		}
		logging.Info("registering REST routes", "route_count", len(validated), "version", version)
		mountHealthEndpoints(engine, version)
		logging.Info("registered reserved health endpoints", "paths", []string{ReservedHealthPath, ReservedReadyPath})
		mountvalidatedRoutes(engine, validated)
		logging.Info("route registration complete", "total_routes", len(validated)+2)
		return nil
	})
	return decorated()
}
