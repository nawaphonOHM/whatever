package rest

import (
	"net/http"

	intprob "github.com/nawaphonOHM/whatever/internal/rest/problem"
)

func newProblemDetails(
	statusCode int,
	code, detail string,
	details ...any,
) Response {
	return intprob.New(statusCode, code, detail, details...)
}

// Error creates an RFC 9457 Problem Details Response.
func Error(statusCode int, code, detail string, details ...any) Response {
	return newProblemDetails(statusCode, code, detail, details...)
}

// BadRequest creates a 400 Bad Request RFC 9457 Problem Details Response.
func BadRequest(code, detail string, details ...any) Response {
	return Error(http.StatusBadRequest, code, detail, details...)
}

// Unauthorized creates a 401 Unauthorized RFC 9457 Problem Details Response.
func Unauthorized(code, detail string, details ...any) Response {
	return Error(http.StatusUnauthorized, code, detail, details...)
}

// Forbidden creates a 403 Forbidden RFC 9457 Problem Details Response.
func Forbidden(code, detail string, details ...any) Response {
	return Error(http.StatusForbidden, code, detail, details...)
}

// NotFound creates a 404 Not Found RFC 9457 Problem Details Response.
func NotFound(code, detail string, details ...any) Response {
	return Error(http.StatusNotFound, code, detail, details...)
}

// InternalServerError creates a 500 Internal Server Error RFC 9457 Problem
// Details Response.
func InternalServerError(code, detail string, details ...any) Response {
	return Error(http.StatusInternalServerError, code, detail, details...)
}
