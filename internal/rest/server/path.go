package server

import (
	"fmt"
	pathpkg "path"
	"slices"
	"strings"
)

// normalizePrefix ensures a leading slash and no trailing slash.
func normalizePrefix(prefix string) string {
	pfx := strings.TrimSpace(prefix)
	if pfx == "" {
		return ""
	}
	if !strings.HasPrefix(pfx, pathSeparator) {
		pfx = pathSeparator + pfx
	}
	return strings.TrimRight(pfx, pathSeparator)
}

// ensureLeadingSlash adds a leading slash when missing.
func ensureLeadingSlash(pth string) string {
	if strings.HasPrefix(pth, pathSeparator) {
		return pth
	}
	return pathSeparator + pth
}

// normalizePath ensures a leading slash and trims trailing slash.
func normalizePath(path string) string {
	pth := strings.TrimSpace(path)
	if pth == "" {
		return ""
	}
	pth = ensureLeadingSlash(pth)
	if len(pth) > 1 {
		pth = strings.TrimRight(pth, pathSeparator)
	}
	return pth
}

// versionToken returns the bare version segment (e.g. "v1").
func versionToken(version APIVersioning) string {
	return fmt.Sprintf("v%d", version)
}

// applyVersion prefixes pfx with /vN when version is set and missing.
func applyVersion(version APIVersioning, pfx string) string {
	if version == 0 {
		return pfx
	}
	vStr := versionToken(version)
	vPrefix := pathSeparator + vStr
	if slices.Contains(strings.Split(pfx, pathSeparator), vStr) {
		return pfx
	}
	return vPrefix + pfx
}

// cleanFullPath cleans and ensures a rooted absolute path.
func cleanFullPath(full string) string {
	if full == "" {
		return pathSeparator
	}
	cleaned := pathpkg.Clean(full)
	if !strings.HasPrefix(cleaned, pathSeparator) {
		cleaned = pathSeparator + cleaned
	}
	return cleaned
}

// isReservedPath reports framework-managed probe endpoints.
func isReservedPath(fullPath string) bool {
	return fullPath == ReservedHealthPath ||
		fullPath == ReservedReadyPath
}

// CalculateFullPath builds the canonical route path.
// Reserved health/readiness endpoints are never version-prefixed.
func CalculateFullPath(
	version APIVersioning,
	prefix, path Pathz,
) string {
	pfx := normalizePrefix(string(prefix))
	pth := normalizePath(string(path))
	fullPath := cleanFullPath(pfx + pth)
	// Reserved probes stay unversioned for kube-style checks.
	if isReservedPath(fullPath) {
		return fullPath
	}
	pfx = applyVersion(version, pfx)
	return cleanFullPath(pfx + pth)
}
