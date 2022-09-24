package gottp_client

import "errors"

// Base request. Supply with method (GET, POST, PATCH, PUT, DELETE) to m, URL to u and a handler to h -> error
func (g *Gottp) New(m string, u string, h ReqHandler) error {
	for _, v := range g.Routes[m] {
		if v.Route == u {
			LogError(m+" route already exists", g.Logger)
			return errors.New("route already exists")
		}
	}

	g.Routes[m] = append(g.Routes[m], Routes{Route: u, Handler: h})

	return nil
}

// Make a Get request to route and pass in a ReqHandler function -> error
func (g *Gottp) Get(route string, handler ReqHandler) error {
	err := g.New("GET", route, handler)
	return err
}

// Make a Post request to route and pass in a ReqHandler function -> error
func (g *Gottp) Post(route string, handler ReqHandler) error {
	err := g.New("POST", route, handler)
	return err
}

// Make a Patch request to route and pass in a ReqHandler function -> error
func (g *Gottp) Patch(route string, handler ReqHandler) error {
	err := g.New("PATCH", route, handler)
	return err
}

// Make a Put request to route and pass in a ReqHandler function -> error
func (g *Gottp) Put(route string, handler ReqHandler) error {
	err := g.New("PUT", route, handler)
	return err
}

// Make a Delete request to route and pass in a ReqHandler function -> error
func (g *Gottp) Delete(route string, handler ReqHandler) error {
	err := g.New("DELETE", route, handler)
	return err
}
