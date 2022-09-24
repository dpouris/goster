package gottp_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Res struct {
	http.ResponseWriter
}

// Supply h with a map[string]string for the headers and s with an int representing the response status code or use the http.Status(...). They keys and values will be translated to the header of the response and the header will be locked afterwards not allowing changes to be made.
func (g *Res) NewHeaders(h map[string]string, s int) {
	for k, v := range h {
		g.Header().Set(k, v)
	}

	g.WriteHeader(s)
}

// Send back a JSON response. Supply j with a value that's valid marsallable(?) to JSON -> error
func (r Res) JSON(j any) error {
	if v, ok := j.([]byte); ok {
		cleanEmptyBytes(&v)
		r.Header().Set("Content-Type", "application/json")
		r.Write(v)
		return nil
	}

	v, err := json.Marshal(j)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}
	r.Header().Set("Content-Type", "application/json")
	_, err = r.Write(v)

	return err
}

func (r Res) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r Res) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r Res) Header() http.Header {
	return r.ResponseWriter.Header()
}

func cleanEmptyBytes(b *[]byte) {
	new_b := []byte{}

	for _, v := range *b {
		if v == 0 {
			break
		}
		new_b = append(new_b, v)
	}
	*b = new_b
}
