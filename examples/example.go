package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	Gottp "github.com/dpouris/gottp-server"
)

type JSONResponse struct {
	Hey string `json:"hey"`
	You string `json:"you"`
}

func main() {
	g := Gottp.Server()

	g.AddGlobalMiddleware(func(r Gottp.Res, req *Gottp.Req) error {
		fmt.Println("middleware")
		return nil
	})

	g.Get("/path", func(r Gottp.Res, req *Gottp.Req) error {

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

	g.Post("/path", func(r Gottp.Res, req *Gottp.Req) error {
		r.Write([]byte("This is a post :P"))
		res := make([]byte, 20*1024)
		req.Body.Read(res)
		fmt.Println(string(res))
		return nil
	})

	g.Get("/hey", func(r Gottp.Res, req *Gottp.Req) error {
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
			fmt.Println(err)
		}
		r.Write(heyPage)
		return nil
	})

	g.Get("/logs", func(r Gottp.Res, req *Gottp.Req) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := r.JSON(log_map)

		if err != nil {
			Gottp.LogError(err.Error(), g.Logger)
		}
		return nil
	})

	g.ListenAndServe(":8088")
}
