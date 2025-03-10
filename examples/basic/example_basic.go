package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()
	err := g.StaticDir("/static")
	if err != nil {
		Goster.LogError("could not set static dir", g.Logger)
	}

	err = g.TemplateDir("/templates")
	if err != nil {
		Goster.LogError("could not set templates dir", g.Logger)
	}

	g.Get("/", func(ctx *Goster.Ctx) error {
		ctx.Text("Welcome to Goster!")
		return nil
	})

	g.Start(":8080")
}
