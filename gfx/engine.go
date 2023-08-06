package gfx

import (
	"net/http"
	"strings"
)

// GFXEngine represents the main engine
type GFXEngine struct {
	routes []Route
}

// NewGFXEngine creates a new GFXEngine
func NewGFXEngine() *GFXEngine {
	return &GFXEngine{}
}

// ServeHTTP serves HTTP requests
func (g *GFXEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestParts := strings.Split(r.URL.Path, "/")

	for _, route := range g.routes {
		if r.Method == route.method && len(requestParts) == len(route.parts) {
			if c := g.processRoute(route, w, r, requestParts); c != nil {
				c.Next()
				return
			}
		}
	}

	http.NotFound(w, r)
}

// Get adds a GET route to the engine
func (g *GFXEngine) Get(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler, nil, nil)
}

// Post adds a POST route to the engine
func (g *GFXEngine) Post(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler, nil, nil)
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
