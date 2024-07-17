package goster

import (
	"strings"
)

// Base request. Supply with method (GET, POST, PATCH, PUT, DELETE) to m, URL to u and a handler to h -> error
func (g *Goster) New(method string, url string, handler RequestHandler) (route Route) {
	for name := range g.Routes[method] {
		if name == url {
			LogError(method+" route already exists", g.Logger)
			return
		}
	}

	routeType := "normal"
	if strings.Contains(url, ":") {
		routeType = "dynamic"
	}

	cleanPath(&url)

	route = Route{Type: routeType, Handler: handler}
	g.Routes[method][url] = route

	return
}

// Make a Get request to route and pass in a ReqHandler function -> error
func (g *Goster) Get(route string, handler RequestHandler) (r Route) {
	r = g.New("GET", route, handler)
	return
}

// Make a Post request to route and pass in a ReqHandler function -> error
func (g *Goster) Post(route string, handler RequestHandler) (r Route) {
	r = g.New("POST", route, handler)
	return
}

// Make a Patch request to route and pass in a ReqHandler function -> error
func (g *Goster) Patch(route string, handler RequestHandler) (r Route) {
	r = g.New("PATCH", route, handler)
	return
}

// Make a Put request to route and pass in a ReqHandler function -> error
func (g *Goster) Put(route string, handler RequestHandler) (r Route) {
	r = g.New("PUT", route, handler)
	return
}

// Make a Delete request to route and pass in a ReqHandler function -> error
func (g *Goster) Delete(route string, handler RequestHandler) (r Route) {
	r = g.New("DELETE", route, handler)
	return
}
