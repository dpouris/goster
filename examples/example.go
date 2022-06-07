package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	Gottp "github.com/dpouris/gottp-server"
)

func main() {
	g := Gottp.Server()

	g.AddGlobalMiddleware(func(r http.ResponseWriter, req *http.Request) error {
		fmt.Println("middleware")
		return nil
	})

	g.Get("/path", func(r http.ResponseWriter, req *http.Request) error {
		res := make(map[string]any, 10)
		res["hey"] = "you"
		marsalled, _ := json.Marshal(res)
		fmt.Println(res["hey"])
		r.Write(marsalled)
		return nil
	})

	g.Post("/path", func(r http.ResponseWriter, req *http.Request) error {
		r.Write([]byte("This is a post :P"))
		res := make([]byte, 20*1024)
		req.Body.Read(res)
		fmt.Println(string(res))
		return nil
	})

	g.Get("/hey", func(r http.ResponseWriter, req *http.Request) error {
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
			fmt.Println(err)
		}
		r.Write(heyPage)
		return nil
	})

	g.ListenAndServe(":8088")
}
