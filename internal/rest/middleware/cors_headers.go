package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	wildcardOrigin = "*"
	commaSep       = ", "
	headerAllowOrg = "Access-Control-Allow-Origin"
	headerAllowMth = "Access-Control-Allow-Methods"
	headerAllowHdr = "Access-Control-Allow-Headers"
	headerExposeHd = "Access-Control-Expose-Headers"
	headerAllowCrd = "Access-Control-Allow-Credentials"
	headerMaxAge   = "Access-Control-Max-Age"
	headerVary     = "Vary"
	headerValOrg   = "Origin"
)

func matchSpecificOrigin(origins []string, origin string) (string, bool) {
	if origin != "" && slices.Contains(origins, origin) {
		return origin, true
	}
	return "", false
}

func matchAllowedOrigin(origins []string, origin string) (string, bool) {
	if slices.Contains(origins, wildcardOrigin) {
		return wildcardOrigin, true
	}
	return matchSpecificOrigin(origins, origin)
}

func setHeaderIfNonEmpty(w http.ResponseWriter, key, val string) {
	if val != "" {
		w.Header().Set(key, val)
	}
}

func setConfigHeaders(w http.ResponseWriter, cfg *CORSConfig) {
	setHeaderIfNonEmpty(w, headerAllowMth, strings.Join(cfg.AllowMethods, commaSep))
	setHeaderIfNonEmpty(w, headerAllowHdr, strings.Join(cfg.AllowHeaders, commaSep))
	setHeaderIfNonEmpty(w, headerExposeHd, strings.Join(cfg.ExposeHeaders, commaSep))
	if cfg.AllowCredentials {
		w.Header().Set(headerAllowCrd, "true")
	}
	setHeaderIfNonEmpty(w, headerMaxAge, cfg.MaxAge)
}

func isOriginToken(token string) bool {
	trimmed := strings.TrimSpace(token)
	return trimmed == wildcardOrigin || strings.EqualFold(trimmed, headerValOrg)
}

func hasOriginToken(val string) bool {
	for token := range strings.SplitSeq(val, ",") {
		if isOriginToken(token) {
			return true
		}
	}
	return false
}

func addVaryOrigin(w http.ResponseWriter) {
	if !slices.ContainsFunc(w.Header().Values(headerVary), hasOriginToken) {
		w.Header().Add(headerVary, headerValOrg)
	}
}

func applyVaryHeader(w http.ResponseWriter, origins []string) {
	if !slices.Contains(origins, wildcardOrigin) {
		addVaryOrigin(w)
	}
}

func hasConfiguredOrigins(cfg *CORSConfig) bool {
	return cfg != nil && len(cfg.AllowOrigins) > 0
}

func setCORSHeaders(c *gin.Context, cfg *CORSConfig) bool {
	if !hasConfiguredOrigins(cfg) {
		return false
	}
	applyVaryHeader(c.Writer, cfg.AllowOrigins)
	allowedOrigin, ok := matchAllowedOrigin(cfg.AllowOrigins, c.Request.Header.Get(headerValOrg))
	if !ok {
		return false
	}
	c.Writer.Header().Set(headerAllowOrg, allowedOrigin)
	setConfigHeaders(c.Writer, cfg)
	return true
}
