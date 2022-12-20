package gottp_client

import (
	"errors"
	"net/http"
)

func HandleMethod(g *Gottp, req *http.Request) (status int, err error) {
	u := req.URL.String()
	m := req.Method

	allowedMethods := make([]string, 0)

	for name := range g.Routes[m] {
		if u == name {
			allowedMethods = append(allowedMethods, m)
		}
	}

	if len(allowedMethods) <= 0 {
		return http.StatusNotFound, errors.New("404 NOT FOUND")
	}

	for _, v := range allowedMethods {
		if v == m {
			return http.StatusOK, nil
		}
	}

	return http.StatusMethodNotAllowed, errors.New("405 METHOD NOT ALLOWED")
}

func HandleLog(route string, method string, err error, g *Gottp) {
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
