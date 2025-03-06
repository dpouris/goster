package goster

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Routes map[string]map[string]Route

// AddStaticDir registers GET routes for all files found under the given directory.
// It uses the current working directory as the base path and creates a route for each file,
// serving its content along with the appropriate Content-Type header.
// The `dir` parameter should be a relative path from the working directory (the directory you'll execute the program).
func (rs *Routes) AddStaticDir(dir string) error {
	exPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine working directory for static dir %s\n", dir)
		return err
	}

	// construct the full path to the static directory
	staticPath := path.Join(exPath, dir)

	// walk the directory and register a GET route for each file found
	err = filepath.WalkDir(staticPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// process only files (skip directories)
		if !d.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot open file %s in static dir %s\n", filePath, dir)
				return err
			}
			defer file.Close()

			// read the file contents
			bytes, err := io.ReadAll(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read file %s in static dir %s\n", filePath, dir)
				return err
			}
			contents := string(bytes)

			// compute the route path relative to the static directory
			relPath, _ := filepath.Rel(staticPath, filePath)
			routePath := filepath.Join(dir, relPath)
			cleanPath(&routePath)

			// register a GET route that serves the file content
			rs.New("GET", routePath, func(ctx *Ctx) error {
				contentType := getContentType(file.Name())
				ctx.Response.Header().Set("Content-Type", contentType)
				ctx.Text(contents)
				return nil
			})
		}

		return nil
	})

	return err
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
