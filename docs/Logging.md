# Logging in Goster

Logging is important for monitoring your application’s behavior and debugging issues. Goster includes a simple logging utility that captures events and also stores a history of requests. This document describes how logging works in Goster and how you can use it.

## How Goster Logs Requests

By default, Goster logs each incoming HTTP request automatically. When any request is handled, Goster records a log entry of the form:
```
[<METHOD>] ON ROUTE <path>
```
for example:
```
[GET] ON ROUTE /users/42
```
These entries are stored in the `Logs` slice on the Goster server (`g.Logs`). They are also printed to the console (stdout) if you run your app in a terminal, via the standard library `log.Logger`.

This means you have two ways to access request logs:
1. **In memory** – `g.Logs` (a slice of strings) contains recent log entries. You could use this to build an admin endpoint to fetch logs or for testing.
2. **In console output** – by default, Goster uses `log.Print` behind the scenes for these entries, so they appear in your application’s standard output.

## Using the Logger for Custom Messages

Goster’s logging utility is accessible through package functions that accept the server’s logger:
- `goster.LogInfo(message, g.Logger)`
- `goster.LogWarning(message, g.Logger)`
- `goster.LogError(message, g.Logger)`

These let you log custom messages at different levels (info, warning, error). All of them ultimately write to `g.Logger` (which is a `*log.Logger`). By default, `g.Logger` is configured to write to os.Stdout with a date/time prefix.

Example:

```go
g := goster.NewServer()
goster.LogInfo("Server started successfully", g.Logger) 

g.Get("/compute", func(ctx *goster.Ctx) error {
    goster.LogInfo("Compute endpoint hit", g.Logger)
    // ... do work ...
    if somethingUnexpected {
        goster.LogWarning("Unexpected condition encountered", g.Logger)
    }
    return ctx.Text("done")
})
```

In the above:
- When the server starts, we log an info message.
- Each time `/compute` is called, we log an info. If an unusual condition occurs, we log a warning.

The output in the console will look like (with timestamps prepended by Go’s log package):
```
2025/03/07 16:45:10 INFO  - Server started successfully
2025/03/07 16:47:05 INFO  - Compute endpoint hit
2025/03/07 16:47:05 WARN  - Unexpected condition encountered
```

Notice the format: by default, Goster prefixes the level (`INFO`, `WARN`, `ERROR`) in the log message. This is done by the `LogInfo/LogWarning/LogError` functions internally.

These messages are also added to the `g.Logs` slice. So `g.Logs` would contain:
```json
[
  "[INFO] Server started successfully",
  "[INFO] Compute endpoint hit",
  "[WARN] Unexpected condition encountered",
  ...
]
```
(The exact format might be slightly different, but conceptually, the log level is included.)

## Accessing Logs Programmatically

Because `g.Logs` holds the log entries, you can expose them via an endpoint for debugging. A simple example from the repository’s usage:

```go
g.Get("/logs", func(ctx *goster.Ctx) error {
    logMap := make(map[int]string, len(g.Logs))
    for i, entry := range g.Logs {
        logMap[i] = entry
    }
    return ctx.JSON(logMap)
})
``` 

This will return a JSON object where each key is an index and the value is the log entry. For instance:
```json
{
  "0": "[GET] ON ROUTE /",
  "1": "[GET] ON ROUTE /logs"
}
```
which shows that the root path `/` was accessed, and then the `/logs` path (to retrieve logs) was accessed.

You can modify this approach as needed:
- You might want to paginate or limit logs if the list grows large.
- You could filter by log level (for example, only errors) by scanning `g.Logs` for entries containing `ERROR`.

## Customizing Logging

As of now, Goster’s logging is relatively simple and **not fully configurable** (the README even marks “Configurable Logging” as a TODO feature). This means:
- The format of the log messages is fixed in the `LogInfo/Warning/Error` functions.
- The output destination of `g.Logger` is stdout by default. If you want to change it (say to log to a file), you could replace `g.Logger` with your own `log.Logger` instance after creating the server. For example:
  ```go
  g := goster.NewServer()
  file, _ := os.Create("app.log")
  g.Logger = log.New(file, "", log.LstdFlags)
  ```
  This would make all Goster logs go to `app.log` instead of the console. Make sure to handle errors and close the file appropriately in a real application.

- There is no built-in log rotation or level filtering. If you need more sophisticated logging (like only enabling debug logs in dev, etc.), you might integrate another logging library or simply conditionalize your calls to `LogInfo/Warning/Error`.

## Logging Middleware vs. Built-in Logging

If you prefer, you can also use middleware to log requests. For example, a global middleware that logs `ctx.Request.Method` and `ctx.Request.URL.Path`. This gives you full control over format and where it’s logged (you could use a third-party logging library inside the middleware). If you go that route, you might want to disable or ignore Goster’s built-in request logging. While you can’t turn it off easily, you could just not use `g.Logs` or ignore its output. Alternatively, clearing `g.Logs` periodically will remove old entries if you solely rely on your own logging.

For most basic uses, however, the built-in logging is sufficient to trace what endpoints are being hit and to sprinkle some custom logs for events.

## Example: Error Logging

If an error occurs in your handler and you catch it, you can use `goster.LogError` to record it:

```go
g.Post("/upload", func(ctx *goster.Ctx) error {
    err := handleUpload(ctx.Request)
    if err != nil {
        goster.LogError("Upload failed: "+err.Error(), g.Logger)
        ctx.Response.WriteHeader(http.StatusInternalServerError)
        return ctx.Text("Upload failed")
    }
    return ctx.Text("Upload successful")
})
```

This will log an ERROR with the error message. In your logs you might see:
```
2025/03/07 17:00:00 ERROR - Upload failed: file too large
```
Meanwhile, the client gets a 500 response with the text "Upload failed". By logging the error, you have a record server-side to investigate later.

## Summary

- Goster automatically logs all requests (method and route) to an internal list and stdout.
- Use `goster.LogInfo`, `LogWarning`, `LogError` for your own log messages in handlers or elsewhere .
- Retrieve logs via `g.Logs` for debugging or exposing through an API if needed.
- The logging system is simple; for advanced needs, you can replace or augment it with custom middleware or a custom logger.

Logging is essential for understanding your running application. Even though Goster’s logging is basic, it provides the hooks you need to implement a solid logging strategy for your service.