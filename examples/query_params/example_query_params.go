package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()

	g.Get("/", func(ctx *Goster.Ctx) error {
		q, exists := ctx.Query.Get("q")
		if exists {
			ctx.Text("Query parameter 'q' is: " + q)
		} else {
			ctx.Text("Query parameter 'q' not found")
		}
		return nil
	})

	g.ListenAndServe(":8080")
}
