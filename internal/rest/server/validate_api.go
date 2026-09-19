package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/pkg/rest"
)

// routeKey builds a unique method+path registration key.
func routeKey(method, path string) string {
	return method + " " + path
}

// errNilAPIAt formats a nil API validation error.
func errNilAPIAt(apiIdx int, prefix rest.Pathz) error {
	return fmt.Errorf(
		"%w: API at index %d in prefix %q is nil",
		ErrNilAPI,
		apiIdx,
		prefix,
	)
}

// errNilHandlerAt formats a nil handler validation error.
func errNilHandlerAt(apiIdx int, prefix rest.Pathz) error {
	return fmt.Errorf(
		"%w: API at index %d in prefix %q has nil handler",
		ErrNilHandler,
		apiIdx,
		prefix,
	)
}

// errInvalidMethodAt formats an invalid method validation error.
func errInvalidMethodAt(apiIdx int, method rest.HTTPMethod) error {
	return fmt.Errorf(
		"%w: API at index %d has invalid method code %d",
		ErrInvalidMethod,
		apiIdx,
		method,
	)
}

// checkAPINonNil ensures API and handler pointers exist.
func checkAPINonNil(api *rest.ExportableAPI, apiIdx int,
	prefix rest.Pathz) error {
	if api == nil {
		return errNilAPIAt(apiIdx, prefix)
	}
	if api.Handler == nil {
		return errNilHandlerAt(apiIdx, prefix)
	}
	return nil
}

// validateAPIEntry checks a single exportable API entry.
func validateAPIEntry(
	api *rest.ExportableAPI,
	apiIdx int,
	prefix rest.Pathz,
) error {
	if err := checkAPINonNil(api, apiIdx, prefix); err != nil {
		return err
	}
	if !api.Method.IsValid() {
		return errInvalidMethodAt(apiIdx, api.Method)
	}
	return nil
}

// checkReservedPath rejects framework-owned probe paths.
func checkReservedPath(fullPath string) error {
	if fullPath != ReservedHealthPath && fullPath != ReservedReadyPath {
		return nil
	}
	return fmt.Errorf(
		"%w: path %q is reserved for framework "+
			"liveness/readiness probes",
		ErrReservedPath,
		fullPath,
	)
}

// checkDuplicate marks a route key or reports collision.
func checkDuplicate(seen map[string]bool, key string) error {
	if !seen[key] {
		seen[key] = true
		return nil
	}
	return fmt.Errorf(
		"%w: %s is registered more than once",
		ErrDuplicateRoute,
		key,
	)
}

// buildValidatedRoute creates a route and checks collisions.
func buildValidatedRoute(
	reg *rest.RestAPIRegistration,
	api *rest.ExportableAPI,
	seen map[string]bool,
) (*validatedRoute, error) {
	fullPath := CalculateFullPath(reg.Version, reg.Prefix, api.Path)
	if err := checkReservedPath(fullPath); err != nil {
		return nil, err
	}
	methodStr := api.Method.GetName()
	key := routeKey(methodStr, fullPath)
	if err := checkDuplicate(seen, key); err != nil {
		return nil, err
	}
	return &validatedRoute{
		Method:      methodStr,
		Path:        fullPath,
		Handler:     api.Handler,
		Middlewares: api.Middleware,
	}, nil
}
