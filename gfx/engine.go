package gfx

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GFXEngine represents the main engine
type GFXEngine struct {
	routes      []Route
	development bool
}

// NewGFXEngine creates a new GFXEngine
func NewGFXEngine() *GFXEngine {
	return &GFXEngine{}
}

func (g *GFXEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now() // Record the start time
	requestParts := strings.Split(r.URL.Path, "/")
	statusCode := http.StatusNotFound // Default status code
	wrappedWriter := newStatusRecorder(w, &statusCode)

	for _, route := range g.routes {
		if r.Method == route.method && len(requestParts) == len(route.parts) {
			if c := g.processRoute(route, wrappedWriter, r, requestParts); c != nil {
				c.Next()
				statusCode = http.StatusOK // Success status code
				break
			}
		}
	}

	if statusCode == http.StatusNotFound {
		http.NotFound(wrappedWriter, r)
	}

	if g.development {
		// Log details to the console if development mode is on
		fmt.Printf("Date: %s, Method: %s, Status code: %d, Time taken: %v\n", time.Now().Format(time.RFC1123), r.Method, statusCode, time.Since(startTime))
	}
}

func (g *GFXEngine) IsDevelopment() {
	g.development = true
}

// Get adds a GET route to the engine
func (g *GFXEngine) Get(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler, nil, nil)
}

// Post adds a POST route to the engine
func (g *GFXEngine) Post(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler, nil, nil)
}

// Patch adds a PATCH route to the engine
func (g *GFXEngine) Patch(path string, handler HandlerFunc) {
	g.addRoute("PATCH", path, handler, nil, nil)
}

// Put adds a PUT route to the engine
func (g *GFXEngine) Put(path string, handler HandlerFunc) {
	g.addRoute("PUT", path, handler, nil, nil)
}

// Delete adds a DELETE route to the engine
func (g *GFXEngine) Delete(path string, handler HandlerFunc) {
	g.addRoute("DELETE", path, handler, nil, nil)
}

// Group creates a new RouteGroup
func (g *GFXEngine) Group(basePath string) *RouteGroup {
	return &RouteGroup{
		engine:   g,
		basePath: basePath,
	}
}

// Run starts the web
func (g *GFXEngine) Run(addr string) error {
	return http.ListenAndServe(addr, g)
}

func (g *GFXEngine) processRoute(route Route, w http.ResponseWriter, r *http.Request, requestParts []string) *Context {
	params := make(map[string]string)
	match := true

	for _, i := range route.paramsIndex {
		params[route.parts[i]] = requestParts[i]
	}

	for i, part := range requestParts {
		if i >= len(route.parts) || (part != route.parts[i] && !contains(route.paramsIndex, i)) {
			match = false
			break
		}
	}

	if match {
		c := &Context{
			Writer:     w,
			Request:    r,
			params:     params,
			Headers:    r.Header,
			index:      0,
			middleware: append(route.middleware, handlerToMiddleware(route.handler)),
		}
		return c
	}

	return nil
}

// addRoute adds a route to the engine
func (g *GFXEngine) addRoute(method string, path string, handler HandlerFunc, middleware []MiddlewareFunc, group *RouteGroup) {
	if group != nil {
		path = group.basePath + path
		middleware = append(group.middleware, middleware...)
	}

	parts := strings.Split(path, "/")
	paramsIndex := []int{}

	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramsIndex = append(paramsIndex, i)
			parts[i] = part[1:] // Remove the ":" prefix
		}
	}

	route := Route{method, handler, middleware, parts, paramsIndex}
	g.routes = append(g.routes, route)
}

func handlerToMiddleware(h HandlerFunc) MiddlewareFunc {
	return func(c *Context) {
		h(c)
	}
}

func contains(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
