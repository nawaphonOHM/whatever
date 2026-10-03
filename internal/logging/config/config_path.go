package config

import "strings"

func hasPathPrefix(path string) bool {
	return strings.HasPrefix(path, "/") ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../")
}

func hasPathSeparator(path string) bool {
	return strings.Contains(path, "/") || strings.Contains(path, "\\")
}

func hasLogExtension(norm string) bool {
	return strings.HasSuffix(norm, ".log") || strings.HasSuffix(norm, ".txt")
}

// isLikelyFilePath checks if the output string represents a file path.
func isLikelyFilePath(path string) bool {
	norm := strings.ToLower(path)
	if hasPathPrefix(path) {
		return true
	}
	if hasPathSeparator(path) {
		return true
	}
	return hasLogExtension(norm)
}
