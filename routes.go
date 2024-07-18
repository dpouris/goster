package goster

import (
	"fmt"
	"strings"
)

type Routes map[string]map[string]Route

// New creates a new Route for the specified method and url using the provided handler. If the Route already exists an error is returned.
func (rs *Routes) New(method string, url string, handler RequestHandler) (err error) {
	for name := range (*rs)[method] {
		if name == url {
			err = fmt.Errorf("[%s] -> [%s] route already exists", method, url)
			return
		}
	}

	routeType := "normal"
	if strings.Contains(url, ":") {
		routeType = "dynamic"
	}

	cleanPath(&url)

	(*rs)[method][url] = Route{Type: routeType, Handler: handler}

	return
}

// Get creates a new Route under the GET method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Get(url string, handler RequestHandler) error {
	return g.Routes.New("GET", url, handler)
}

// Post creates a new Route under the POST method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Post(path string, handler RequestHandler) error {
	return g.Routes.New("POST", path, handler)
}

// Patch creates a new Route under the PATCH method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Patch(path string, handler RequestHandler) error {
	return g.Routes.New("PATCH", path, handler)
}

// Put creates a new Route under the PUT method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Put(path string, handler RequestHandler) error {
	return g.Routes.New("PUT", path, handler)
}

// Delete creates a new Route under the DELETE method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Delete(path string, handler RequestHandler) error {
	return g.Routes.New("DELETE", path, handler)
}
