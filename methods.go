package gottp_client

import (
	"errors"
	"strings"
)

// Base request. Supply with method (GET, POST, PATCH, PUT, DELETE) to m, URL to u and a handler to h -> error
func (g *Gottp) New(m string, u string, h RequestHandler) error {
	for name := range g.Routes[m] {
		if name == u {
			LogError(m+" route already exists", g.Logger)
			return errors.New("route already exists")
		}
	}

	url_type := "normal"
	if strings.Contains(u, ":") {
		url_type = "dynamic"
	}

	if u[0] != '/' {
		u = "/" + u
	}

	if u[len(u)-1] == '/' {
		u = u[:len(u)-1]
	}

	g.Routes[m][u] = Route{Type: url_type, Handler: h}

	return nil
}

// Make a Get request to route and pass in a ReqHandler function -> error
func (g *Gottp) Get(route string, handler RequestHandler) error {
	err := g.New("GET", route, handler)
	return err
}

// Make a Post request to route and pass in a ReqHandler function -> error
func (g *Gottp) Post(route string, handler RequestHandler) error {
	err := g.New("POST", route, handler)
	return err
}

// Make a Patch request to route and pass in a ReqHandler function -> error
func (g *Gottp) Patch(route string, handler RequestHandler) error {
	err := g.New("PATCH", route, handler)
	return err
}

// Make a Put request to route and pass in a ReqHandler function -> error
func (g *Gottp) Put(route string, handler RequestHandler) error {
	err := g.New("PUT", route, handler)
	return err
}

// Make a Delete request to route and pass in a ReqHandler function -> error
func (g *Gottp) Delete(route string, handler RequestHandler) error {
	err := g.New("DELETE", route, handler)
	return err
}
