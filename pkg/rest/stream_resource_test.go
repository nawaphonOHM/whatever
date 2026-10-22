package rest

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testResourcePayload = "payload"

var (
	errStreamReadFailed  = errors.New("stream read failed")
	errStreamCloseFailed = errors.New("stream close failed")
)

type spyCloser struct {
	reader io.Reader
	err    error
	closed bool
	reads  int
	closes int
}

func (s *spyCloser) Read(p []byte) (int, error) {
	if s.closed {
		return 0, io.ErrClosedPipe
	}
	s.reads++
	return s.reader.Read(p)
}

func (s *spyCloser) Close() error {
	s.closed = true
	s.closes++
	return s.err
}

type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

func TestStream_CloseOnWrite(t *testing.T) {
	spy := &spyCloser{reader: strings.NewReader(testResourcePayload)}
	c, w := newTestContext()
	Stream(http.StatusOK, spy).Write(c)
	assert.True(t, spy.closed)
	assert.Equal(t, 1, spy.closes)
	assert.Equal(t, testResourcePayload, w.Body.String())
	assert.Empty(t, c.Errors)
}

func TestStream_CloseOnNilContext(t *testing.T) {
	spy := &spyCloser{reader: strings.NewReader(testResourcePayload)}
	Stream(http.StatusOK, spy).Write(nil)
	assert.True(t, spy.closed)
	assert.Equal(t, 1, spy.closes)
}

func TestStream_ReadErrorRecorded(t *testing.T) {
	spy := &spyCloser{reader: errReader{err: errStreamReadFailed}}
	c, _ := newTestContext()
	Stream(http.StatusOK, spy).Write(c)
	assert.True(t, spy.closed)
	require.Len(t, c.Errors, 1)
	assert.Equal(t, errStreamReadFailed, c.Errors[0].Err)
}

func TestStream_CloseErrorRecorded(t *testing.T) {
	spy := &spyCloser{reader: strings.NewReader(testResourcePayload), err: errStreamCloseFailed}
	c, _ := newTestContext()
	Stream(http.StatusOK, spy).Write(c)
	assert.True(t, spy.closed)
	require.Len(t, c.Errors, 1)
	assert.Equal(t, errStreamCloseFailed, c.Errors[0].Err)
}

func TestStream_ReadAndCloseErrorRecorded(t *testing.T) {
	spy := &spyCloser{reader: errReader{err: errStreamReadFailed}, err: errStreamCloseFailed}
	c, _ := newTestContext()
	Stream(http.StatusOK, spy).Write(c)
	assert.True(t, spy.closed)
	require.Len(t, c.Errors, 2)
	assert.Equal(t, errStreamReadFailed, c.Errors[0].Err)
	assert.Equal(t, errStreamCloseFailed, c.Errors[1].Err)
}

func TestStream_NilWriter(t *testing.T) {
	spy := &spyCloser{reader: strings.NewReader(testResourcePayload)}
	c, _ := newTestContext()
	c.Writer = nil
	assert.NotPanics(t, func() { Stream(http.StatusOK, spy).Write(c) })
	assert.True(t, spy.closed)
}
