package gottp_server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type ReqHandler func(Res, *Req) error

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
	Header func() http.Header

	Write func([]byte) (int, error)

	WriteHeader func(statusCode int)
}

// Send back a JSON response. Supply j with a value that's valid marsallable(?) to JSON -> error
func (r *Res) JSON(j any) error {
	v, err := json.Marshal(j)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}
	r.Header().Set("Content-Type", "application/json")
	_, err = r.Write(v)

	return err
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

// New Gottp.Server instance -> *Gottp
func Server() *Gottp {
	logger := log.New(os.Stdout, "[SERVER] - ", log.LstdFlags)
	return &Gottp{Routes: make(map[string][]Routes, 10), Middleware: make([]ReqHandler, 0), Logger: logger}
}

// Pass in a ReqHandler or ...ReqHandler type function(s) to handle incoming http requests on every single request
func (g *Gottp) AddGlobalMiddleware(m ...ReqHandler) {
	g.Middleware = append(g.Middleware, m...)
}

func (g *Gottp) ListenAndServe(p string) {
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

func (g *Gottp) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m := req.Method
	u := req.URL.String()

	// Middleware that handles validity of incoming request method
	status, err := HandleMethod(g, req)

	// Logger middleware
	HandleLog(u, m, err, g)

	// Transform the ResponseWriter and Request params to be more manageable by end users and adds some useful function middleware
	n_res, n_req := TransformReq(res, req)

	if err != nil {
		n_res.WriteHeader(status)
		return
	}

	// Write successful header if all went ok
	head := n_res.Header()
	DefaultHeader(&head)

	if len(g.Middleware) > 0 {
		for _, m := range g.Middleware {
			m(n_res, &n_req)
		}
	}

	for _, v := range g.Routes[m] {
		if stringU := req.URL.String(); stringU == v.Route {
			v.Handler(n_res, &n_req)
		}
	}
}
