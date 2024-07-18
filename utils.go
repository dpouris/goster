package goster

import (
	"strings"
)

// Adds basic headers
func DefaultHeader(c *Ctx) {
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

func cleanPath(path *string) {
	if len(*path) == 0 {
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	*path = strings.TrimSuffix(*path, "/")
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
