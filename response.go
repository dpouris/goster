package goster

import (
	"net/http"
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
