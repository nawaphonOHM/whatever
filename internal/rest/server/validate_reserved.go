package server

import (
	"fmt"
	"net/http"
	"strings"
)

// isProbeSubpath checks if a path falls under reserved probe namespaces.
func isProbeSubpath(path string) bool {
	return strings.HasPrefix(path, ReservedHealthPath+"/") ||
		strings.HasPrefix(path, ReservedReadyPath+"/")
}

// isRootCatchAll checks if path is a root catch-all wildcard.
func isRootCatchAll(path string) bool {
	return strings.HasPrefix(path, "/*")
}

// errReservedProbePath formats an error for reserved probe paths and subpaths.
func errReservedProbePath(fullPath string) error {
	return fmt.Errorf(
		"%w: path %q is reserved for framework liveness/readiness probes",
		ErrReservedPath,
		fullPath,
	)
}

// errProbeConflict formats an error for routes that conflict with probe routes.
func errProbeConflict(fullPath string) error {
	return fmt.Errorf(
		"%w: root catch-all %q conflicts with framework probe routes",
		ErrReservedPath,
		fullPath,
	)
}

// checkGetProbeConflict checks probe subpath and catch-all conflicts for GET method.
func checkGetProbeConflict(fullPath string) error {
	if isProbeSubpath(fullPath) {
		return errReservedProbePath(fullPath)
	}
	if isRootCatchAll(fullPath) {
		return errProbeConflict(fullPath)
	}
	return nil
}

// checkReservedPath rejects framework probe paths and GET probe conflicts.
func checkReservedPath(method, fullPath string) error {
	if isReservedPath(fullPath) {
		return errReservedProbePath(fullPath)
	}
	if method == http.MethodGet {
		return checkGetProbeConflict(fullPath)
	}
	return nil
}
