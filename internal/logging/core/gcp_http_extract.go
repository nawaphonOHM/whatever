package core

import "log/slog"

func setStatus(req *HttpRequest, a slog.Attr) bool {
	if req.Status != 0 {
		return false
	}
	st, ok := parseHTTPStatus(a.Value.Resolve().Any())
	if !ok {
		return false
	}
	req.Status = st
	return true
}

func applyStatus(req *HttpRequest, a slog.Attr) bool {
	if !isStatusKey(a.Key) {
		return false
	}
	return setStatus(req, a)
}

func applyStatusOrLatency(req *HttpRequest, a slog.Attr, latencyFromMs *bool) bool {
	if applyStatus(req, a) {
		return true
	}
	return applyLatencyAttr(req, a, latencyFromMs)
}

// httpExtractState carries the mutable state of a single HTTP extraction pass.
type httpExtractState struct {
	req           *HttpRequest
	pathVal       string
	queryVal      string
	remaining     []slog.Attr
	isHTTP        bool
	latencyFromMs bool
}

func (s *httpExtractState) parseExplicit(a slog.Attr) bool {
	ok, leftover := parseExplicitHTTPAttr(s.req, a)
	if leftover != nil {
		s.remaining = append(s.remaining, *leftover)
	}
	return ok
}

func (s *httpExtractState) tryExplicit(a slog.Attr) bool {
	if !isExplicitHTTPKey(a.Key) {
		return false
	}
	return s.parseExplicit(a)
}

func (s *httpExtractState) tryFlat(a slog.Attr) bool {
	if !s.isHTTP {
		return false
	}
	ctx := httpAttrContext{pathVal: &s.pathVal, queryVal: &s.queryVal, latencyFromMs: &s.latencyFromMs}
	return applyHTTPAttr(s.req, a, &ctx)
}

func (s *httpExtractState) process(a slog.Attr) {
	if s.tryExplicit(a) || s.tryFlat(a) {
		return
	}
	s.remaining = append(s.remaining, a)
}

func extractHTTPRequest(attrs []slog.Attr, isHTTP bool) (*HttpRequest, []slog.Attr) {
	st := &httpExtractState{req: &HttpRequest{}, remaining: make([]slog.Attr, 0, len(attrs)), isHTTP: isHTTP}
	for _, a := range attrs {
		st.process(a)
	}
	finalizeRequestURL(st.req, st.pathVal, st.queryVal)
	if st.req.IsEmpty() {
		return nil, st.remaining
	}
	return st.req, st.remaining
}
