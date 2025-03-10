# Getting Started with Goster

This guide will help you set up a basic web server using Goster. You’ll go from installation to having a running server that responds to HTTP requests.

## Installation

Make sure you have **Go 1.18+** installed on your system. Initialize a Go module for your project if you haven’t already:

```bash
go mod init myapp   # replace 'myapp' with your module name
```

Then add Goster to your project:

```bash
go get -u github.com/dpouris/goster
```

This will download the Goster module and update your `go.mod`. You’re now ready to use Goster in your code.

## Basic Setup

Let’s create a simple web server with one route. Create a file `main.go` with the following:

```go
package main

import "github.com/dpouris/goster"

func main() {
    // Initialize a new Goster server
    g := goster.NewServer()

    // Define a route for GET /
    g.Get("/", func(ctx *goster.Ctx) error {
        return ctx.Text("Welcome to Goster!")
    })

    // Start the server on port 8080
    g.Start(":8080")
}
```

**Run the server:** 

```bash
go run main.go
```

Open your browser (or use `curl`) and visit `http://localhost:8080`. You should see the response **“Welcome to Goster!”** displayed. Congratulations – you’ve just set up a Goster server!

## What’s Happening?

- We created a Goster server with `goster.NewServer()`. This gives us an instance `g` that will handle HTTP requests.
- We added a route using `g.Get("/")`. The first argument is the path (`"/"` for the root). The second argument is a **handler function** that Goster will call when a request comes in for that path. Our handler function uses `ctx.Text` to send a plain-text response.
- Finally, `g.Start(":8080")` starts an HTTP server on port 8080 and begins listening for requests. Under the hood, this uses Go’s `http.ListenAndServe`, passing Goster’s router as the handler.

When you visited the URL, Goster received the request, matched it to the `/` route, and executed your handler, which wrote “Welcome to Goster!” back to the client.

## Secure Server with TLS

Goster also supports running an HTTPS server using TLS. For example, set up your certificate and key files and start the server with:

```go
package main

import "github.com/dpouris/goster"

func main() {
    g := goster.NewServer()
    // Replace with the actual paths to your certificate and key
    g.StartTLS(":8443", "path/to/cert.pem", "path/to/key.pem")
}
```

## Next Steps

Now that you have a basic server running, you can start adding more routes and functionality:

- Define more routes (e.g., a `/hello` route) – see the [Routing](Routing.md) documentation.
- Add middleware for logging or authentication – see [Middleware](Middleware.md).
- Serve static files (like images or CSS) – see [Static Files](Static_Files.md).
- Set up HTML templates for more complex responses – see [Templates](Templates.md).

Goster is designed to scale from this simple example to a structure with multiple routes and middleware. Continue reading the documentation for guidance on each topic, or check out the `examples/` directory in the GitHub repository for additional sample programs.