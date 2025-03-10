package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()
	err := g.TemplateDir("/templates")
	if err != nil {
		Goster.LogError("could not set templates dir", g.Logger)
	}

	g.Get("/greet/:name", func(ctx *Goster.Ctx) error {
		name, exists := ctx.Path.Get("name")
		if exists {
			ctx.Template("index.gohtml", name)
		} else {
			ctx.Text("Name not found in path")
		}
		return nil
	})

	g.Start(":8080")
}
