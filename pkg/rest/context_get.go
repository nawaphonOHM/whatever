package rest

import "time"

// GetString returns a context value as a string.
func (c *Context) GetString(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.GetString(key)
	}
	return ""
}

// GetBool returns a context value as a boolean.
func (c *Context) GetBool(key string) bool {
	if c.ginCtx != nil {
		return c.ginCtx.GetBool(key)
	}
	return false
}

// GetInt returns a context value as an integer.
func (c *Context) GetInt(key string) int {
	if c.ginCtx != nil {
		return c.ginCtx.GetInt(key)
	}
	return 0
}

// GetInt64 returns a context value as int64.
func (c *Context) GetInt64(key string) int64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetInt64(key)
	}
	return 0
}

// GetUint returns a context value as uint.
func (c *Context) GetUint(key string) uint {
	if c.ginCtx != nil {
		return c.ginCtx.GetUint(key)
	}
	return 0
}

// GetUint64 returns a context value as uint64.
func (c *Context) GetUint64(key string) uint64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetUint64(key)
	}
	return 0
}

// GetFloat64 returns a context value as float64.
func (c *Context) GetFloat64(key string) float64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetFloat64(key)
	}
	return 0
}

// GetTime returns a context value as time.Time.
func (c *Context) GetTime(key string) time.Time {
	if c.ginCtx != nil {
		return c.ginCtx.GetTime(key)
	}
	return time.Time{}
}

// GetDuration returns a context value as time.Duration.
func (c *Context) GetDuration(key string) time.Duration {
	if c.ginCtx != nil {
		return c.ginCtx.GetDuration(key)
	}
	return 0
}

// GetStringSlice returns a context value as a string slice.
func (c *Context) GetStringSlice(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringSlice(key)
	}
	return nil
}

// GetStringMap returns a context value as a map of interfaces.
func (c *Context) GetStringMap(key string) map[string]any {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMap(key)
	}
	return nil
}

// GetStringMapString returns a context value as a string map.
func (c *Context) GetStringMapString(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMapString(key)
	}
	return nil
}

// GetStringMapStringSlice returns a context value as a map of string slices.
func (c *Context) GetStringMapStringSlice(key string) map[string][]string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMapStringSlice(key)
	}
	return nil
}
