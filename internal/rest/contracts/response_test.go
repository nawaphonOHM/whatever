package contracts

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stringerVal struct{ text string }

func (s stringerVal) String() string { return s.text }

type stringerPtr struct{ text string }

func (s *stringerPtr) String() string { return s.text }

type panicStringer struct{}

func (panicStringer) String() string { panic("stringer exploded") }

func TestSafeStringer(t *testing.T) {
	t.Run("valid struct value", func(t *testing.T) {
		str, ok := SafeStringer(stringerVal{text: "val"})
		assert.True(t, ok)
		assert.Equal(t, "val", str)
	})

	t.Run("valid pointer value", func(t *testing.T) {
		str, ok := SafeStringer(&stringerPtr{text: "ptr"})
		assert.True(t, ok)
		assert.Equal(t, "ptr", str)
	})

	t.Run("nil interface", func(t *testing.T) {
		str, ok := SafeStringer(nil)
		assert.False(t, ok)
		assert.Empty(t, str)
	})

	t.Run("typed nil pointer", func(t *testing.T) {
		var s *stringerPtr
		str, ok := SafeStringer(s)
		assert.False(t, ok)
		assert.Empty(t, str)
	})

	t.Run("panicking stringer propagates panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = SafeStringer(panicStringer{})
		})
	})
}

func TestSetHeader(t *testing.T) {
	t.Run("empty key is ignored", func(t *testing.T) {
		var h map[string]string
		SetHeader(&h, "", "val")
		assert.Nil(t, h)
	})

	t.Run("canonicalizes keys and overwrites deterministic", func(t *testing.T) {
		var h map[string]string
		SetHeader(&h, "x-custom-header", "v1")
		assert.Equal(t, "v1", h["X-Custom-Header"])

		SetHeader(&h, "X-CUSTOM-HEADER", "v2")
		assert.Len(t, h, 1)
		assert.Equal(t, "v2", h["X-Custom-Header"])
	})
}

func TestWriteHeaders(t *testing.T) {
	t.Run("nil context or writer", func(t *testing.T) {
		assert.NotPanics(t, func() {
			WriteHeaders(nil, map[string]string{"K": "V"})
			WriteHeaders(&gin.Context{}, map[string]string{"K": "V"})
		})
	})

	t.Run("empty headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		WriteHeaders(c, nil)
		WriteHeaders(c, make(map[string]string))
		assert.Empty(t, w.Header())
	})

	t.Run("valid and empty headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		headers := map[string]string{
			"X-Custom": "value",
			"X-Empty":  "",
			"":         "invalid-key",
		}
		WriteHeaders(c, headers)
		assert.Equal(t, "value", w.Header().Get("X-Custom"))
		assert.Equal(t, "", w.Header().Get("X-Empty"))
		assert.Contains(t, w.Header(), "X-Empty")
		assert.NotContains(t, w.Header(), "")
	})
}
