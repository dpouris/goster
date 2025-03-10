package main

import (
	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()

	g.Get("/json", func(ctx *Goster.Ctx) error {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "Hello, JSON!",
		}
		ctx.JSON(response)
		return nil
	})

	g.Start(":8080")
}
