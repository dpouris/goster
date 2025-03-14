package main

import (
	"net/http"
	"strconv"

	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()
	err := g.TemplateDir("/templates")
	if err != nil {
		Goster.LogError("could not set templates dir", g.Logger)
	}

	g.Get("/greet/:name", func(ctx *Goster.Ctx) error {
		name, _ := ctx.Path.Get("name")
		age, exists := ctx.Query.Get("age")
		if !exists {
			ctx.Response.WriteHeader(http.StatusUnprocessableEntity)
			ctx.JSON(map[string]string{
				"status":  strconv.Itoa(http.StatusUnprocessableEntity),
				"context": "`age` query parameter not specified",
			})
			return nil
		}

		ctx.Template("index.gohtml", map[string]string{
			"name": name,
			"age":  age,
		})
		return nil
	})

	g.Start(":8080")
}
