package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	Gottp "github.com/dpouris/gottp-client"
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

		r.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, 200)
		r.JSON(res)

		return nil
	})

	g.Post("/path", func(r Gottp.Res, req *Gottp.Req) error {
		db := make([]byte, req.ContentLength)
		req.Body.Read(db)
		err := ioutil.WriteFile("./examples/fake_db.txt", db, 0666)

		if err != nil {
			return err
		}

		r.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, http.StatusCreated)

		r.JSON(db)

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
