package gottp_client

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type Gottp struct {
	Routes     map[string]map[string]Route
	Middleware []RequestHandler
	Logger     *log.Logger
	Logs       []string
}

type Route struct {
	Type       string
	Handler    RequestHandler
	identifier string
}

type RequestHandler func(ctx *Ctx) error

type Ctx struct {
	Request        *http.Request
	ResponseWriter Res
	CtxMeta
}

type CtxMeta struct {
	Params
}

type Params struct {
	values map[string]string
}

func (p *Params) Get(id string) string {
	return p.values[id]
}

// New Gottp.Server instance -> *Gottp
func Server() *Gottp {
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

func (g *Gottp) ListenAndServe(p string) {
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

func (g *Gottp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	u := r.URL.String()
	n_res := Res{w}

	route, exists := g.Routes[m][u]

	if !exists {
		for name, route := range g.Routes[m] {
			if route.Type == "dynamic" {
				if new_path, err := match_paths(u, name); err == nil && new_path == u {
					identifier := strings.Trim((parse_url(name)[1]), ":")
					new_route := Route{
						identifier: identifier,
						Type:       "normal",
						Handler:    route.Handler,
					}
					g.Routes[m][new_path] = new_route

					meta := CtxMeta{Params: Params{map[string]string{g.Routes[m][u].identifier: parse_url(u)[1]}}}
					defer route.Handler(&Ctx{ResponseWriter: n_res, Request: r, CtxMeta: meta})
					break
				}
			}
		}
	} else {
		if len(route.identifier) > 0 {
			meta := CtxMeta{Params: Params{map[string]string{g.Routes[m][u].identifier: parse_url(u)[1]}}}
			defer route.Handler(&Ctx{ResponseWriter: n_res, Request: r, CtxMeta: meta})
		} else {
			defer route.Handler(&Ctx{ResponseWriter: n_res, Request: r})
		}

	}

	// Middleware that handles validity of incoming request method
	status, err := HandleMethod(g, r)

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
		m(&Ctx{ResponseWriter: n_res, Request: r})
	}

}

func parse_url(url string) []string {
	r := strings.Split(url, "/")
	temp_r := make([]string, 0)

	for _, v := range r {
		if len(v) != 0 {
			temp_r = append(temp_r, v)
		}
	}

	return temp_r
}

func match_paths(full string, dyn string) (string, error) {
	parsed_full := parse_url(full)
	parsed_dyn := parse_url(dyn)

	if len(parsed_full) <= 1 || len(parsed_dyn) <= 1 {
		return "", errors.New("not matching url's")
	}
	p := strings.Replace(dyn, parsed_dyn[1], parsed_full[1], -1)
	return p, nil
}
