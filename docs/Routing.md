# Routing in Goster

Routing is at the core of Goster. It determines how incoming HTTP requests are matched to the code that handles them. This document explains how to define routes, handle path parameters, and understand Goster’s routing mechanisms.

## Defining Routes

You register routes on a Goster server by calling methods named after HTTP verbs:

- `g.Get(path, handler)`
- `g.Post(path, handler)`
- `g.Put(path, handler)`
- `g.Patch(path, handler)`
- `g.Delete(path, handler)`

Each of these methods takes a URL path and a **handler function**. The handler should have the signature `func(ctx *goster.Ctx) error`. The handler is called when a request with the corresponding HTTP method and path is received.

**Example – GET route:**

```go
g := goster.NewServer()

g.Get("/hello", func(ctx *goster.Ctx) error {
    return ctx.Text("Hello, world!")
})
``` 

In this example, any `GET` request to `/hello` will trigger the handler and respond with “Hello, world!”.

Under the hood, Goster stores routes in a routing table (a map). When you call `g.Get` (or any method), Goster adds an entry associating the path to your handler. If you try to register the same path twice for the same method, Goster will return an error to prevent duplicates.

## Dynamic Routes and Path Parameters

Often you need routes that capture parts of the URL (e.g., an ID in the path). Goster supports **dynamic routing** using path parameters. A path parameter is indicated by a `:` prefix in the route definition.

**Example – dynamic route:**

```go
g.Get("/users/:id", func(ctx *goster.Ctx) error {
    id, _ := ctx.Path.Get("id")
    return ctx.Text("User ID is " + id)
})
``` 

Here, `:id` in the path means that any value in that segment of the URL will match. For a request to `/users/42`, for instance, the handler will be called and `ctx.Path.Get("id")` will return `"42"`. 

You can use multiple parameters and static parts in one path. For example, `/users/:uid/books/:bid` would capture two parameters, `uid` and `bid`. In the handler, you’d retrieve both with `ctx.Path.Get("uid")` and `ctx.Path.Get("bid")`. The order and names should match what you put in the route pattern.

**Note:** Path parameters only match up to the next `/` or the end of the path. For instance, in `/files/:name.txt`, the parameter would include “.txt” as part of the value (because the dot is not a separator). Generally, define parameters between slashes, like `/files/:name`.

## Route Handlers and the Context

A route handler is a function with signature `func(ctx *goster.Ctx) error`. When a request comes in, Goster creates a new context (`ctx`) and passes it to your handler. This context contains:

- `ctx.Request` – the original `*http.Request`.
- `ctx.Response` – a response writer (through which you send output).
- `ctx.Path` – a map of path parameters for this request.
- `ctx.Query` – a map of query string parameters for this request.

Typically, you won’t need to access `ctx.Request` or `ctx.Response` directly for basic tasks, because Goster provides helper methods on `ctx` (like `ctx.Text`, `ctx.JSON`, etc.). But they are available if you need lower-level control.

**Writing responses:** A handler should return an `error`. If you encounter an error during processing, you can return it and handle it as needed (for example, you might use middleware to catch errors and return a JSON error response). If no error occurs, return `nil` (as in the examples above).

Goster doesn’t automatically do anything with returned errors yet (beyond logging them), so you can also handle errors inside the handler (e.g., send an HTTP 500). This might be improved in the future.

## Order of Routes and Priority

Goster’s routing is straightforward: each route is registered to an HTTP method and exact path (except dynamic segments which match variable text). When a request comes in, Goster checks the routing table for that HTTP method:

1. If an exact match is found for the path, it uses that handler.
2. If not, it tries to match dynamic routes. Goster will iterate through registered dynamic routes to find one that fits the request path. For each dynamic route pattern, it checks if the incoming path can be parsed to match that pattern. If a match is found, it stops searching further.
3. If no route matches, Goster will return a 404 Not Found (by default, simply an empty response with that status).

Currently, dynamic route matching in Goster is simple and checks patterns in the order they were added. It’s a good practice to avoid overly ambiguous patterns. For example, if you have both `/users/:id` and `/users/profile`, Goster might interpret `/users/profile` as matching the `:id` route if added first. In such cases, check ordering or adjust your route design (maybe use `/users/:id/profile` for profiles to include an ID).

## Adding Routes for Other Methods

We demonstrated `Get`. The same concept applies for other methods:

```go
g.Post("/users", func(ctx *goster.Ctx) error {
    // Create a new user
    return ctx.Text("User created")
})

g.Delete("/users/:id", func(ctx *goster.Ctx) error {
    // Delete user with given id
    return ctx.Text("User deleted")
})
```

Use the appropriate method name for the type of request you want to handle. If a client sends a request with a method that you haven’t defined, Goster will reply with a 405 Method Not Allowed for that path (assuming the path exists for a different method, otherwise 404).

## Wildcard Routing (Not supported)

Some frameworks support wildcards or catch-all parameters (like `/files/*path` to match `/files/any/arbitrary/path`). **Goster does not currently support wildcard segments.** Every route segment is either static or a single parameter. If you need to serve arbitrary nested paths, you might consider using the static file serving feature (see [Static Files](Static_Files.md)) or handle it in your handler by parsing `ctx.Request.URL` yourself.

## Summary

- Use `g.<Method>` to register routes for different HTTP methods.
- Include `:param` in the path to capture dynamic path parameters. Retrieve them in the handler with `ctx.Path.Get`.
- The `ctx` (context) passed to handlers provides request data and helper methods for responses.
- Goster matches routes by method and then by path, supporting dynamic segments. If no match, it returns 404 by default.
- Keep route patterns unambiguous to avoid confusion between static and dynamic routes.

With routing set up, you can now move on to processing requests (e.g., reading request bodies, returning JSON) which is covered in [Context and Responses](Context_and_Responses.md), or adding [Middleware](Middleware.md) to extend functionality before/after handlers.