package urlutil

import "strings"

func isLowerAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUpperAlpha(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isAlpha(c byte) bool {
	return isLowerAlpha(c) || isUpperAlpha(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isSchemeSpecial(c byte) bool {
	return c == '+' || c == '-' || c == '.'
}

func isValidSchemeFollowChar(c byte) bool {
	if isAlpha(c) || isDigit(c) {
		return true
	}
	return isSchemeSpecial(c)
}

func hasValidSchemeChars(scheme string) bool {
	for i := 1; i < len(scheme); i++ {
		if !isValidSchemeFollowChar(scheme[i]) {
			return false
		}
	}
	return true
}

// IsValidScheme checks if the provided scheme complies with RFC 3986 section 3.1:
// scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
func IsValidScheme(scheme string) bool {
	if scheme == "" {
		return false
	}
	if !isAlpha(scheme[0]) {
		return false
	}
	return hasValidSchemeChars(scheme)
}

// HasSchemePrefix returns true if raw begins with a valid RFC 3986 scheme followed by "://".
func HasSchemePrefix(raw string) bool {
	idx := strings.Index(raw, schemeSep)
	if idx <= 0 {
		return false
	}
	return IsValidScheme(raw[:idx])
}

// HasSchemeSingleSlash checks if raw begins with a valid RFC 3986 scheme followed by ":/" but NOT "://".
// Examples: "https:/host", "http:/admin:pass@minio:9000", "s3:/bucket/x", "x:/host".
func HasSchemeSingleSlash(raw string) bool {
	idx := strings.Index(raw, ":/")
	if idx <= 0 {
		return false
	}
	if !IsValidScheme(raw[:idx]) {
		return false
	}
	return !strings.HasPrefix(raw[idx:], schemeSep)
}
