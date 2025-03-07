package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()
	g.StaticDir("/static")
	g.TemplateDir("/templates")

	g.Get("/", func(ctx *Goster.Ctx) error {
		ctx.Text("Welcome to Goster!")
		return nil
	})

	g.ListenAndServe(":8080")
}
