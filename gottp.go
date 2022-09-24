package gottp_client

import (
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

func (g *Gottp) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	m := r.Method
	u := r.URL.String()
	n_res := Res{res}

	// Middleware that handles validity of incoming request method
	status, err := HandleMethod(g, r)

	// Logger middleware
	HandleLog(u, m, err, g)

	// Transform the ResponseWriter and Request params to be more manageable by end users and adds some useful function middleware
	n_req := TransformReq(r)

	if err != nil {
		res.WriteHeader(status)
		return
	}

	// Write successful header if all went ok
	head := res.Header()
	DefaultHeader(&head)

	for _, m := range g.Middleware {
		m(n_res, &n_req)

	}

	for _, v := range g.Routes[m] {
		if stringU := r.URL.String(); stringU == v.Route {
			v.Handler(n_res, &n_req)
		}
	}
}
