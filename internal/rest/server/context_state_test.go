package server

import (
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
)

func TestContext_StateStore(t *testing.T) {
	gc, _ := newTestContext()
	c := newContext(gc)
	now := time.Now().UTC().Truncate(time.Second)
	setBasicState(c, now)
	value, exists := c.Get(testKey)
	assert.True(t, exists)
	assertBasicState(t, c, value, now)
}

func setBasicState(c *Context, now time.Time) {
	c.Set(testKey, testValue)
	c.Set("bool", true)
	c.Set("int", testIntValue)
	c.Set("time", now)
}

func assertBasicState(t *testing.T, c rest.Context, value any, now time.Time) {
	assert.Equal(t, testValue, value)
	assert.Equal(t, testValue, c.MustGet(testKey))
	assert.True(t, c.GetBool("bool"))
	assert.Equal(t, testIntValue, c.GetInt("int"))
	assert.Equal(t, now, c.GetTime("time"))
}

func TestContext_TypedStateValues(t *testing.T) {
	gc, _ := newTestContext()
	c := newContext(gc)
	setTypedState(c)
	assertTypedState(t, c)
}

func setTypedState(c *Context) {
	c.Set("int64", int64(testInt64Value))
	c.Set("uint", uint(testUintValue))
	c.Set("uint64", uint64(testUint64Value))
	c.Set("float", testFloatValue)
	c.Set("duration", time.Minute)
	c.Set("slice", []string{testValue})
}

func assertTypedState(t *testing.T, c rest.Context) {
	assert.Equal(t, int64(testInt64Value), c.GetInt64("int64"))
	assert.Equal(t, uint(testUintValue), c.GetUint("uint"))
	assert.Equal(t, uint64(testUint64Value), c.GetUint64("uint64"))
	assert.Equal(t, testFloatValue, c.GetFloat64("float"))
	assert.Equal(t, time.Minute, c.GetDuration("duration"))
	assert.Equal(t, []string{testValue}, c.GetStringSlice("slice"))
}

func TestContext_StateMapsAndMissingValues(t *testing.T) {
	gc, _ := newTestContext()
	c := newContext(gc)
	setMapState(c)
	assert.Equal(t, map[string]any{testKey: testValue}, c.GetStringMap("map"))
	assertMapState(t, c)
	_, exists := c.Get("missing")
	assert.False(t, exists)
	assert.Panics(t, func() { c.MustGet("missing") })
}

func assertMapState(t *testing.T, c rest.Context) {
	stringsMap := map[string]string{testKey: testValue}
	slicesMap := map[string][]string{testKey: {testValue}}
	assert.Equal(t, stringsMap, c.GetStringMapString("strings"))
	assert.Equal(t, slicesMap, c.GetStringMapStringSlice("slices"))
}

func setMapState(c *Context) {
	c.Set("map", map[string]any{testKey: testValue})
	c.Set("strings", map[string]string{testKey: testValue})
	c.Set("slices", map[string][]string{testKey: {testValue}})
}
