package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
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
