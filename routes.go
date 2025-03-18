package goster

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Route represents an HTTP route with a type and a handler function.
type Route struct {
	Type    string         // Type specifies the type of the route (e.g., TypeStatic, TypeDynamic, TypeWildcard).
	Handler RequestHandler // Handler is the function that handles the route.
}

type Routes map[string]map[string]Route

func (rs *Routes) prepareStaticRoutes(dir string) (err error) {
	staticPaths := engine.Config.StaticFilePaths
	for relPath := range staticPaths {
		staticPath := staticPaths[relPath]
		file, err := os.Open(staticPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open static file `%s`\n", file.Name())
			return err
		}

		// register a GET route that serves the static file
		routePath := filepath.Join(dir, relPath)
		cleanPath(&routePath)
		err = rs.New("GET", routePath, func(ctx *Ctx) error {
			return staticFileHandler(ctx, file)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't add route `%s`. Most likely there's a duplicate entry\n", routePath)
		}
	}

	return
}

// New creates a new Route for the specified method and url using the provided handler. If the Route already exists an error is returned.
func (rs *Routes) New(method string, url string, handler RequestHandler) (err error) {
	for name := range (*rs)[method] {
		if name == url {
			err = fmt.Errorf("[%s] -> [%s] route already exists", method, url)
			return
		}
	}

	routeType := TypeStatic
	if strings.Contains(url, ":") {
		routeType = TypeDynamic
	} else if strings.Contains(url, "*") {
		routeType = TypeWildcard
	}

	cleanPath(&url)
	(*rs)[method][url] = Route{Type: routeType, Handler: handler}

	return
}

// Get creates a new Route under the GET method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Get(url string, handler RequestHandler) error {
	return g.Routes.New("GET", url, handler)
}

// Post creates a new Route under the POST method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Post(path string, handler RequestHandler) error {
	return g.Routes.New("POST", path, handler)
}

// Patch creates a new Route under the PATCH method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Patch(path string, handler RequestHandler) error {
	return g.Routes.New("PATCH", path, handler)
}

// Put creates a new Route under the PUT method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Put(path string, handler RequestHandler) error {
	return g.Routes.New("PUT", path, handler)
}

// Delete creates a new Route under the DELETE method for `path`. If the Route aleady exists an error is returned.
func (g *Goster) Delete(path string, handler RequestHandler) error {
	return g.Routes.New("DELETE", path, handler)
}

func staticFileHandler(ctx *Ctx, file *os.File) (err error) {
	// read the file contents
	_, err = file.Seek(0, 0)
	if err != nil {
		return
	}
	fInfo, _ := file.Stat()
	fSize := fInfo.Size()
	buffer := make([]byte, fSize)
	_, _ = io.ReadFull(file, buffer)

	// prepare and write response
	contentType := getContentType(file.Name())
	ctx.Response.Header().Set("Content-Type", contentType)
	_, err = ctx.Response.Write(buffer)
	return
}
