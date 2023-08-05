package gfx

// RouteGroup represents a group of routes
type RouteGroup struct {
	engine     *GFXEngine
	basePath   string
	middleware []MiddlewareFunc
}

// UseMiddleware adds a middleware to the group
func (rg *RouteGroup) UseMiddleware(middleware MiddlewareFunc) {
	rg.middleware = append(rg.middleware, middleware)
}

// Get adds a GET route to the group
func (rg *RouteGroup) Get(path string, handler HandlerFunc) {
	rg.engine.addRoute("GET", path, handler, nil, rg)
}

// Post adds a POST route to the group
func (rg *RouteGroup) Post(path string, handler HandlerFunc) {
	rg.engine.addRoute("POST", path, handler, nil, rg)
}
