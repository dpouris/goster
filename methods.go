package goster

import (
	"strings"
)

// Base request. Supply with method (GET, POST, PATCH, PUT, DELETE) to m, URL to u and a handler to h -> error
func (g *Goster) New(m string, u string, h RequestHandler) (route Route) {
	for name := range g.Routes[m] {
		if name == u {
			LogError(m+" route already exists", g.Logger)
			return
		}
	}

	routeType := "normal"
	if strings.Contains(u, ":") {
		routeType = "dynamic"
	}

	parsePath(&u)

	route = Route{Type: routeType, Handler: h}
	g.Routes[m][u] = route

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
