package main

import (
	"log"

	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()

	// Middleware to log requests
	g.UseGlobal(func(ctx *Goster.Ctx) error {
		log.Printf("Received request for %s", ctx.Request.URL.Path)
		return nil
	})

	g.Get("/", func(ctx *Goster.Ctx) error {
		ctx.Text("Middleware example")
		return nil
	})

	g.Start(":8080")
}
