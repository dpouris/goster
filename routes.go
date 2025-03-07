package goster

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Routes map[string]map[string]Route

func (rs *Routes) prepareStaticRoutes(dir string) (err error) {
	staticPaths := engine.Config.StaticFilePaths

	for relPath := range staticPaths {
		file, err := os.Open(staticPaths[relPath])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open static file `%s` in dir `%s`\n", file.Name(), relPath)
			return err
		}
		defer file.Close()

		// read the file contents
		bytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read static file `%s` in dir `%s`\n", file.Name(), relPath)
			return err
		}

		// compute the route path relative to the static directory
		routePath := filepath.Join(dir, relPath)
		cleanPath(&routePath)

		// register a GET route that serves the file content
		rs.New("GET", routePath, func(ctx *Ctx) error {
			contentType := getContentType(file.Name())
			ctx.Response.Header().Set("Content-Type", contentType)
			ctx.Response.Write(bytes)
			return nil
		})

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

	routeType := "normal"
	if strings.Contains(url, ":") {
		routeType = "dynamic"
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
