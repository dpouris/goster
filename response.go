package goster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	http.ResponseWriter
}

// Supply h with a map[string]string for the headers and s with an int representing the response status code or use the http.Status(...). They keys and values will be translated to the header of the response and the header will be locked afterwards not allowing changes to be made.
func (g *Response) NewHeaders(h map[string]string, s int) {
	for k, v := range h {
		g.Header().Set(k, v)
	}

	g.WriteHeader(s)
}

// Send back a JSON response. Supply j with a value that's valid marsallable(?) to JSON -> error
func (r Response) JSON(j any) error {
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

func (r Response) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r Response) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r Response) Header() http.Header {
	return r.ResponseWriter.Header()
}

func cleanEmptyBytes(b *[]byte) {
	cleaned := []byte{}

	for _, v := range *b {
		if v == 0 {
			break
		}
		cleaned = append(cleaned, v)
	}
	*b = cleaned
}
