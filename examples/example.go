package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()

	g.AddGlobalMiddleware(func(ctx *Goster.Ctx) error {
		fmt.Println("global middleware")
		return nil
	})

	g.Get("/", func(ctx *Goster.Ctx) error {
		q, exists := ctx.Get("q")
		msg := "Hello and welcome to the test server of Goster :D"
		if exists {
			if q == "69" {
				msg = "AHA! You found the secret message with the code 69! Your treasure is this 8====D~ 8="
			} else {
				msg = fmt.Sprintf("Almost there. %s isn't correct. You could try again but I wouldn't blame you if you gave up :c", q)
			}
		}
		ctx.ResponseWriter.Write([]byte(msg))
		return nil
	})

	g.Get("db/", func(ctx *Goster.Ctx) error {
		name, _ := ctx.Meta.Get("yourName")
		res := struct {
			Hey string `json:"hey"`
			You string `json:"you"`
		}{
			Hey: "Hello",
			You: name,
		}

		ctx.ResponseWriter.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, 200)
		ctx.ResponseWriter.JSON(res)

		return nil
	})

	g.Get("db/kati/:id", func(ctx *Goster.Ctx) error {
		db, exists := ctx.Meta.Get("db")

		if !exists {
			db = "{}"
		}
		ctx.ResponseWriter.Write([]byte(fmt.Sprintf("hello this is a multi route page at db/%s", db)))

		return nil
	})

	g.Get("path/:name", func(ctx *Goster.Ctx) error {
		name, exists := ctx.Params.Get("name")

		if !exists {
			msg := "please specify a corrent route"
			ctx.ResponseWriter.Write([]byte(msg))
			return errors.New(msg)
		}
		ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Hi, my name is %s", name)))
		return nil
	})

	g.Post("db/", func(ctx *Goster.Ctx) error {
		db := make([]byte, ctx.Request.ContentLength)
		ctx.Request.Body.Read(db)
		err := ioutil.WriteFile("./fake_db.txt", db, 0666)

		if err != nil {
			err_json := struct {
				Msg string `json:"msg"`
			}{
				Msg: err.Error(),
			}
			ctx.ResponseWriter.JSON(err_json)
			ctx.ResponseWriter.WriteHeader(500)
			return err
		}

		ctx.ResponseWriter.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, http.StatusCreated)

		ctx.ResponseWriter.JSON(db)

		return nil
	})

	g.Get("hey/", func(ctx *Goster.Ctx) error {
		heyPage, err := ioutil.ReadFile("./hey.html")

		if err != nil {
			fmt.Println(err)
		}
		ctx.ResponseWriter.Write(heyPage)
		return nil
	})

	g.Get("logs/", func(ctx *Goster.Ctx) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := ctx.ResponseWriter.JSON(log_map)

		if err != nil {
			Goster.LogError(err.Error(), g.Logger)
		}
		return nil
	})

	g.ListenAndServe(":8088")
}
