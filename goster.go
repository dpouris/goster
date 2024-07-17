package goster

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Goster is the main structure of the package. It handles the addition of new routes and middleware, and manages logging.
type Goster struct {
	Routes     map[string]map[string]Route // Routes is a map of HTTP methods to their respective route handlers.
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
// which is joined to Engine.Config.BaseStaticDir (default is the current working directory).
//
// This instructs the engine where to look for static files like .html, .gohtml, .css, .js, etc.
// If the directory doesn't exist, it will return an appropriate error.
//
// TODO: Should revise this function
func (g *Goster) TemplateDir(d string) (err error) {
	err = engine.SetTemplateDir(d)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	return
}

// ListenAndServe starts listening for incoming requests on the specified port (e.g., ":8080").
func (g *Goster) ListenAndServe(p string) {
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
			Query: Params{
				values: make(map[string]string),
			},
			Path: Path{
				values: make(map[string]string),
			},
		},
	}
	// Parse the URL and extract query parameters into the Meta struct
	url := ctx.Request.URL.String()
	ctx.Meta.ParseUrl(&url) // TODO: handle this error

	DefaultHeader(&ctx)

	// Construct normal routes from URL path if they match a specific dynamic route
	g.constructPathRoutes(&ctx, url)

	// Validate the route based on the HTTP method and URL path
	status := g.validateRoute(ctx.Request.Method, url)

	if status != http.StatusOK {
		ctx.Response.WriteHeader(status)
		return
	}
	g.launchHandler(&ctx, url)

	// Execute global middleware handlers
	for _, rh := range g.Middleware["*"] {
		rh(&ctx)
	}
}

// ------------------------------------------Private Methods--------------------------------------------------- //

// launchHandler launches the necessary handler for the incoming request based on the route.
func (g *Goster) launchHandler(ctx *Ctx, url string) {
	method := ctx.Request.Method

	HandleLog(ctx, g, nil)

	route := g.Routes[method][url]

	defer route.Handler(ctx)

	// Run all route-specific middleware defined by the user
	for _, rh := range g.Middleware[url] {
		rh(ctx)
	}
}

// constructPathRoutes constructs normal routes from URL path if they match a specific dynamic route.
func (g *Goster) constructPathRoutes(ctx *Ctx, url string) {
	method := ctx.Request.Method

	for path, route := range g.Routes[method] {
		if route.Type != "dynamic" {
			continue
		}

		if isDynamicRouteMatch(ctx, url, path) {
			newRoute := Route{
				Type:    "normal",
				Handler: route.Handler,
			}
			engine.Goster.Routes[method][url] = newRoute

			break
		}
	}
}

// validateRoute checks if the route "u" exists inside the `map[string]map[string]Route` collection
// and under the method "m" so that the following expression evaluates to true:
//
//	if _, exists := g.Routes[m][u]; exists {
//		// some code
//	}
//
// If "u" exists but not under the method "m", then the status `http.StatusMethodNotAllowed` is returned
//
// If "u" doesn't exist then the status `http.StatusNotFound` is returned
func (g *Goster) validateRoute(m, u string) int {
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
		return http.StatusNotFound
	}

	if !methodAllowed {
		return http.StatusMethodNotAllowed
	}
	return http.StatusOK
}

// isDynamicRouteMatch checks if the raw (stripped from Query Parameters) URL path `url` matches a Dynamic Route path
// `dynPath`. A Dynamic Route is a path string that has the following format: "path/anotherPath/:variablePathname" where `:variablePathname`
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
func isDynamicRouteMatch(ctx *Ctx, url string, dynPath string) (match bool) {
	match = false
	dynPathPattern := regexp.MustCompile(`\:\w+`)
	// For example in the dynamic url path "greet/:name":
	// 	- variablePathStart would be 6 (the index of the char ':')
	// 	- variablePath would be "name"
	variablePathStart := dynPathPattern.FindStringIndex(dynPath)[0]
	variablePath := strings.Trim(dynPathPattern.FindString(dynPath), ":")

	if len(url) < variablePathStart {
		return
	}

	variablePathValue := url[variablePathStart:]
	stopIdx := strings.IndexFunc(variablePathValue, func(r rune) bool { return r == '/' })
	if stopIdx > 0 {
		variablePathValue = variablePathValue[:stopIdx]
	}

	constructedPath := dynPathPattern.ReplaceAllString(dynPath, variablePathValue)

	if constructedPath == url {
		match = true
		ctx.Meta.Path.values[variablePath] = variablePathValue
		return
	}

	return
}
