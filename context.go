package goster

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type Ctx struct {
	Request  *http.Request
	Response Response
	Meta
}

func NewContext(r *http.Request, w http.ResponseWriter) Ctx {
	return Ctx{
		Request:  r,
		Response: Response{w},
		Meta: Meta{
			Query: make(map[string]string),
			Path:  make(map[string]string),
		},
	}
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

func (c *Ctx) Redirect(url string, code int) {
	if u, err := urlpkg.Parse(url); err == nil {
		// If url was relative, make its path absolute by
		// combining with request path.
		// The client would probably do this for us,
		// but doing it ourselves is more reliable.
		// See RFC 7231, section 7.1.2
		if u.Scheme == "" && u.Host == "" {
			oldpath := c.Request.URL.Path
			if oldpath == "" { // should not happen, but avoid a crash if it does
				oldpath = "/"
			}

			// no leading http://server
			if url == "" || url[0] != '/' {
				// make relative path absolute
				olddir, _ := path.Split(oldpath)
				url = olddir + url
			}

			var query string
			if i := strings.Index(url, "?"); i != -1 {
				url, query = url[:i], url[i:]
			}

			// clean up but preserve trailing slash
			trailing := strings.HasSuffix(url, "/")
			url = path.Clean(url)
			if trailing && !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += query
		}
	}

	h := c.Response.Header()

	// RFC 7231 notes that a short HTML body is usually included in
	// the response because older user agents may not understand 301/307.
	// Do it only if the request didn't already have a Content-Type headec.Request.
	_, hadCT := h["Content-Type"]

	h.Set("Location", hexEscapeNonASCII(url))
	if !hadCT && (c.Request.Method == "GET" || c.Request.Method == "HEAD") {
		h.Set("Content-Type", "text/html; charset=utf-8")
	}
	c.Response.WriteHeader(code)

	// Shouldn't send the body for POST or HEAD; that leaves GET.
	if !hadCT && c.Request.Method == "GET" {
		body := "<a href=\"" + htmlEscape(url) + "\">" + http.StatusText(code) + "</a>.\n"
		fmt.Fprintln(c.Response, body)
	}
}
