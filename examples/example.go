package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	Gottp "github.com/dpouris/gottp-client"
)

func main() {
	g := Gottp.Server()

	g.AddGlobalMiddleware(func(ctx *Gottp.Ctx) error {
		fmt.Println("middleware")
		return nil
	})

	g.Get("path/", func(ctx *Gottp.Ctx) error {

		res := struct {
			Hey string `json:"hey"`
			You string `json:"you"`
		}{
			Hey: "Hello",
			You: "World",
		}

		ctx.ResponseWriter.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, 200)
		ctx.ResponseWriter.JSON(res)

		return nil
	})

	g.Get("path/:name", func(ctx *Gottp.Ctx) error {
		name := ctx.Params.Get("name")
		ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Hi, my name is %s", name)))
		return nil
	})

	g.Post("path/", func(ctx *Gottp.Ctx) error {
		db := make([]byte, ctx.Request.ContentLength)
		ctx.Request.Body.Read(db)
		err := ioutil.WriteFile("./examples/fake_db.txt", db, 0666)

		if err != nil {
			return err
		}

		ctx.ResponseWriter.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, http.StatusCreated)

		ctx.ResponseWriter.JSON(db)

		return nil
	})

	g.Get("hey/", func(ctx *Gottp.Ctx) error {
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
			fmt.Println(err)
		}
		ctx.ResponseWriter.Write(heyPage)
		return nil
	})

	g.Get("logs/", func(ctx *Gottp.Ctx) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := ctx.ResponseWriter.JSON(log_map)

		if err != nil {
			Gottp.LogError(err.Error(), g.Logger)
		}
		return nil
	})

	g.ListenAndServe(":8088")
}
