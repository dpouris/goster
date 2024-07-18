package main

import (
	"fmt"
	"strconv"

	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.NewServer()

	g.Get("/greet", func(ctx *Goster.Ctx) error {
		ctx.Text("Hey there stranger!\nGo to `/greet/:yourName` to see a message just for you!")
		return nil
	})

	g.Get("/greet/:name", func(ctx *Goster.Ctx) error {
		name, _ := ctx.Path.Get("name")
		ctx.Text(fmt.Sprintf("Hello there %s!\nGo to `/greet/:yourName/:yourAge` to see another message just for you!", name))
		return nil
	})

	g.Get("/greet/:name/:age", func(ctx *Goster.Ctx) error {
		name, _ := ctx.Path.Get("name")
		ageStr, _ := ctx.Path.Get("age")

		age, err := strconv.Atoi(ageStr)

		if err != nil {
			ctx.Text(fmt.Sprintf("%s the value `%s` is not a valid age value :C", name, ageStr))
		}

		ctx.Text(fmt.Sprintf("Hello %s who is %d years old!", name, age))
		return nil
	})

	g.ListenAndServe(":8080")
}
