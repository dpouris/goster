package gottp_server

import (
	"errors"
	"log"
	"net/http"
)

type ReqHandler func(http.ResponseWriter, *http.Request) error

type Routes struct {
	route   string
	handler ReqHandler
}

type Gottp struct {
	routes     map[string][]Routes
	middleware []ReqHandler
}

func Server() *Gottp {
	return &Gottp{routes: make(map[string][]Routes, 10), middleware: make([]ReqHandler, 0)}
}

func (g *Gottp) AddGlobalMiddleware(middleware ...ReqHandler) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Gottp) ListenAndServe(port string) {
	LogInfo("LISTENING ON http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, g))
}

func (g *Gottp) Base(method string, addr string, handler ReqHandler) error {
	for _, v := range g.routes[method] {
		if v.route == addr {
			LogError(method + " route already exists")
			return errors.New("route already exists")
		}
	}

	g.routes[method] = append(g.routes[method], Routes{route: addr, handler: handler})

	return nil
}

func (g *Gottp) Get(route string, handler ReqHandler) error {
	err := g.Base("GET", route, handler)
	return err
}

func (g *Gottp) Post(route string, handler ReqHandler) error {
	err := g.Base("POST", route, handler)
	return err
}

func (g *Gottp) Patch(route string, handler ReqHandler) error {
	err := g.Base("PATCH", route, handler)
	return err
}

func (g *Gottp) Put(route string, handler ReqHandler) error {
	err := g.Base("PUT", route, handler)
	return err
}

func (g *Gottp) Delete(route string, handler ReqHandler) error {
	err := g.Base("DELETE", route, handler)
	return err
}

func (g *Gottp) ServeHTTP(r http.ResponseWriter, req *http.Request) {
	m := req.Method
	u := req.URL.String()
	// Middleware that handles validity of incoming request method
	status, err := HandleMethod(g, req)

	if err != nil {
		LogError(err.Error())
		r.WriteHeader(status)
		return
	}

	// Logger middleware
	HandleLog(u, m)

	head := r.Header()
	// Write successful header if all went ok
	MakeSuccessHeader(&head)
	r.WriteHeader(status)

	if len(g.middleware) > 0 {
		for _, m := range g.middleware {
			m(r, req)
		}
	}

	for _, v := range g.routes[m] {
		if stringU := req.URL.String(); stringU == v.route {
			v.handler(r, req)
		}
	}
}
