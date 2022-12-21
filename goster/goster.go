package goster

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Goster struct {
	Context    *Ctx
	Routes     map[string]map[string]Route
	Middleware map[string][]RequestHandler
	Logger     *log.Logger
	Logs       []string
}

type Route struct {
	Type         string
	Handler      RequestHandler
	DynamicRoute DynamicRoute
}

type RequestHandler func(ctx *Ctx) error

type DynamicRoute struct {
	FullPath        string
	DynPath         string
	Identifier      string
	IdentifierValue string
}

type Ctx struct {
	Request  *http.Request
	Response Response
	Meta
}

type Meta struct {
	Params
}

type Params struct {
	values map[string]string
}

// --------------------------------------------------------------------------------------------- //

var e engine

// New Goster.NewServer instance -> *Goster
func NewServer() *Goster {
	g := e.Init()
	return g
}

// Pass in a ReqHandler or ...ReqHandler type function(s) to handle incoming http requests on every single request
func (g *Goster) UseGlobal(m ...RequestHandler) {
	g.Middleware["*"] = append(g.Middleware["*"], m...)
}

func (g *Goster) Use(path string, m ...RequestHandler) {
	parsePath(&path)
	g.Middleware[path] = m
	// for method, routes := range g.Routes {
	// 	if _, exists := g.Routes[method][path]; exists {
	// 		route := routes[path]
	// 		route.Middleware = append(route.Middleware, m...)
	// 		routes[path] = route
	// 	}
	// }
}

// Start listening for incoming requests
func (g *Goster) ListenAndServe(p string) {
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

// In order to inherit the Handler interface, we must include a method called ServeHTTP. This is where the magic happens
func (g *Goster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Ctx{
		Request:  r,
		Response: Response{w},
	}

	g.Context = &c

	defer DefaultHeader(&c)

	g.launchHandler()

	for _, rh := range g.Middleware["*"] {
		rh(&c)
	}

}

// Launches the necessary handler for the incoming request
func (g *Goster) launchHandler() {
	c := g.Context
	m := c.Request.Method
	u := c.Request.URL.String()

	// parses query params if there are any and removes them from url effectively transforming it from /path/path2?q=v -> /path/path2
	params, err := parseParams(&u)
	meta := Meta{}

	if err == nil {
		meta.Params = params
	} else {
		meta.Params.values = map[string]string{}
	}

	// we need to pass the new stripped path in handleRoute as parseParams has cleaned it and will match the routes that are in g.Routes
	err = g.handleRoute(c, u)

	HandleLog(c, g, err)

	route, routeExists := g.Routes[m][u]

	if !routeExists {
		for name, route := range g.Routes[m] {
			if route.Type != "dynamic" {
				continue
			}

			matchedRoute, err := matchDynamicRoute(u, name)
			if err != nil {
				fmt.Fprintln(e.Goster.Logger.Writer(), fmt.Errorf("error: %s", err.Error()))
				continue
			}

			newRoute := Route{
				Type:         "normal",
				Handler:      route.Handler,
				DynamicRoute: matchedRoute,
			}
			e.Goster.Routes[m][matchedRoute.FullPath] = newRoute

			meta.Params.values[e.Goster.Routes[m][u].DynamicRoute.Identifier] = matchedRoute.IdentifierValue
			c.Meta = meta
			defer route.Handler(c)

			break
		}
	} else {
		if len(route.DynamicRoute.Identifier) > 0 {
			meta.Params.values[route.DynamicRoute.Identifier] = g.Routes[m][u].DynamicRoute.IdentifierValue
		}
		c.Meta = meta
		defer route.Handler(c)
	}

	for _, rh := range g.Middleware[u] {
		rh(c)
	}
}

func (g *Goster) handleRoute(c *Ctx, u string) (err error) {
	m := c.Request.Method
	methodAllowed := false
	routeExists := false

	for method := range g.Routes {
		if _, exists := g.Routes[method][u]; exists && method == m {
			routeExists = true
			methodAllowed = true
			break
		} else if exists && method != m {
			routeExists = true
		}
	}

	if !routeExists {
		c.Response.WriteHeader(http.StatusNotFound)
		err = errors.New("404 not found")
		return
	}

	if !methodAllowed {
		c.Response.WriteHeader(http.StatusMethodNotAllowed)
		err = errors.New("405 method not allowed")
		return
	}
	return
}

// Get tries to find if id is passed in to the url as a query param or as a dynamic route. If the specified id isn't found <e> will be false
func (p *Params) Get(id string) (i string, e bool) {
	id, exists := p.values[id]

	return id, exists
}
