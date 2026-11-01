package urlutil

import (
	"net/url"
	"strings"
)

func isBareHostName(name string) bool {
	return name == "http" || name == "https"
}

func isBareSchemeHost(u *url.URL) bool {
	if isBareHostName(u.Hostname()) {
		return true
	}
	return u.Host == "http:" || u.Host == "https:"
}

func isColonHost(host string) bool {
	return strings.HasSuffix(host, ":") || strings.HasPrefix(host, ":")
}

func isMalformedHost(u *url.URL) bool {
	// Leading or trailing colons and bare scheme names can be misparsed as hosts.
	if isColonHost(u.Host) {
		return true
	}
	return isBareSchemeHost(u)
}

func isMalformedRaw(raw string) bool {
	// Single-slash schemes (for example, "https:/host") are ambiguous.
	if HasSchemeSingleSlash(raw) {
		return true
	}
	// Protocol-relative and backslash-prefixed inputs can hide userinfo from the
	// parsed authority (for example, "//user:pass@host" or "\\host").
	return strings.HasPrefix(raw, "/") || strings.Contains(raw, "\\")
}

func isMalformedHostOrScheme(raw string, u *url.URL) bool {
	if isMalformedRaw(raw) {
		return true
	}
	return isMalformedHost(u)
}

func hasPathAtChar(u *url.URL) bool {
	if strings.Contains(u.Path, atChar) || strings.Contains(u.RawPath, atChar) {
		return true
	}
	return strings.Contains(u.Opaque, atChar)
}

func hasMalformedHostAuthority(host string) bool {
	return host == "" || strings.Contains(host, atChar)
}

func hasMalformedPathAuthority(u *url.URL) bool {
	// A path pseudo-scheme such as "a/b://user:pass@host" can hide credentials
	// in the path instead of the authority.
	return hasPathAtChar(u) || strings.Contains(u.Path, schemeSep)
}

func hasMalformedAuthority(u *url.URL) bool {
	if hasMalformedHostAuthority(u.Host) {
		return true
	}
	return hasMalformedPathAuthority(u)
}

func isLeakedUserInfo(raw string, u *url.URL) bool {
	if !strings.Contains(raw, atChar) {
		return false
	}
	return hasMalformedAuthority(u)
}

// isAmbiguousMalformed identifies parsed forms where malformed authority or
// scheme structure could leave userinfo outside the URL's authority.
func isAmbiguousMalformed(raw string, u *url.URL) bool {
	if isMalformedHostOrScheme(raw, u) {
		return true
	}
	return isLeakedUserInfo(raw, u)
}
