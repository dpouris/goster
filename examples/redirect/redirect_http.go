package main

import (
	"net/http"
	"net/url"

	Goster "github.com/dpouris/goster"
)

const HOST = "https://test.com"

func main() {
	g := Goster.NewServer()

	g.Get("*path", func(ctx *Goster.Ctx) error {
		pathname, _ := ctx.Path.Get("path")
		status := http.StatusMovedPermanently
		if len(pathname) > 1 {
			status = http.StatusTemporaryRedirect
		}
		location, _ := url.JoinPath(HOST, pathname)
		ctx.Redirect(location, status)
		return nil
	})

	g.UseGlobal(func(ctx *Goster.Ctx) error {
		urlPath, _ := ctx.Path.Get("path")
		println("Path: " + urlPath)
		return nil
	})

	g.Start(":8001")
}
