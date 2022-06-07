package gottp_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type ReqHandler func(http.ResponseWriter, *http.Request) error

type Req struct {
	Method        string
	URL           *url.URL
	Header        http.Header
	Body          io.ReadCloser
	GetBody       func() (io.ReadCloser, error)
	ContentLength int64
	Close         bool
	Form          url.Values
	PostForm      url.Values
	MultipartForm *multipart.Form
	RemoteAddr    string
	RequestURI    string
	Response      *http.Response
}

type Res struct {
	Response       http.Response
	ResponseWriter http.ResponseWriter
	// Header() http.Header

	// Write([]byte) (int, error)

	// WriteHeader(statusCode int)
}

func (r *Res) JSON(j any) error {
	v, err := json.Marshal(r)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}

	_, err = r.ResponseWriter.Write(v)

	return err
}

func (r *Res) Write(b []byte) (int, error) {
	br, err := r.Response.Body.Read(b)

	fmt.Println(string(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 0, err
	}

	return br, nil
}

type Routes struct {
	Route   string
	Handler ReqHandler
}

type Gottp struct {
	Routes     map[string][]Routes
	Middleware []ReqHandler
	Logger     *log.Logger
	Logs       []string
}

func Server() *Gottp {
	logger := log.New(os.Stdout, "[SERVER] - ", log.LstdFlags)
	return &Gottp{Routes: make(map[string][]Routes, 10), Middleware: make([]ReqHandler, 0), Logger: logger}
}

func (g *Gottp) AddGlobalMiddleware(middleware ...ReqHandler) {
	g.Middleware = append(g.Middleware, middleware...)
}

func (g *Gottp) ListenAndServe(port string) {
	LogInfo("LISTENING ON http://127.0.0.1"+port, g.Logger)
	log.Fatal(http.ListenAndServe(port, g))
}

func (g *Gottp) Base(method string, addr string, handler ReqHandler) error {
	for _, v := range g.Routes[method] {
		if v.Route == addr {
			LogError(method+" route already exists", g.Logger)
			return errors.New("route already exists")
		}
	}

	g.Routes[method] = append(g.Routes[method], Routes{Route: addr, Handler: handler})

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

	// Logger middleware
	HandleLog(u, m, err, g)

	if err != nil {
		r.WriteHeader(status)
		return
	}

	// Write successful header if all went ok
	head := r.Header()
	DefaultHeader(&head)

	if len(g.Middleware) > 0 {
		for _, m := range g.Middleware {
			m(r, req)
		}
	}

	for _, v := range g.Routes[m] {
		if stringU := req.URL.String(); stringU == v.Route {
			v.Handler(r, req)
		}
	}
}
