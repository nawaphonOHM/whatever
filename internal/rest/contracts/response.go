package contracts

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// MediaTypeProblemJSON is the canonical media type for RFC 9457 Problem
// Details.
const MediaTypeProblemJSON = "application/problem+json"

// Response represents any HTTP response that can write itself to a Gin
// context.
type Response interface {
	// StatusCode returns the HTTP status code of the response.
	StatusCode() int
	// Write renders the response to the given Gin context.
	Write(*gin.Context)
	// WithHeader attaches a custom key-value header to the response.
	WithHeader(string, string) Response
	// WithHeaderInt attaches a custom header whose integer value is converted to string.
	WithHeaderInt(string, int) Response
	// WithHeaderStringer attaches a custom header whose value is obtained from fmt.Stringer.
	WithHeaderStringer(string, fmt.Stringer) Response
}

// SafeStringer converts a fmt.Stringer to string safely handling untyped nil
// and typed-nil pointers.
func SafeStringer(s fmt.Stringer) (string, bool) {
	if isNilStringer(s) {
		return "", false
	}
	return s.String(), true
}

func isNilStringer(s fmt.Stringer) bool {
	if s == nil {
		return true
	}
	return isNilKind(s)
}

func isNilKind(v any) bool {
	val := reflect.ValueOf(v)
	if !isNillableKind(val.Kind()) {
		return false
	}
	return val.IsNil()
}

func isNillableKind(k reflect.Kind) bool {
	switch k {
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

// SetHeader attaches a canonicalized key-value header into the target headers map, allocating if necessary.
func SetHeader(headers *map[string]string, key, value string) {
	if key == "" {
		return
	}
	if *headers == nil {
		*headers = make(map[string]string)
	}
	(*headers)[http.CanonicalHeaderKey(key)] = value
}

// WriteHeaders writes the provided headers map to the gin.Context Writer header map.
func WriteHeaders(c *gin.Context, headers map[string]string) {
	if !canWriteHeaders(c, headers) {
		return
	}
	applyHeaders(c.Writer.Header(), headers)
}

func canWriteHeaders(c *gin.Context, headers map[string]string) bool {
	if c == nil {
		return false
	}
	if c.Writer == nil {
		return false
	}
	return len(headers) > 0
}

func applyHeaders(h http.Header, headers map[string]string) {
	for k, v := range headers {
		applyHeader(h, k, v)
	}
}

func applyHeader(h http.Header, key, value string) {
	if key != "" {
		h.Set(key, value)
	}
}
