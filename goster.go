package goster

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Goster is the main structure of the package. It handles the addition of new routes and middleware, and manages logging.
type Goster struct {
	Routes     Routes                      // Routes is a map of HTTP methods to their respective route handlers.
	Middleware map[string][]RequestHandler // Middleware is a map of routes to their respective middleware handlers.
	Logger     *log.Logger                 // Logger is used for logging information and errors.
	Logs       []string                    // Logs stores logs for future reference.
}

// Route represents an HTTP route with a type and a handler function.
type Route struct {
	Type    string         // Type specifies the type of the route (e.g., "static", "dynamic").
	Handler RequestHandler // Handler is the function that handles the route.
}

// RequestHandler is a type for functions that handle HTTP requests within a given context.
type RequestHandler func(ctx *Ctx) error

// ------------------------------------------Public Methods--------------------------------------------------- //

// NewServer creates a new Goster instance.
func NewServer() *Goster {
	g := engine.init()
	return g
}

// UseGlobal adds middleware handlers that will be applied to every single request.
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
		fmt.Fprint(os.Stderr, err)
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
		fmt.Fprint(os.Stderr, err)
	}

	err = g.Routes.prepareStaticRoutes(dir)
	if err != nil {
		return fmt.Errorf("could not prepare routes for static files: %s", err)
	}

	return
}

// ListenAndServe starts listening for incoming requests on the specified port (e.g., ":8080").
func (g *Goster) ListenAndServe(p string) {
	g.cleanUp()
	LogInfo("LISTENING ON http://127.0.0.1"+p, g.Logger)
	log.Fatal(http.ListenAndServe(p, g))
}

// ServeHTTP is the handler for incoming HTTP requests to the server.
// It parses the request, manages routing, and is required to implement the http.Handler interface.
func (g *Goster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx{
		Request:  r,
		Response: Response{w},
		Meta: Meta{
			Query: make(map[string]string),
			Path:  make(map[string]string),
		},
	}
	// Parse the URL and extract query parameters into the Meta struct
	reqURL := ctx.Request.URL.String()
	reqMethod := ctx.Request.Method
	DefaultHeader(&ctx)

	// Construct a normal route from URL path if it matches a specific dynamic route
	var dynPath string
	for dynPathURL, route := range g.Routes[reqMethod] {
		if route.Type != "dynamic" {
			continue
		}
		if g.resolveDynRoute(reqMethod, reqURL, dynPathURL, route) {
			dynPath = dynPathURL
			break
		}
	}

	// Validate the route based on the HTTP method and URL path
	status := g.validateRoute(reqMethod, reqURL)
	if status != http.StatusOK {
		ctx.Response.WriteHeader(status)
		return
	}

	ctx.prepare(reqURL, dynPath)
	g.launchHandler(&ctx, reqMethod, reqURL)
	// Execute global middleware handlers
	for _, rh := range g.Middleware["*"] {
		rh(&ctx)
	}
}

// ------------------------------------------Private Methods--------------------------------------------------- //

// launchHandler launches the necessary handler for the incoming request based on the route.
func (g *Goster) launchHandler(ctx *Ctx, reqMethod, reqURL string) {
	cleanPath(&reqURL)
	HandleLog(ctx, g, nil) // TODO: ???????

	route := g.Routes[reqMethod][reqURL]
	defer route.Handler(ctx)
	// Run all route-specific middleware defined by the user
	for _, rh := range g.Middleware[reqURL] {
		rh(ctx)
	}
}

// resolveDynRoute constructs a normal route from URL path if it matches a specific dynamic route.
func (g *Goster) resolveDynRoute(reqMethod, reqURL, dynPathURL string, route Route) bool {
	cleanPath(&reqURL)

	if g.isDynRouteMatch(reqURL, dynPathURL) {
		g.Routes.New(reqMethod, reqURL, route.Handler)
		return true
	}

	return false
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
func (g *Goster) validateRoute(reqMethod, reqURL string) int {
	cleanPath(&reqURL)
	methodAllowed := false
	routeExists := false

	for method := range g.Routes {
		if _, exists := g.Routes[method][reqURL]; exists && method == reqMethod {
			routeExists = true
			methodAllowed = true
			break
		} else if exists && method != reqMethod {
			routeExists = true
		}
	}

	if !routeExists {
		return http.StatusNotFound
	}

	if !methodAllowed {
		return http.StatusMethodNotAllowed
	}
	return http.StatusOK
}

// isDynRouteMatch checks URL path `reqURL` matches a Dynamic Route path
// `dynPathURL`. A Dynamic Route is a path string that has the following format: "path/anotherPath/:variablePathname" where `:variablePathname`
// is a catch-all identifier that matches any route with the same structure up to that point.
//
// Ex:
//
//	var ctx = ...
//	var url = "path/anotherPath/andYetAnotherPath"
//	var dynPath = "path/anotherPath/:identifier"
//	if !isDynamicRouteMatch(&ctx, url, dynPath) {
//			panic(...)
//	}
//
// The above code will not panic as the isDynamicRouteMatch will evaluate to `true`
func (g *Goster) isDynRouteMatch(reqURL string, dynPathURL string) (match bool) {
	cleanPath(&reqURL)
	cleanPath(&dynPathURL)

	match = true
	_, err := matchDynPathValue(dynPathURL, reqURL)

	if err != nil {
		match = false
		return
	}

	return
}

func (g *Goster) cleanUp() {
	if engine.Config.BaseTemplateDir == "" {
		LogInfo("No specified template directory. Defaulting to `templates/`...", g.Logger)
		engine.SetTemplateDir("templates")
	}
}
