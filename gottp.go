package gottp_client

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Gottp struct {
	Routes     map[string]map[string]Route
	Middleware []RequestHandler
	Logger     *log.Logger
	Logs       []string
}

type Route struct {
	Type         string
	Handler      RequestHandler
	DynamicRoute DynamicRoute `json: omitempty`
}

type RequestHandler func(ctx *Ctx) error

type DynamicRoute struct {
	FullPath        string
	DynPath         string
	Identifier      string
	IdentifierValue string
}

type Ctx struct {
	Request        *http.Request
	ResponseWriter Response
	Meta
}

type Meta struct {
	Params
}

type Params struct {
	values map[string]string
}

// Get tries to find if id is passed in to the url as a query param or as a dynamic route. If the specified id isn't found <e> will be false
func (p *Params) Get(id string) (i string, e bool) {
	id, exists := p.values[id]

	return id, exists
}

// New Gottp.NewServer instance -> *Gottp
func NewServer() *Gottp {
	logger := log.New(os.Stdout, "[SERVER] - ", log.LstdFlags)
	methods := make(map[string]map[string]Route)
	methods["GET"] = make(map[string]Route)
	methods["POST"] = make(map[string]Route)
	methods["PUT"] = make(map[string]Route)
	methods["PATCH"] = make(map[string]Route)
	methods["DELETE"] = make(map[string]Route)
	return &Gottp{Routes: methods, Middleware: make([]RequestHandler, 0), Logger: logger}
}

// Pass in a ReqHandler or ...ReqHandler type function(s) to handle incoming http requests on every single request
func (g *Gottp) AddGlobalMiddleware(m ...RequestHandler) {
	g.Middleware = append(g.Middleware, m...)
}

// Start listening for incoming requests
func (g *Gottp) ListenAndServe(p string) {
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

func (g *Gottp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	u := r.URL.String()
	res := Response{w}

	// parses query params if there are any and removes them from url effectively transforming it from /path/path2?q=v -> /path/path2
	params, err := parseParams(&u)
	meta := Meta{}

	if err == nil {
		meta.Params = params
	} else {
		meta.Params.values = map[string]string{}
	}

	route, routeExists := g.Routes[m][u]

	if !routeExists {
		for name, route := range g.Routes[m] {
			if route.Type != "dynamic" {
				continue
			}

			matchedRoute, err := matchDynamicRoute(u, name)
			if err != nil {
				fmt.Fprint(g.Logger.Writer(), fmt.Errorf("error: %s", err.Error()))
				continue
			}

			newRoute := Route{
				Type:         "normal",
				Handler:      route.Handler,
				DynamicRoute: matchedRoute,
			}
			g.Routes[m][matchedRoute.FullPath] = newRoute
			meta.Params.values[g.Routes[m][u].DynamicRoute.Identifier] = matchedRoute.IdentifierValue

			defer route.Handler(&Ctx{ResponseWriter: res, Request: r, Meta: meta})
			break
		}
	} else {
		if len(route.DynamicRoute.Identifier) > 0 {
			meta.Params.values[route.DynamicRoute.Identifier] = g.Routes[m][u].DynamicRoute.IdentifierValue
		}
		defer route.Handler(&Ctx{ResponseWriter: res, Request: r, Meta: meta})
	}

	// Middleware that handles validity of incoming request method
	status, err := HandleMethod(g, u, m)

	// Logger middleware
	HandleLog(u, m, err, g)

	if err != nil {
		w.WriteHeader(status)
		return
	}

	// Write successful header if all went ok
	head := w.Header()
	DefaultHeader(&head)

	for _, m := range g.Middleware {
		m(&Ctx{ResponseWriter: res, Request: r})
	}

}
