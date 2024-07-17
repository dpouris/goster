# Goster 🚀
[![GoDoc](https://godoc.org/github.com/gomarkdown/markdown?status.svg)](https://pkg.go.dev/github.com/dpouris/goster)
[![Go Report Card](https://goreportcard.com/badge/github.com/dpouris/goster)](https://goreportcard.com/report/github.com/dpouris/goster)
[![License](https://img.shields.io/github/license/dpouris/goster)](https://github.com/dpouris/goster/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/dpouris/goster)

Welcome to **Goster**, the lightweight and efficient web framework for Go! 🌟

## Why Goster?

- 🚀 **Fast and Lightweight**: Build with simplicity in mind, Goster provides a minimalistic abstraction on top of the built-in http package.
- 📊 **Intuitive API**: Easy-to-use API that simplifies web development without sacrificing flexibility.
- 🛠 **Extensible Middleware**: Seamlessly add middleware to enhance your application's functionality.
- 🔍 **Dynamic Routing**: Effortlessly handle both static and dynamic routes.
- 🧪 **Configurable Logging** (TODO): Powerful and customizable logging to keep track of your application's activity.

## Installation

Install Goster using `go get`:

```sh
go get -u github.com/dpouris/goster
```

## Quick Start

Create your first Goster server:

```go
package main

import (
    "github.com/dpouris/goster"
)

func main() {
    g := goster.NewServer()

    g.Get("/", func(ctx *goster.Ctx) error {
        ctx.Text("Hello, Goster!")
        return nil
    })

    g.ListenAndServe(":8080")
}
```

## **Usage**

- Create a new server:
	```go
	g := Goster.NewServer()
	```

- Add a `GET` Route:
	```go
	g.Get("/path", func(ctx *Goster.Ctx) error {
		// Handler logic
	})
	```

- Add a `Dynamic Route`:
	```go
	g.Post("/path/:id", func(ctx *Goster.Ctx) error {
			// Handler logic
		})
	```

- Add `Global Middleware`:
	```go
	g.UseGlobal(func(ctx *Goster.Ctx) error {
		// middleware logic
		})
	```

- Add `Path specific Middleware`:
	```go
	g.Use("/path", func(ctx *Goster.Ctx) error {
		// middleware logic
		})
	```

- Logging:

	By default Goster handles all incoming requests and Logs the info on the Logs field. On the example below, we craete a new instance of Goster server and supply `Goster.Logger` to the Log functions.
	```go
	import Goster "github.com/dpouris/goster/goster"

	func main() {
		g := Goster.NewServer()

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

- All logs:

	You can access all the logs on the `Goster.Logs` field.

	```go
	g.Get("/logs", func(ctx *Goster.Ctx) error {
			log_map := make(map[int]any, len(g.Logs))
			for i, v := range g.Logs {
				log_map[i] = v
			}
			r.JSON(log_map)

			return nil
		})
	```

	Sample Response

	```json
	// Logs are stored in the Logs field of Goster instance
	{
		"0": "[GET] ON ROUTE /hey",
		"1": "[GET] ON ROUTE /logs"
	}
	```

## Examples

Check out these examples to get started quickly:

- [Basic Setup](examples/basic/example_basic.go)
- [Query Parameters](examples/query_params/example_query_params.go)
- [Dynamic Routes](examples/dynamic_routes/example_dynamic_route.go)
- [JSON Responses](examples/json_response/example_json_response.go)
- [HTML Templates](examples/html_template/example_html_template.go)
- [Middleware Usage](examples/middleware/example_middleware.go)

## Contributing

I welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more information.

## License

Goster is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.