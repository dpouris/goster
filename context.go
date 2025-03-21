package goster

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

type Ctx struct {
	Request  *http.Request
	Response Response
	Meta
}

// Send an HTML template t file to the client. If template not in template dir then will return error.
func (c *Ctx) Template(t string, data any) (err error) {
	templatePaths := engine.Config.TemplatePaths

	// iterate through all known templates
	for tmplId := range templatePaths {
		// if given template matches a known template get the template path, parse it and write it to response
		if tmplId == t {
			tmpl := template.Must(template.ParseFiles(templatePaths[tmplId]))
			err = tmpl.Execute(c.Response, data)

			if err != nil {
				return err
			}

		}
	}

	return
}

// Send an HTML template t file to the client. TemplateWithFuncs supports functions to be embedded in the html template for use. If template not in template dir then will return error.
func (c *Ctx) TemplateWithFuncs(t string, data any, funcMap template.FuncMap) (err error) {
	templatePaths := engine.Config.TemplatePaths

	// iterate through all known templates
	for tmplId := range templatePaths {
		// if given template matches a known template get the template path, parse it and write it to response
		if tmplId == t {
			tmplFile := templatePaths[tmplId]
			baseFilename := filepath.Base(tmplId)
			tmpl := template.Must(template.New(baseFilename).Funcs(funcMap).ParseFiles(tmplFile))
			err = tmpl.Execute(c.Response, data)

			if err != nil {
				return err
			}

		}
	}

	return
}

// Send an HTML f file to the client. If if file not in FilesDir dir then will return error.
func (c *Ctx) HTML(t string) (err error) {
	templatePaths := engine.Config.TemplatePaths

	for tmplId := range templatePaths {
		if tmplId == t {
			// open template
			file, err := os.Open(templatePaths[tmplId])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			defer file.Close()

			// read file
			fInfo, _ := file.Stat()
			buf := make([]byte, fInfo.Size())
			_, err = io.ReadFull(file, buf)
			if err != nil {
				return err
			}

			t := string(buf)

			// set headers
			contentType := getContentType(file.Name())
			c.Response.Header().Set("Content-Type", contentType)

			fmt.Fprint(c.Response.ResponseWriter, t) // write response
		}
	}
	return
}

// Send plain text to the client
func (c *Ctx) Text(s string) {
	c.Response.Header().Set("Content-Length", fmt.Sprint(len(s)))
	c.Response.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(c.Response.ResponseWriter, s)
}

// Send back a JSON response. Supply j with a value that's valid marsallable(?) to JSON -> error
func (c *Ctx) JSON(j any) (err error) {
	if v, ok := j.([]byte); ok {
		v = slices.DeleteFunc(v, func(b byte) bool {
			return b == 0
		})
		c.Response.Header().Set("Content-Type", "application/json")
		_, err = c.Response.Write(v)
		return
	}

	v, err := json.Marshal(j)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	c.Response.Header().Set("Content-Type", "application/json")
	_, err = c.Response.Write(v)
	return
}
