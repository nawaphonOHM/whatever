// Package contracts provides core REST types and abstractions.
package contracts

import (
	"fmt"
	"net/http"
)

// HTTPMethod represents supported HTTP request methods.
type HTTPMethod int

// Supported HTTP method constants.
const (
	GET HTTPMethod = iota
	HEAD
	POST
	PUT
	PATCH
	DELETE
	OPTIONS
	CONNECT
	TRACE
)

var httpMethodNames = map[HTTPMethod]string{
	GET:     http.MethodGet,
	HEAD:    http.MethodHead,
	POST:    http.MethodPost,
	PUT:     http.MethodPut,
	PATCH:   http.MethodPatch,
	DELETE:  http.MethodDelete,
	OPTIONS: http.MethodOptions,
	CONNECT: http.MethodConnect,
	TRACE:   http.MethodTrace,
}

// GetName returns the standard uppercase HTTP method string.
// Panics if the method is unrecognized.
func (m HTTPMethod) GetName() string {
	name, ok := httpMethodNames[m]
	if !ok {
		panic(fmt.Sprintf("unknown HTTP method: %d", m))
	}
	return name
}

// String implements fmt.Stringer for HTTPMethod.
func (m HTTPMethod) String() string {
	if m.IsValid() {
		return m.GetName()
	}
	return fmt.Sprintf("UNKNOWN(%d)", m)
}

// IsValid returns true if the HTTPMethod is recognized.
func (m HTTPMethod) IsValid() bool {
	return m >= GET && m <= TRACE
}
