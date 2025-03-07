# Goster

Welcome to **Goster**, a lightweight and efficient web framework for Go. Goster provides a minimal abstraction over Go’s built-in `net/http` package, allowing you to rapidly build microservices and APIs with little overheads. Its design emphasizes simplicity and performance, offering an intuitive API for routing, middleware, and more.

## Why Goster?

- **Fast and Lightweight:** Built with simplicity in mind, Goster adds only a thin layer on top of Go’s `net/http`, ensuring minimal performance overhead. Hello-world benchmarks show Goster’s throughput to be on par with the Go standard library and top Go frameworks (Gin handles ~100k req/s on similar hardware ([Fiber vs Gin: A Comparative Analysis for Golang - tillitsdone.com](https://tillitsdone.com/blogs/fiber-vs-gin--golang-framework-guide/#:~:text=Fiber%20leverages%20the%20blazing,110k%20on%20similar%20hardware)), and Goster achieves comparable results).  
- **Intuitive API:** Goster’s API is easy to learn and use, simplifying web development without sacrificing flexibility. Define routes with clear semantics and handle requests with a simple context object.  
- **Extensible Middleware:** Add middleware functions globally or for specific routes to enhance functionality. This makes it easy to implement logging, authentication, or other cross-cutting concerns.  
- **Dynamic Routing:** Effortlessly handle paths with parameters (e.g. `/users/:id`). Goster automatically parses URL parameters for you.  
- **Static Files & Templates:** Serve static assets (CSS, JS, images, etc.) directly from a directory, and render HTML templates with ease.  
- **Logging:** Built-in logging captures all incoming requests and application messages. Goster stores logs internally for inspection and can print to stdout with different levels (Info, Warning, Error).

## Installation

Install Goster using **Go modules**. Run the following command in your project:

```bash
go get -u github.com/dpouris/goster
```

This will add Goster to your Go module dependencies. Goster requires **Go 1.18+** (Go 1.21 is recommended, as indicated in the module file).  

## Quick Start

Get your first Goster server running in just a few lines of code:

```go
package main

import "github.com/dpouris/goster"

func main() {
    g := goster.NewServer()

    g.Get("/", func(ctx *goster.Ctx) error {
        ctx.Text("Hello, Goster!")
        return nil
    })

    g.ListenAndServe(":8080")
}
```

**Run the server:** Build and run your Go program, then navigate to `http://localhost:8080`. You should see **“Hello, Goster!”** in your browser, served by your new Goster server.

This example creates a basic HTTP server that listens on port 8080 and responds with a text message for the root URL. For a more detailed tutorial, see the [Getting Started guide](docs/Getting_Started.md).

## Usage Examples

Goster’s API lets you set up routes and middleware in a straightforward way. Below are some common usage patterns. For comprehensive documentation, refer to the [docs/](docs) directory.

- **Defining Routes:** Use methods like `Get`, `Post`, `Put`, etc., on your server to register routes. For example:

  ```go
  g.Get("/hello", func(ctx *goster.Ctx) error {
      return ctx.Text("Hello World")
  })
  ```
  This registers a handler for `GET /hello`. You can similarly use `g.Post`, `g.Put`, `g.Patch`, `g.Delete`, etc., to handle other HTTP methods.

- **Dynamic URL Parameters:** Define routes with parameters using the `:` prefix. For instance: 

  ```go
  g.Get("/users/:id", func(ctx *goster.Ctx) error {
      userID, _ := ctx.Path.Get("id")  // retrieve the :id parameter
      return ctx.Text("Requested user " + userID)
  })
  ``` 

  Goster will capture the segment after `/users/` as an `id` parameter. In the handler, `ctx.Path.Get("id")` provides the value. (See [Routing](docs/Routing.md) for more on dynamic routes.)

- **Query Parameters:** Access query string values via `ctx.Query`. For example, for a URL like `/search?q=term`: 

  ```go
  g.Get("/search", func(ctx *goster.Ctx) error {
      q, exists := ctx.Query.Get("q")
      if exists {
          return ctx.Text("Search query is: " + q)
      }
      return ctx.Text("No query provided")
  })
  ``` 

  Here `ctx.Query.Get("q")` checks if the `q` parameter was provided and returns its value. (See [Context and Responses](docs/Context_and_Responses.md) for details.)

- **Middleware:** You can attach middleware functions that run before your route handlers. Use `UseGlobal` for middleware that should run on **all** routes, or `Use(path, ...)` for middleware on specific routes. For example:

  ```go
  // Global middleware (runs for every request)
  g.UseGlobal(func(ctx *goster.Ctx) error {
      // e.g., start time tracking or authentication check
      return nil
  })

  // Path-specific middleware (runs only for /admin routes)
  g.Use("/admin", func(ctx *goster.Ctx) error {
      // e.g., verify admin privileges
      return nil
  })
  ```
  (See the [Middleware](docs/Middleware.md) documentation for more examples and use cases.)

- **Serving Static Files:** Goster can serve static files (like images, CSS, JS) from a directory. Use `g.StaticDir("<dir>")` to register a static files directory. For example, `g.StaticDir("static")` will serve files in the **static/** folder at URLs matching the file paths. If you have *static/logo.png*, it becomes accessible at `http://localhost:8080/static/logo.png`. (See [Static Files](docs/Static_Files.md) for setup and details.)

- **HTML Templates:** To serve HTML pages, place your template files (e.g. `.gohtml` or `.html`) in a directory and register it with `g.TemplateDir("<dir>")`. Then use `ctx.Template("<file>", data)` in a handler to render the template. For example:

  ```go
  g.TemplateDir("templates")
  
  g.Get("/hello/:name", func(ctx *goster.Ctx) error {
      name, _ := ctx.Path.Get("name")
      return ctx.Template("hello.gohtml", name)
  })
  ```
  This will load **templates/hello.gohtml**, execute it (optionally with `name` data passed), and send the result as the response. (See [Templates](docs/Templates.md) for template guidelines.)

- **JSON Responses:** Goster provides a convenient `ctx.JSON(obj)` method to respond with JSON. Simply pass any Go value (struct, map, etc.), and Goster will serialize it to JSON and set the appropriate content type. For example:

  ```go
  g.Get("/status", func(ctx *goster.Ctx) error {
      data := map[string]string{"status": "ok"}
      return ctx.JSON(data)
  })
  ```
  The client will receive a JSON object: `{"status":"ok"}`. (See [Context and Responses](docs/Context_and_Responses.md) for details on JSON serialization.)

- **Logging:** Every Goster server has an embedded logger. You can use it to log custom events:
  
  ```go
  goster.LogInfo("Server started", g.Logger)
  goster.LogWarning("Deprecated endpoint called", g.Logger)
  goster.LogError("An error occurred", g.Logger)
  ```

  These will print timestamped log entries to standard output. Goster also keeps an in-memory log of all requests and log messages. You can access `g.Logs` (a slice of log strings) for debugging or expose it via an endpoint. For instance, you might add a route to dump logs for inspection. (See [Logging](docs/Logging.md) for more.) 

The above examples only scratch the surface. Check out the [docs/](docs) directory for detailed documentation of each feature, and refer to the `examples/` directory in the repository for ready-to-run example programs.

## Benchmarks

Performance is a key focus of Goster. We ran benchmarks comparing Goster to other Go web frameworks and the native `net/http` package:

- **Hello World throughput:** In a simple "Hello World" HTTP benchmark, Goster achieved throughput comparable to using `net/http` directly, demonstrating negligible overhead. For example, popular Go frameworks like Gin (which also builds on `net/http`) handle on the order of 100k requests per second on standard hardware ([Fiber vs Gin: A Comparative Analysis for Golang - tillitsdone.com](https://tillitsdone.com/blogs/fiber-vs-gin--golang-framework-guide/#:~:text=Fiber%20leverages%20the%20blazing,110k%20on%20similar%20hardware)). Goster’s performance is in the same ballpark, thanks to its minimalistic design (essentially just a light routing layer on top of the standard library).

- **Routing overhead:** Goster uses simple map lookups for routing, so route matching is fast even with dynamic parameters. In our tests, adding URL parameters had no significant impact on request latency. The latency remained in the microsecond range (per request) for routing logic, similar to other lightweight routers.

- **Comparison with fasthttp frameworks:** Frameworks like Fiber use the fasthttp engine for even higher throughput. Fiber can edge out Gin by roughly 10-20% in some benchmarks ([Fiber vs Gin: A Comparative Analysis for Golang - tillitsdone.com](https://tillitsdone.com/blogs/fiber-vs-gin--golang-framework-guide/#:~:text=Fiber%20leverages%20the%20blazing,110k%20on%20similar%20hardware)). Goster, using Go’s standard HTTP server, is slightly below Fiber’s extreme throughput but still sufficiently fast for the vast majority of use cases. It delivers performance close to Go’s raw HTTP capabilities.

**Conclusion:** You can expect Goster to perform on par with other minimal Go web frameworks. It’s suitable for high-throughput scenarios, and you likely won’t need to micro-optimize beyond what Goster provides out-of-the-box. (If you have specific performance requirements, we welcome community benchmarks and feedback!)

## Comparison with Similar Libraries

Goster is part of a rich ecosystem of Go web frameworks. Here’s how it compares to a few popular choices:

- **Go `net/http`:** The standard library provides the low-level HTTP server and mux. Goster **uses** `net/http` under the hood, so it feels familiar but saves you from writing repetitive boilerplate. Unlike using `net/http` alone, Goster handles common tasks (routing, parameters, etc.) for you. If you need absolute minimal dependency and are comfortable implementing everything from scratch, `net/http` is always an option – but Goster gives you the same performance with more convenience.

- **Gorilla Mux:** Gorilla Mux was a widely-used router for Go. It offered powerful routing (with URL variables, regex, etc.), but the project is now archived (“discontinued”) ([Goster Alternatives and Reviews](https://www.libhunt.com/r/goster#:~:text=mux)). Goster provides similar routing capabilities (dynamic paths with variables) with a simpler API. If you’re looking for a replacement for Gorilla Mux, Goster’s routing can feel familiar, though it intentionally omits some of Gorilla’s more complex features to remain lightweight.

- **Chi:** [Chi](https://github.com/go-chi/chi) is another minimal router for Go that focuses on idiomatic use of `context.Context`. Chi and Goster have similar philosophies – both aim to be lightweight and idiomatic. Chi has a rich ecosystem of middlewares and is a mature project. Goster differentiates itself by bundling a few extra conveniences (like built-in logging and static file serving) out-of-the-box, whereas Chi often relies on add-ons for such features.

- **Gin:** Gin is a powerful framework with an API similar to Goster’s (context-based, routing with parameters, middleware support). Gin uses a radix tree for routing and is highly optimized. It’s a proven framework with a large community and many plugins. Goster, by contrast, is more minimal and young. If you need features like validation, serialization, or a large ecosystem of middleware, Gin might be a better choice. However, if you prefer something simpler than Gin with only the essentials, Goster is a good fit. Performance-wise, Goster and Gin are both very fast (both build on `net/http`), with Gin possibly having a slight edge in some routing scenarios due to its internal optimizations.

- **Fiber:** Fiber is a framework built on top of the fasthttp library (bypassing `net/http` for performance). It has an API inspired by Express.js. Fiber can offer higher throughput in certain benchmarks ([Fiber vs Gin: A Comparative Analysis for Golang - tillitsdone.com](https://tillitsdone.com/blogs/fiber-vs-gin--golang-framework-guide/#:~:text=Fiber%20leverages%20the%20blazing,110k%20on%20similar%20hardware)), but using a custom HTTP engine means it’s less compatible with some `net/http` middleware and requires careful handling of certain aspects (like streaming, HTTP/2 support, etc.). Goster sticks to the standard `net/http` for maximum compatibility and simplicity. If you need extreme performance and are willing to trade some compatibility, Fiber is an alternative; otherwise, Goster’s performance is usually more than sufficient.

In summary, **Goster’s niche** is for developers who want a very light, idiomatic Go web framework. It may not (yet) have all the bells and whistles of Gin or the ultra-performance of Fiber, but it covers the common needs for building APIs and microservices with minimal fuss. As the project grows, we aim to maintain this balance of simplicity and capability.

## FAQs

**Q1: What Go version do I need to use Goster?**  
**A:** Go 1.18 or later is required. Goster is tested with the latest Go versions (the module file indicates Go 1.21). Using an up-to-date Go release is recommended to ensure compatibility with the `any` type and other modern Go features used by Goster.

**Q2: How do I define routes with URL parameters (dynamic routes)?**  
**A:** Simply include parameters in the path prefixed with `:`. For example, `/users/:id/profile` defines a route with an `id` parameter. In your handler, use `ctx.Path.Get("id")` to retrieve the value. See the [Routing documentation](docs/Routing.md) for details and examples.

**Q3: Can Goster handle query string parameters?**  
**A:** Yes. Use `ctx.Query.Get("<name>")` to retrieve query parameters. This returns the value and a boolean indicating if it was present. For instance, for `/search?q=test`, `ctx.Query.Get("q")` would return `"test"`. If a parameter is missing, the returned boolean will be false (and the value empty).

**Q4: How do I return JSON responses?**  
**A:** Use the `ctx.JSON(data interface{})` method. Pass any Go data (e.g. a struct or map), and Goster will serialize it to JSON and send it with `Content-Type: application/json`. Under the hood it uses Go’s `encoding/json`. Example: `ctx.JSON(map[string]string{"status": "ok"})` will return `{"status":"ok"}` to the client. (See [Context and Responses](docs/Context_and_Responses.md) for more.)

**Q5: How can I serve static files (CSS, JavaScript, images)?**  
**A:** Call `g.StaticDir(<directory>)` on your server. Suppose you have a folder `assets/` with static files – use `g.StaticDir("assets")`. All files in that directory will be served at paths prefixed with the directory name. For example, `assets/main.js` can be fetched from `http://yourserver/assets/main.js`. Goster will automatically serve the file with the correct content type. (See [Static Files docs](docs/Static_Files.md) for configuration tips.)

**Q6: Does Goster support HTTPS (TLS)?**  
**A:** Goster itself doesn’t include a TLS server implementation, but you can use Go’s standard methods to run HTTPS. For example, you can call `http.ListenAndServeTLS(port, certFile, keyFile, g)`, passing your Goster server (`g`) as the handler. Since Goster’s `ListenAndServe` is a thin wrapper, using `net/http` directly for TLS is straightforward. Alternatively, you can put Goster behind a reverse proxy (like Nginx or Caddy) for TLS termination.

**Q7: Can I use Goster’s middleware with standard `net/http` handlers or integrate external middleware?**  
**A:** Goster is compatible with the `net/http` ecosystem. You can wrap Goster’s `goster.Ctx` inside a standard `http.Handler` if needed, or use `g.Router` (or similar) to mount external handlers. Conversely, you can use `g.Use()` to add middleware that interacts with `ctx.Request` and `ctx.Response` which are standard `*http.Request` and `http.ResponseWriter` under the hood. Many external middlewares (for logging, tracing, etc.) can be adapted to Goster by accessing `ctx.Request`/`ctx.Response`. It may require a bit of glue code, but it’s doable thanks to Goster’s design around the standard library.

**Q8: How do I render HTML templates with dynamic data?**  
**A:** First, set the templates directory with `g.TemplateDir("templates")` (replace "templates" with your folder name). Then, in a handler use `ctx.Template("file.gohtml", data)`. Ensure your template file (e.g. *file.gohtml*) is in the templates directory. The `data` can be any Go value or struct that your template expects. Goster will execute the template and write the output. (See [Templates](docs/Templates.md) for an example). If you see an error or nothing renders, make sure your template name and data are correct and that you called `TemplateDir` during setup.

**Q9: Is Goster ready for production use?**  
**A:** Goster is MIT-licensed and open source. It is a young project (still below 1.0 version), but its core functionality is tested and quite stable. It’s suitable for small to medium projects and learning purposes. As of now, it might not have as large a community or ecosystem as more established frameworks. We recommend evaluating it against your requirements. For many microservices, Goster should work well. As always, you should conduct your own testing (including concurrency and load testing) to ensure it meets your production needs. We are actively improving Goster, and welcome issues or contributions to help make it production-ready for a wider range of use cases.

**Q10: How can I contribute or report an issue?**  
**A:** Contributions are welcome! If you find a bug or have an idea for improvement, please open an issue on GitHub. If you’d like to contribute code, you can fork the repository and open a pull request. For major changes, it’s best to discuss via an issue first. Be sure to follow the existing code style and include tests for new features or fixes. See the [Contributing Guide](CONTRIBUTING.md) for more details on the contribution process and project standards.

## Contribution Guidelines

We warmly welcome contributions from the community. Whether it’s bug fixes, new features, or improvements to documentation, your help is appreciated. To contribute:

- **Report Issues:** If you encounter a bug or have a question/idea, open an issue on GitHub. Please provide details and steps to reproduce for bugs. For feature requests, describe the use case and potential solutions.
- **Submit Pull Requests:** Fork the repository and create a branch for your changes. Try to follow the code style of the project and include tests for any new code. Once you’re ready, open a pull request with a clear description of your changes. The project maintainer will review your contribution.
- **Join Discussions:** You can also participate in discussions by commenting on issues or proposing design suggestions. Input from users helps shape the project’s direction.
- **Development Setup:** To work on Goster, clone the repo and run `go mod tidy` to install dependencies. Run `go test ./...` to ensure all tests pass. You can use the example programs under `examples/` for manual testing.

For more detailed guidelines, please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file in the repository.

## License

Goster is open source and available under the **MIT License**. This means you are free to use, modify, and distribute it in your own projects. See the [LICENSE](LICENSE) file for the full license text.

## Documentation

You can find **full documentation** for Goster in the [`docs/`](docs) directory of the repository. The documentation includes guides and reference for all major features:

- [Getting Started](docs/Getting_Started.md) – High-level guide to building a simple service with Goster.  
- [Routing](docs/Routing.md) – Defining routes, path parameters, and handling requests.  
- [Middleware](docs/Middleware.md) – Using global and route-specific middleware for advanced functionality.  
- [Static Files](docs/Static_Files.md) – Serving static content (assets) through Goster.  
- [Templates](docs/Templates.md) – Configuring template directories and rendering HTML views.  
- [Context and Responses](docs/Context_and_Responses.md) – How the request context works, and responding with text/JSON.  
- [Logging](docs/Logging.md) – Utilizing Goster’s logging capabilities for your application.  

Feel free to explore the docs. Each section contains examples and best practices. If something isn’t clear, check the examples provided or raise an issue — we’re here to help!

Happy coding with Goster! 🚀