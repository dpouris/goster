# Goster
[![GoDoc](https://godoc.org/github.com/gomarkdown/markdown?status.svg)](https://pkg.go.dev/github.com/dpouris/goster)
[![Go Report Card](https://goreportcard.com/badge/github.com/dpouris/goster)](https://goreportcard.com/report/github.com/dpouris/goster)
[![License](https://img.shields.io/github/license/dpouris/goster)](https://github.com/dpouris/goster/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/dpouris/goster)




Goster is a siple HTTP library that can be used to serve static files and make simple API routes. It provides an abstraction on top of the built in http package to get up and running in no time.
-
<br>


## **INSTALLATION**

```shell
$ go get -u github.com/dpouris/goster
```
<br>

## **EXAMPLE**

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	Goster "github.com/dpouris/goster"
)

func main() {
	g := Goster.Server()

	g.Get("/hey", func(r Goster.Res, req *Goster.Req) error {
        // Loads the HTML file
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
            Goster.LogError(err.Error(), g.Logger)
		}

        // Write the HTML to the response body
		r.Write(heyPage)
		return nil
	})

	g.Post("/hey", func(r Goster.Res, req *Goster.Req) error {
		// A map that is marshalled to JSON
		log_map := map[string]string{
			"hey": "you",
			"hello": "world",
		}

		// Write the JSON to the resonse body
		err := r.JSON(log_map)

		if err != nil {
            Goster.LogError(err.Error(), g.Logger)
		}

		return nil
	})

    // Listen on port 8088
	g.ListenAndServe(":8088")
}

```
<br>

# **USAGE**

### **New Server**
```go
g := Goster.Server()
```

### **GET**
```go
g.Get("/path", func(r Goster.Res, req *Goster.Req) error {
	// Handler logic
})
```

### **POST**
```go
g.Post("/path", func(r Goster.Res, req *Goster.Req) error {
		// Handler logic
	})
```

### **ListenAndServe**
```go
g.ListenAndServe(":8000") //Pass in whatever port is free
```

### **Global Middleware**
```go
g.AddGlobalMiddleware(func(r http.ResponseWriter, req *http.Request) error {
    // middleware logic
	})
```
<br>

## **LOGGING**

By default Goster handles all incoming requests and Logs the info on the Logs field. On the example bellow I create a new instance of Goster server and supply `Goster.Logger` to the Log functions.
```go
import Goster "github.com/dpouris/goster"

func main() {
	g := Goster.Server()

    // Logs to stdout
    Goster.LogInfo("This is an info message", g.Logger)
    Goster.LogWarning("This is an warning message", g.Logger)
    Goster.LogError("This is an error message", g.Logger)
}
```
```shell
// OUTPUT

2022/06/07 11:45:40 INFO  - This is an info message
2022/06/07 11:45:40 WARN  - This is an warning message
2022/06/07 11:45:40 ERROR - This is an error message
```

### **All logs**

You can access all the logs on the `Goster.Logs` field.

```go
g.Get("/logs", func(r Goster.Res, req *Goster.Req) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := r.JSON(log_map)

		if err != nil {
			Goster.LogError(err.Error(), g.Logger)
		}
		return nil
	})
```

 - ### Sample Response

	```json
	{
		"0": "[GET] ON ROUTE /hey",
		"1": "[GET] ON ROUTE /logs"
	}	// Logs are stored in the Logs field of Goster instance
	```
