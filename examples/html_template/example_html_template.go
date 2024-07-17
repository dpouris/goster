package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()
	g.TemplateDir("templates")

	g.Get("/greet/:name", func(ctx *Goster.Ctx) error {
		name, exists := ctx.Path.Get("name")
		if exists {
			ctx.Template("index.gohtml", name)
		} else {
			ctx.Text("Name not found in path")
		}
		return nil
	})

	g.ListenAndServe(":8080")
}
