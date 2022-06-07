package gottp_server

import (
	"errors"
	"net/http"
)

func HandleMethod(g *Gottp, req *http.Request) (status int, err error) {
	u := req.URL.String()
	m := req.Method

	allowedMethods := make([]string, 0)

	for k, v := range g.Routes {
		for _, route := range v {
			if u == route.Route {
				allowedMethods = append(allowedMethods, k)
			}
		}
	}

	if len(allowedMethods) <= 0 {
		return 404, errors.New("404 NOT FOUND")
	}

	for _, v := range allowedMethods {
		if v == m {
			return 200, nil
		}
	}

	return 405, errors.New("405 METHOD NOT ALLOWED")
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
