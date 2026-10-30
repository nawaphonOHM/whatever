package core

import (
	"net/url"
)

func parsedEscapedPath(u *url.URL) string {
	if u.RawPath != "" || u.Path != "" {
		return u.EscapedPath()
	}
	return ""
}

func parseURLPath(p string) string {
	u, err := url.Parse(p)
	if err != nil {
		return ""
	}
	return parsedEscapedPath(u)
}

func escapePath(p string) string {
	if p == "" {
		return ""
	}
	if res := parseURLPath(p); res != "" {
		return res
	}
	return (&url.URL{Path: p}).EscapedPath()
}

func composeURL(pathVal, queryVal string) string {
	escapedPath := escapePath(pathVal)
	if queryVal != "" {
		return escapedPath + "?" + queryVal
	}
	return escapedPath
}

func assignComposedURL(req *HttpRequest, pathVal, queryVal string) {
	if pathVal != "" {
		req.RequestURL = composeURL(pathVal, queryVal)
		return
	}
	if queryVal != "" {
		req.RequestURL = "?" + queryVal
	}
}

func finalizeRequestURL(req *HttpRequest, pathVal, queryVal string) {
	if req.RequestURL != "" {
		return
	}
	assignComposedURL(req, pathVal, queryVal)
}
