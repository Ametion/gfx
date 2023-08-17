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

// Patch adds a PATCH route to the engine
func (rg *RouteGroup) Patch(path string, handler HandlerFunc) {
	rg.engine.addRoute("PATCH", path, handler, nil, nil)
}

// Put adds a PUT route to the engine
func (rg *RouteGroup) Put(path string, handler HandlerFunc) {
	rg.engine.addRoute("PUT", path, handler, nil, nil)
}

// Delete adds a DELETE route to the engine
func (rg *RouteGroup) Delete(path string, handler HandlerFunc) {
	rg.engine.addRoute("DELETE", path, handler, nil, nil)
}
