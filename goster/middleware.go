package goster

import (
	"errors"
	"net/http"
)

func HandleMethod(g *Goster, url, method string) (status int, err error) {
	allowedMethods := make([]string, 0)

	for name := range g.Routes[method] {
		if url == name {
			allowedMethods = append(allowedMethods, method)
		}
	}

	if len(allowedMethods) <= 0 {
		return http.StatusNotFound, errors.New("404 not found")
	}

	for _, v := range allowedMethods {
		if v == method {
			return http.StatusOK, nil
		}
	}

	return http.StatusMethodNotAllowed, errors.New("405 method not allowed")
}

func HandleLog(route string, method string, err error, g *Goster) {
	if err != nil {
		l := err.Error()
		g.Logs = append(g.Logs, l)
		LogError(l, g.Logger)
		return
	}
	l := "[" + method + "]" + " ON ROUTE " + route
	g.Logs = append(g.Logs, l)
	LogInfo(l, g.Logger)
}

// Adds basic headers
func DefaultHeader(h *http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Cross-Origin-Opener-Policy", "same-origin")

	h.Set("Connection", "Keep-Alive")
	h.Set("Keep-Alive", "timeout=5, max=997")
}
