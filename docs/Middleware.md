# Middleware in Goster

Middleware allows you to execute code *before* or *after* your route handlers, enabling cross-cutting features like logging, authentication, error handling, etc. Goster provides a simple way to add middleware functions either globally (for all routes) or for specific routes.

## What is Middleware?

In web frameworks, middleware is a function that sits in the request handling chain. It can inspect or modify the request/response, and decide whether to continue to the next handler or stop the chain (for example, to return an error or redirect).

In Goster, a middleware is simply a function with the same signature as a handler: `func(ctx *goster.Ctx) error`. Middleware runs before the main handler for a route, in the order they were added.

## Global Middleware

**Global middleware** runs on every request, regardless of the path or HTTP method. This is useful for things like logging every request or enforcing site-wide security checks.

To add a global middleware, use the `UseGlobal` method:

```go
g := goster.NewServer()

g.UseGlobal(func(ctx *goster.Ctx) error {
    // Example: simple request logger
    println("Received request:", ctx.Request.Method, ctx.Request.URL.Path)
    return nil  // continue to next handler (or next middleware/route)
})
```

You can call `UseGlobal` multiple times to add multiple middleware. They will execute in the order added. If a middleware returns a non-nil error, Goster will consider the request handling failed at that point (you might handle this by logging or sending an error response).

Common use cases for global middleware:
- Logging requests (as in the example).
- Setting up common response headers (like security headers).
- Authentication checks for all routes (if not many public routes).

## Route-Specific Middleware

Sometimes you only want middleware on certain routes or groups of routes. Use `Use(path, middleware...)` to attach middleware to a specific path (or prefix). 

```go
// Middleware only for paths under /admin
g.Use("/admin", func(ctx *goster.Ctx) error {
    if !isUserAdmin(ctx) {
        ctx.Response.WriteHeader(403) // Forbidden
        ctx.Text("Forbidden")
        return fmt.Errorf("unauthorized")  // stop further handling
    }
    return nil
})
```

In this snippet, any request to a path that starts with `/admin` will go through the middleware. If the function returns an error (as in the case of a non-admin user), the main handler for the route will not run. If it returns nil, Goster will proceed to the next middleware (if any) or the final handler.

**How specific does the path need to be?** 

The `Use` method uses simple prefix matching for the path you provide:
- If you provide an exact path (e.g., `"/dashboard"`), it will apply to that path’s routes.
- If you provide a prefix (e.g., `"/admin"` as above), it will apply to all routes that have that prefix (like `/admin`, `/admin/settings`, `/admin/users/123`). 

Internally, Goster stores middleware in a map keyed by the path or prefix you give. On each request, after the main handler is chosen, Goster runs:
1. All global middleware (in order added).
2. Any middleware whose path prefix matches the request path (in the order they were added, but note that only one middleware function list will match – exactly or by prefix).

**Order and Matching Detail:** If you added both `Use("/admin", ...)` and `Use("/admin/settings", ...)`, a request to `/admin/settings` would trigger *only* the middleware associated with the exact `/admin/settings` match (because Goster’s implementation stores the middleware for exact path separately from prefix matches). In general, use either broad prefixes or exact paths to avoid confusion. Goster does not currently support middleware for multiple arbitrary patterns or regex.

## Middleware Execution Flow

For a given request, the flow is:

1. **Global middleware** – all functions added via `UseGlobal` run, in the order added.
2. **Route-specific middleware** – if the request path matches a key used in `Use(path, ...)`, those middleware functions run (in order).
3. **Route handler** – finally, the main handler for the route executes.

All middleware and the handler share the same `ctx` (context) for the request, so they can communicate via `ctx`. For example, a logging middleware could set a value in `ctx.Meta` or add to `ctx.Logs` that a later middleware or handler could use.

If any middleware returns an error, Goster will still call the remaining middleware in that list but will skip the main route handler (see the use of `defer route.Handler(ctx)` in the internal implementation). It’s up to you how to handle errors: you might simply log them, or a middleware could send an early response (like `ctx.Text("Forbidden")` as shown above).

## Examples

**1. Logging middleware (global):**

```go
g.UseGlobal(func(ctx *goster.Ctx) error {
    start := time.Now()
    err := ctx.Next() // Note: Goster doesn't have ctx.Next(); this is a conceptual example
    duration := time.Since(start)
    fmt.Printf("%s %s completed in %v\n", ctx.Request.Method, ctx.Request.URL.Path, duration)
    return err
})
```

*(The above uses a conceptual `ctx.Next()` to illustrate timing around the handler, but Goster’s middleware does not currently provide a built-in next mechanism since it runs all middleware automatically.)* In practice, to measure time you could record `start` in the first middleware and in a **deferred** function, calculate the duration after the handler runs (since the handler runs after middleware).

**2. Authentication middleware (specific path):**

```go
g.Use("/api", func(ctx *goster.Ctx) error {
    token := ctx.Request.Header.Get("Authorization")
    if !validateToken(token) {
        ctx.Response.WriteHeader(401)
        ctx.Text("Unauthorized")
        return fmt.Errorf("auth failed")
    }
    return nil
})
```

This will protect all routes beginning with `/api`. If the token is missing or invalid, it returns 401 Unauthorized and stops the request from reaching the actual API handler.

**3. Middleware stacking:** 

You can attach multiple middleware to the same path or globally. They will run in the order added:

```go
g.UseGlobal(mw1)
g.UseGlobal(mw2)
// mw1 runs, then mw2, then the route handler.

g.Use("/reports", mwA)
g.Use("/reports", mwB)
// For any /reports request, mwA runs then mwB then the handler.
```

Be mindful of the order, especially if one middleware’s behavior affects another. For example, if `mwA` modifies something that `mwB` relies on.

## Best Practices

- **Keep middleware focused:** Each middleware should ideally do one thing (logging, auth check, etc.). This makes it easier to compose and reuse.
- **Performance:** Remember that global middleware runs for every request. Don’t put extremely heavy processing in middleware (or guard it so it only runs when needed).
- **Error Handling:** Decide how you want to handle errors in middleware. One pattern is to have middleware *not* return errors upward, but rather handle errors internally (for example, by writing a response directly). Another approach is to use a final middleware that catches errors from previous ones or the handler. Goster doesn’t have a special error-catching middleware mechanism, but you can structure your code to check for error returns after `ListenAndServe` (since Goster’s ListenAndServe currently doesn’t propagate handler errors, you might handle errors inside the middleware itself).

- **ctx values:** You can use the `ctx.Meta` map if you need to pass information from middleware to handlers (for example, user info after authentication). Alternatively, since `ctx.Request` is available, you could use Go’s standard `Context` (in `ctx.Request.Context()`) to store values, but that’s typically not necessary for simple cases.

Middleware can greatly enhance your application by separating concerns. With Goster’s `UseGlobal` and `Use` methods, you have the flexibility to apply middleware broadly or narrowly as needed. Continue to [Static Files](Static_Files.md) or other docs to explore more Goster features.