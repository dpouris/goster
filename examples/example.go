package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	Gottp "github.com/dpouris/gottp-server"
)

type JSONResponse struct {
	Hey string `json:"hey"`
	You string `json:"you"`
}

func main() {
	g := Gottp.Server()

	g.AddGlobalMiddleware(func(r http.ResponseWriter, req *http.Request) error {
		fmt.Println("middleware")
		return nil
	})

	g.Get("/path", func(r http.ResponseWriter, req *http.Request) error {

		res := JSONResponse{
			Hey: "Hello",
			You: "World",
		}
		marsalled, err := json.Marshal(res)

		if err != nil {
			Gottp.LogError(err.Error(), g.Logger)
		}

		r.Header().Set("Content-Type", "application/json")
		r.WriteHeader(http.StatusOK)
		r.Write(marsalled)

		return nil
	})

	g.Post("/path", func(r http.ResponseWriter, req *http.Request) error {
		r.Write([]byte("This is a post :P"))
		res := make([]byte, 20*1024)
		req.Body.Read(res)
		fmt.Println(string(res))
		return nil
	})

	g.Get("/hey", func(r http.ResponseWriter, req *http.Request) error {
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
			fmt.Println(err)
		}
		r.Write(heyPage)
		return nil
	})

	g.Get("/logs", func(r http.ResponseWriter, req *http.Request) error {
		_, err := r.Write([]byte(strings.Join(g.Logs, "\n")))

		if err != nil {
			Gottp.LogError(err.Error(), g.Logger)
		}

		r.Header().Set("Content-Type", "text/plain")
		return nil
	})

	g.ListenAndServe(":8088")
}
