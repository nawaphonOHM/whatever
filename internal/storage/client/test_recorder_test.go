package client_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const recordedBodyLimit = 1 << 20

type recordedRequest struct {
	ReadErr       error
	Header        http.Header
	Query         url.Values
	Method        string
	Scheme        string
	Host          string
	Path          string
	Body          []byte
	ContentLength int64
	BodyLen       int64
}

type headWriter struct {
	buf bytes.Buffer
	n   int64
}

func (w *headWriter) Write(p []byte) (int, error) {
	w.n += int64(len(p))
	room := recordedBodyLimit - w.buf.Len()
	if room <= 0 {
		return len(p), nil
	}
	_, err := w.buf.Write(p[:min(room, len(p))])
	return len(p), err
}

func newRecordedRequest(req *http.Request) recordedRequest {
	rec := recordedRequest{
		Header: req.Header.Clone(), Query: req.URL.Query(), Method: req.Method,
		Scheme: req.URL.Scheme, Host: req.URL.Host, Path: req.URL.Path, ContentLength: req.ContentLength,
	}
	if req.Body != nil {
		w := &headWriter{}
		_, rec.ReadErr = io.Copy(w, req.Body)
		rec.Body, rec.BodyLen = w.buf.Bytes(), w.n
	}
	return rec
}

// recorder is a mock S3 transport that records every request and answers through respond.
type recorder struct {
	respond func(recordedRequest) *http.Response
	reqs    []recordedRequest
	mu      sync.Mutex
}

func (r *recorder) roundTrip(req *http.Request) (*http.Response, error) {
	rec := newRecordedRequest(req)
	r.mu.Lock()
	r.reqs = append(r.reqs, rec)
	r.mu.Unlock()
	return r.respond(rec), nil
}

func (r *recorder) all() []recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedRequest(nil), r.reqs...)
}

func (r *recorder) httpClient() *http.Client {
	return newMockHTTPClient(r.roundTrip)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func statusResponder(status int) func(recordedRequest) *http.Response {
	return func(recordedRequest) *http.Response {
		return newHTTPResponse(status, "", nil)
	}
}
