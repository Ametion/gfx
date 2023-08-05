package gfx

import (
	"encoding/json"
	"net/http"
)

// Context represents request context
type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Headers    http.Header
	params     map[string]string
	index      int
	middleware []MiddlewareFunc
}

// Next proceeds to the next middleware
func (c *Context) Next() {
	if c.index < len(c.middleware) {
		middleware := c.middleware[c.index]
		c.index++
		middleware(c)
	}
}

// Query gets a query value
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Param gets a path parameter
func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) SetBody(v interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	defer c.Request.Body.Close()
	return decoder.Decode(v)
}

// SendJSON sends a SendJSON response
func (c *Context) SendJSON(statusCode int, v interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	json.NewEncoder(c.Writer).Encode(v)
}
