package goster

import (
	"fmt"
	"log"
	"net/http"
)

const (
	TypeStatic   = "static"
	TypeDynamic  = "dynamic"
	TypeWildcard = "wildcard"
)

// Goster is the main structure of the package. It handles the addition of new routes and middleware, and manages logging.
type Goster struct {
	Routes     Routes                      // Routes is a map of HTTP methods to their respective route handlers.
	Middleware map[string][]RequestHandler // Middleware is a map of routes to their respective middleware handlers.
	Logger     *log.Logger                 // Logger is used for logging information and errors.
	Logs       []string                    // Logs stores logs for future reference.
}

// RequestHandler is a type for functions that handle HTTP requests within a given context.
type RequestHandler func(ctx *Ctx) error

// ------------------------------------------Public Methods--------------------------------------------------- //

// NewServer creates a new Goster instance.
func NewServer() *Goster {
	g := engine.init()
	return g
}

// UseGlobal adds middleware handlers that will be applied to every single incoming request.
func (g *Goster) UseGlobal(m ...RequestHandler) {
	g.Middleware["*"] = append(g.Middleware["*"], m...)
}

// Use adds middleware handlers that will be applied to specific routes/paths.
func (g *Goster) Use(path string, m ...RequestHandler) {
	cleanPath(&path)
	g.Middleware[path] = m
}

// TemplateDir extends the engine's file paths with the specified directory `d`,
// which is joined to Engine.Config.BaseStaticDir (default is the execution path of the program).
//
// This instructs the engine where to look for template files like .html, .gohtml.
// If the directory doesn't exist, it will return an appropriate error.
func (g *Goster) TemplateDir(d string) (err error) {
	err = engine.SetTemplateDir(d)

	if err != nil {
		LogError(err.Error(), g.Logger)
	}

	return
}

// StaticDir sets the directory from which static files like .css, .js, etc are served.
//
// It integrates the specified directory into the server's static file handling
// by invoking AddStaticDir on the Routes collection.
//
// If an error occurs during this process, the error is printed to the standard error output.
// The function returns the error encountered, if any.
func (g *Goster) StaticDir(dir string) (err error) {
	err = engine.SetStaticDir(dir)
	if err != nil {
		LogError(err.Error(), g.Logger)
	}

	err = g.Routes.prepareStaticRoutes(dir)
	if err != nil {
		return fmt.Errorf("could not prepare routes for static files: %s", err)
	}

	return
}

// Start starts listening for incoming requests on the specified port (e.g., ":8080").
func (g *Goster) Start(p string) {
	g.cleanUp()
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

func (g *Goster) StartTLS(addr string, certFile string, keyFile string) {
	g.cleanUp()
	LogInfo("LISTENING ON https://127.0.0.1"+addr, g.Logger)
	log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, g))
}

// ServeHTTP is the handler for incoming HTTP requests to the server.
// It parses the request, manages routing, and is required to implement the http.Handler interface.
func (g *Goster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	// Parse the URL and extract query parameters into the Meta struct
	urlPath := ctx.Request.URL.EscapedPath()
	method := ctx.Request.Method
	DefaultHeader(&ctx)

	// Construct a static route from URL path if it matches a specific dynamic or wildcard route
	for routePath, route := range g.Routes[method] {
		if route.Type == TypeStatic {
			continue
		}

		if urlMatchesRoute(urlPath, routePath) {
			ctx.Meta.ParsePath(urlPath, routePath)
			err := g.Routes.New(method, urlPath, route.Handler)
			if err != nil {
				_ = fmt.Errorf("route %s is duplicate", urlPath) // TODO: it is duplicate, handle
			}
			break
		}
	}

	// Validate the route based on the HTTP method and URL
	status := g.validateRoute(method, urlPath)
	if status != http.StatusOK {
		ctx.Response.WriteHeader(status)
		return
	}

	// Parses query params if any and adds them to query map
	ctx.Meta.ParseQueryParams(r.URL.String())

	// Execute global middleware handlers
	for _, middleware := range g.Middleware["*"] {
		err := middleware(&ctx)
		if err != nil {
			LogError(fmt.Sprintf("error occured while running global middleware: %s", err.Error()), g.Logger)
		}
	}
	logRequest(&ctx, g, nil) // TODO: streamline builtin middleware

	g.launchHandler(&ctx, method, urlPath)
}

// ------------------------------------------Private Methods--------------------------------------------------- //

// launchHandler launches the necessary handler for the incoming request based on the route.
func (g *Goster) launchHandler(ctx *Ctx, method, urlPath string) {
	cleanPath(&urlPath)
	route := g.Routes[method][urlPath]
	defer func() {
		err := route.Handler(ctx)
		// TODO: figure out what to do with handler error
		if err != nil {
			LogError(err.Error(), g.Logger)
		}
	}()
	// Run all route-specific middleware defined by the user
	for _, rh := range g.Middleware[urlPath] {
		err := rh(ctx)
		if err != nil {
			LogError(fmt.Sprintf("error occured while running middleware: %s", err.Error()), g.Logger)
		}
	}
}

// validateRoute checks if the route "reqURL" exists inside the `g.Routes` collection
// and under the method "reqMethod" so that the following expression evaluates to true:
//
//	if _, exists := g.Routes[m][u]; exists {
//		// some code
//	}
//
// If "reqURL" exists but not under the method "reqMethod", then the status `http.StatusMethodNotAllowed` is returned
//
// If "reqURL" doesn't exist then the status `http.StatusNotFound` is returned
func (g *Goster) validateRoute(method, urlPath string) int {
	cleanPath(&urlPath)
	for m := range g.Routes {
		_, exists := g.Routes[m][urlPath]
		if exists && m == method {
			return http.StatusOK
		} else if exists && m != method {
			return http.StatusMethodNotAllowed
		}
	}

	return http.StatusNotFound
}

func (g *Goster) cleanUp() {
	if engine.Config.BaseTemplateDir == "" {
		LogInfo("No specified template directory. Defaulting to `templates/`...", g.Logger)
		err := engine.SetTemplateDir("templates")
		if err != nil {
			LogError(err.Error(), g.Logger)
		}
	}
}
