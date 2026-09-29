package server

import (
	"fmt"
	pathpkg "path"
	"strings"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// versionToken returns the bare version segment (e.g. "v1").
func versionToken(version contracts.APIVersioning) string {
	return fmt.Sprintf("v%d", version)
}

// versionString returns the version token if version > 0.
func versionString(version contracts.APIVersioning) string {
	if version > 0 {
		return versionToken(version)
	}
	return ""
}

// isLeadingToken reports if seg matches the api segment or version token.
func isLeadingToken(seg, apiSeg, vStr string) bool {
	if seg == apiSeg {
		return true
	}
	return vStr != "" && seg == vStr
}

// extractSegments splits a path into cleaned non-empty segments.
func extractSegments(p string) []string {
	trimmed := strings.Trim(pathpkg.Clean(p), pathSeparator)
	if trimmed == "" || trimmed == "." {
		return nil
	}
	return strings.Split(trimmed, pathSeparator)
}

// stripLeadingPrefixes removes redundant leading api and version tokens.
func stripLeadingPrefixes(segments []string, version contracts.APIVersioning) []string {
	apiSeg := strings.TrimPrefix(defaultAPIPrefix, pathSeparator)
	vStr := versionString(version)
	start := 0
	for start < len(segments) {
		if !isLeadingToken(segments[start], apiSeg, vStr) {
			break
		}
		start++
	}
	return segments[start:]
}

// applyAPIPrefix builds the canonical /api/v{N}/... or /api/... route path.
func applyAPIPrefix(version contracts.APIVersioning, rawPath string) string {
	segments := stripLeadingPrefixes(extractSegments(rawPath), version)
	result := []string{strings.TrimPrefix(defaultAPIPrefix, pathSeparator)}
	if version > 0 {
		result = append(result, versionToken(version))
	}
	result = append(result, segments...)
	return pathSeparator + strings.Join(result, pathSeparator)
}

// cleanFullPath cleans and ensures a rooted absolute path.
func cleanFullPath(prefix, path string) string {
	return pathpkg.Clean(pathSeparator + prefix + pathSeparator + path)
}

// isReservedPath reports framework-managed probe endpoints.
func isReservedPath(fullPath string) bool {
	return fullPath == ReservedHealthPath || fullPath == ReservedReadyPath
}

// CalculateFullPath builds the canonical route path.
// Reserved health/readiness endpoints are never version-prefixed or API-prefixed.
func CalculateFullPath(
	version contracts.APIVersioning,
	prefix, path contracts.Pathz,
) string {
	fullPath := cleanFullPath(string(prefix), string(path))
	if isReservedPath(fullPath) {
		return fullPath
	}
	return applyAPIPrefix(version, fullPath)
}
