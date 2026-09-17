package server

import (
	"errors"
	"fmt"
	"net/http"
	pathpkg "path"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/health"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Handler defines the function signature for route handlers returning a standardized Response.
type Handler func(c *Context) response.Response

// Middleware defines the function signature for route-level middlewares.
type Middleware func(c *Context)

// HTTPMethod represents supported HTTP request methods.
type HTTPMethod int

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

// GetName returns the standard uppercase HTTP method string.
// Panics if the method is unrecognized.
func (m HTTPMethod) GetName() string {
	switch m {
	case GET:
		return http.MethodGet
	case HEAD:
		return http.MethodHead
	case POST:
		return http.MethodPost
	case PUT:
		return http.MethodPut
	case PATCH:
		return http.MethodPatch
	case DELETE:
		return http.MethodDelete
	case OPTIONS:
		return http.MethodOptions
	case CONNECT:
		return http.MethodConnect
	case TRACE:
		return http.MethodTrace
	default:
		panic(fmt.Sprintf("unknown HTTP method: %d", m))
	}
}

// String implements fmt.Stringer for HTTPMethod.
func (m HTTPMethod) String() string {
	if m.IsValid() {
		return m.GetName()
	}
	return fmt.Sprintf("UNKNOWN(%d)", m)
}

// IsValid returns true if the HTTPMethod is one of the recognized constants.
func (m HTTPMethod) IsValid() bool {
	return m >= GET && m <= TRACE
}

// ApiVersioning represents the API major version (e.g. 1 for v1). 0 indicates unversioned.
type ApiVersioning uint

// Pathz represents a URL path segment.
type Pathz string

// ExportableApi defines a single API route endpoint with method, path, middlewares, and handler.
type ExportableApi struct {
	Path       Pathz
	Method     HTTPMethod
	Middleware []Middleware
	Handler    Handler
}

// RestApiRegistration groups multiple exportable APIs under a common prefix and version.
type RestApiRegistration struct {
	Version ApiVersioning
	Prefix  Pathz
	Apis    []*ExportableApi
}

// Reserved path constants for framework-managed health probes.
const (
	ReservedHealthPath = "/health"
	ReservedReadyPath  = "/ready"
)

// Sentinel validation errors.
var (
	ErrReservedPath    = errors.New("route conflicts with reserved health/readiness endpoint")
	ErrDuplicateRoute  = errors.New("duplicate route registration")
	ErrNilRegistration = errors.New("registration cannot be nil")
	ErrNilAPI          = errors.New("exportable API entry cannot be nil")
	ErrNilHandler      = errors.New("handler cannot be nil")
	ErrInvalidMethod   = errors.New("invalid HTTP method")
)

// ValidatedRoute represents a pre-validated and normalized route ready for registration on Gin.
type ValidatedRoute struct {
	Method      string
	Path        string
	Middlewares []Middleware
	Handler     Handler
}

// CalculateFullPath joins version, prefix, and path into a normalized canonical URL path.
// No implicit "/api" segment is inserted — callers supply the full prefix they want.
// When version > 0, a "/v{N}" segment is prepended unless the prefix already contains it.
func CalculateFullPath(version ApiVersioning, prefix Pathz, path Pathz) string {
	pfx := strings.TrimSpace(string(prefix))
	pth := strings.TrimSpace(string(path))

	if pfx != "" {
		if !strings.HasPrefix(pfx, "/") {
			pfx = "/" + pfx
		}
		pfx = strings.TrimRight(pfx, "/")
	}

	if version > 0 {
		vStr := fmt.Sprintf("v%d", version)
		vPrefix := "/" + vStr

		hasVersion := slices.Contains(strings.Split(pfx, "/"), vStr)

		if !hasVersion {
			if pfx != "" {
				pfx = vPrefix + pfx
			} else {
				pfx = vPrefix
			}
		}
	}

	if pth != "" {
		if !strings.HasPrefix(pth, "/") {
			pth = "/" + pth
		}
		if len(pth) > 1 {
			pth = strings.TrimRight(pth, "/")
		}
	}

	full := pfx + pth
	if full == "" {
		return "/"
	}

	cleaned := pathpkg.Clean(full)
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	return cleaned
}

// ValidateRegistrations inspects a slice of RestApiRegistration entries and verifies:
// 1. No nil registrations, nil APIs, or nil handlers
// 2. All HTTP methods are valid
// 3. No collision with framework-reserved /health and /ready endpoints
// 4. No duplicate concrete method/path combinations
func ValidateRegistrations(registrations []*RestApiRegistration) ([]*ValidatedRoute, error) {
	var validated []*ValidatedRoute
	seen := make(map[string]bool)

	for regIdx, reg := range registrations {
		if reg == nil {
			return nil, fmt.Errorf("%w: registration at index %d is nil", ErrNilRegistration, regIdx)
		}

		for apiIdx, api := range reg.Apis {
			if api == nil {
				return nil, fmt.Errorf("%w: API at index %d in prefix %q is nil", ErrNilAPI, apiIdx, reg.Prefix)
			}

			if api.Handler == nil {
				return nil, fmt.Errorf("%w: API at index %d in prefix %q has nil handler", ErrNilHandler, apiIdx, reg.Prefix)
			}

			if !api.Method.IsValid() {
				return nil, fmt.Errorf("%w: API at index %d has invalid method code %d", ErrInvalidMethod, apiIdx, api.Method)
			}

			fullPath := CalculateFullPath(reg.Version, reg.Prefix, api.Path)

			if fullPath == ReservedHealthPath || fullPath == ReservedReadyPath {
				return nil, fmt.Errorf("%w: path %q is reserved for framework liveness/readiness probes", ErrReservedPath, fullPath)
			}

			methodStr := api.Method.GetName()
			routeKey := methodStr + " " + fullPath

			if seen[routeKey] {
				return nil, fmt.Errorf("%w: %s is registered more than once", ErrDuplicateRoute, routeKey)
			}
			seen[routeKey] = true

			validated = append(validated, &ValidatedRoute{
				Method:      methodStr,
				Path:        fullPath,
				Middlewares: api.Middleware,
				Handler:     api.Handler,
			})
		}
	}

	return validated, nil
}

// RegisterRoutes validates the provided registrations, mounts framework health endpoints,
// and attaches all validated routes and middlewares to the Gin engine.
func RegisterRoutes(engine *gin.Engine, registrations []*RestApiRegistration) error {
	return RegisterRoutesWithVersion(engine, registrations, "")
}

// RegisterRoutesWithVersion validates the provided registrations, mounts framework health endpoints
// with version metadata, and attaches all validated routes and middlewares to the Gin engine.
func RegisterRoutesWithVersion(engine *gin.Engine, registrations []*RestApiRegistration, version string) error {
	validated, err := ValidateRegistrations(registrations)
	if err != nil {
		return err
	}

	// Mount reserved health endpoints (library-owned; not overridable by importers).
	hh := health.New(version)
	engine.GET(ReservedHealthPath, hh.Health)
	engine.GET(ReservedReadyPath, hh.Ready)

	// Mount validated routes
	for _, route := range validated {
		handlers := make([]gin.HandlerFunc, 0, len(route.Middlewares)+1)
		for _, mw := range route.Middlewares {
			if mw != nil {
				m := mw
				handlers = append(handlers, func(ginCtx *gin.Context) {
					c := NewContext(ginCtx)
					m(c)
				})
			}
		}
		h := route.Handler
		handlers = append(handlers, func(ginCtx *gin.Context) {
			c := NewContext(ginCtx)
			resp := h(c)
			if resp != nil && !ginCtx.IsAborted() {
				resp.Write(ginCtx)
			}
		})
		engine.Handle(route.Method, route.Path, handlers...)
	}

	return nil
}
