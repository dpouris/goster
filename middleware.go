package gottp_server

import (
	"errors"
	"net/http"
)

func HandleMethod(g *Gottp, req *http.Request) (status int, err error) {
	u := req.URL.String()
	m := req.Method

	allowedMethods := make([]string, 0)

	for k, v := range g.routes {
		for _, route := range v {
			if u == route.route {
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

func HandleLog(route string, method string) {
	LogInfo("[" + method + "]" + " ON ROUTE " + route)
}
