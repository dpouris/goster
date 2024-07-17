package goster

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

type Ctx struct {
	Request  *http.Request
	Response Response
	Meta
}

// Send an HTML template t file to the client. If template not in template dir then will return error.
func (c *Ctx) Template(t string, data any) (err error) {
	files := engine.Config.FilePaths

	for tmpl := range files {
		if tmpl == path.Join(engine.Config.BaseStaticDir, t) {
			tm := template.Must(template.ParseFiles(tmpl))
			err = tm.Execute(c.Response, data)

			if err != nil {
				return err
			}

		}
	}

	return
}

// Send an HTML f file to the client. If if file not in FilesDir dir then will return error.
func (c *Ctx) HTML(f string) (err error) {
	files := engine.Config.FilePaths

	for fp := range files {
		if path.Join(engine.Config.BaseStaticDir, f) == fp {
			file, err := os.Open(fp)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}

			defer file.Close()

			f_size, _ := file.Stat()
			buf := make([]byte, f_size.Size())
			io.ReadFull(file, buf)

			t := string(buf)
			c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")

			fmt.Fprintln(c.Response.ResponseWriter, t)
		}
	}
	return
}

// Send plain text to the client
func (c *Ctx) Text(s string) {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.Header().Set("Content-Length", fmt.Sprint(len(s)))
	fmt.Fprint(c.Response.ResponseWriter, s)
}

// Send back a JSON response. Supply j with a value that's valid marsallable(?) to JSON -> error
func (c *Ctx) JSON(j any) error {
	if v, ok := j.([]byte); ok {
		cleanEmptyBytes(&v)
		c.Response.Header().Set("Content-Type", "application/json")
		c.Response.Write(v)
		return nil
	}

	v, err := json.Marshal(j)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	c.Response.Header().Set("Content-Type", "application/json")
	_, err = c.Response.Write(v)

	return err
}
