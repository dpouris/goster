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

	g.Use("/", func(ctx *Goster.Ctx) error {
		fmt.Println("\"/\" middleware")
		return nil
	})

	g.UseGlobal(func(ctx *Goster.Ctx) error {
		fmt.Println("global!!")
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
		ctx.Response.Write([]byte(msg))
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

		ctx.Response.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, 200)
		ctx.JSON(res)

		return nil
	})

	g.Get("greet/:name", func(ctx *Goster.Ctx) (err error) {
		name, exists := ctx.Get("name")
		if !exists {
			ctx.Text("Please navigate to /greet/<yourName>")
			return
		}

		err = g.TemplateDir("templates")

		if err != nil {
			ctx.Text(err.Error())
			return nil
		}

		err = ctx.Template("index.gohtml", name)

		if err != nil {
			fmt.Println(err)
		}

		return nil
	})

	g.Get("age/", func(ctx *Goster.Ctx) (err error) {
		age, exists := ctx.Get("age")

		if !exists {
			msg := "please specify an ?age param in the url"
			ctx.Text(msg)
			err = errors.New(msg)
			return
		}

		ctx.Text(fmt.Sprintf("Hi, I'm %s years old!", age))
		return
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
			ctx.JSON(err_json)
			ctx.Response.WriteHeader(500)
			return err
		}

		ctx.Response.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}, http.StatusCreated)

		ctx.JSON(db)

		return nil
	})

	g.Get("hey/", func(ctx *Goster.Ctx) error {
		err := g.TemplateDir(".")

		if err != nil {
			panic(err)
		}

		if err := ctx.HTML("hey.html"); err != nil {
			panic(err)
		}
		return nil
	})

	g.Get("logs/", func(ctx *Goster.Ctx) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := ctx.JSON(log_map)

		if err != nil {
			Goster.LogError(err.Error(), g.Logger)
		}
		return nil
	})

	g.ListenAndServe(":8088")
}
