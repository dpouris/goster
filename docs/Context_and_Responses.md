# Context and Response Handling in Goster

Goster’s handlers revolve around the **context** (`Ctx`), which encapsulates the request and response. Understanding the `Ctx` structure and its methods will help you effectively read incoming data and write outgoing responses.

## The Ctx Structure

When Goster receives a request, it creates a `Ctx` object and passes a pointer to it into your handler. The `Ctx` provides the following important fields and methods:

- **Request Data:**  
  - `ctx.Request` – the raw `*http.Request`. You can use this to read headers, the request body, etc., just as you would in any Go `net/http` handler.  
  - `ctx.Query` – a helper to access query parameters. It behaves like a map of string to string for URL query strings. Use `ctx.Query.Get("key")` to retrieve a parameter value.  
  - `ctx.Path` – a helper to access path parameters (from dynamic routes). Use `ctx.Path.Get("paramName")` to get the value of a URL parameter.

- **Response Tools:**  
  - `ctx.Response` – a wrapper around `http.ResponseWriter` that lets you write back to the client. In most cases, you won’t use `ctx.Response` directly; instead, you use methods like `ctx.Text`, `ctx.JSON`, `ctx.Template`, which under the hood write to `ctx.Response`. However, `ctx.Response` is available if you need advanced control (e.g., streaming data or setting headers manually).  
  - `ctx.Text(content string) error` – sends a plain text response. It will set the `Content-Type` to `text/plain; charset=utf-8` and write the provided text. Example: `ctx.Text("Hello")` will send “Hello” to the client.  
  - `ctx.JSON(data interface{}) error` – sends a JSON response. It will marshal the given data to JSON (using Go’s encoding/json by default) and set `Content-Type: application/json`. Example: `ctx.JSON(map[string]int{"status": 200})` sends `{"status":200}`. If JSON marshaling fails (e.g., due to an unsupported data type), this method will return an error.  
  - `ctx.Template(name string, data interface{}) error` – renders an HTML template (previously loaded by `TemplateDir`) and sends it. It sets `Content-Type: text/html; charset=utf-8`. Example: `ctx.Template("home.gohtml", user)` will fill the `home.gohtml` template with `user` data and send the result. Errors can occur if the template is not found or fails to execute.  

- **Meta Information and Logs:**  
  - `ctx.Meta` – holds internal metadata like the `Path` and `Query` maps. In most cases you won’t interact with `ctx.Meta` directly, but it’s where Goster stores parsed parameters.  
  - `ctx.Logs` – (if accessible) references the server’s log storage. Actually, logs are stored in the `Goster` server (`g.Logs`), not directly in `ctx`. But every request automatically logs a basic entry like `"[GET] ON ROUTE /example"` to `g.Logs`. If needed, you could correlate `ctx` to the server to fetch logs, but typically you use `g.Logs` via an admin route as shown in examples.

## Reading Request Data

**Query Parameters:** As mentioned, use `ctx.Query.Get("name")`. This returns two values: the value and a boolean `exists`. If `exists` is false, the parameter was not present in the URL. Example:

```go
q, ok := ctx.Query.Get("q")
if ok {
    // use q
} else {
    // parameter not provided
}
```

Internally, when the request comes in, Goster parses the raw query string into the `ctx.Meta.Query` map, so `ctx.Query.Get` is just a convenience accessor for that map.

**Path Parameters:** Use `ctx.Path.Get("param")` similarly. Path params are captured from dynamic routes. If a route is not dynamic or the param name is wrong, `exists` will be false. Example:

```go
id, ok := ctx.Path.Get("id")
if ok {
    // id contains the path segment value
}
```

Path params are parsed during routing, right before your handler is called. Goster populates `ctx.Meta.Path` with the values, which `ctx.Path.Get` accesses.

**Headers:** Use `ctx.Request.Header.Get("Header-Name")` to retrieve header values. For example, `ctx.Request.Header.Get("Content-Type")` or custom headers like `Authorization`. Goster doesn’t wrap header access — you use the standard `http.Request` methods.

**Body:** To read the request body (for POST/PUT, etc.), you can use `ctx.Request.Body`. For instance, if you expect JSON input, you might do:

```go
var input SomeStruct
err := json.NewDecoder(ctx.Request.Body).Decode(&input)
if err != nil {
    // handle JSON parse error
}
```

Remember to handle cases where the body might be empty or an error. Goster doesn’t automatically parse the body for you; you have full control via `ctx.Request`.

**Form data:** If you’re dealing with form submissions (traditional HTML forms or `multipart/form-data` for file uploads), you can use `ctx.Request.ParseForm()` or `ctx.Request.ParseMultipartForm()` and then access `ctx.Request.Form` or `ctx.Request.MultipartForm`. Again, this is using Go’s standard library capabilities directly on the `http.Request`.

## Writing Responses

The simplest way to respond is to use the helper methods:

- `ctx.Text`: for plaintext.
- `ctx.JSON`: for JSON.
- `ctx.Template`: for HTML content.

Each of these will automatically set a proper HTTP status code if not already set. By default, if you haven’t written any headers yet, writing through these methods will result in an implicit 200 OK status (unless an error occurs during JSON marshaling or template execution, in which case you should handle that by perhaps setting a 500).

**Custom Status Codes:** If you need to set a status code (like 201 Created, 204 No Content, 400 Bad Request, etc.), you have two options:
1. Use `ctx.Response.WriteHeader(code)` before writing the body. For example: 
   ```go
   ctx.Response.WriteHeader(http.StatusCreated)
   ctx.Text("resource created")
   ``` 
   Or for JSON:
   ```go
   ctx.Response.WriteHeader(http.StatusBadRequest)
   ctx.JSON(map[string]string{"error": "bad input"})
   ```
   Make sure to call `WriteHeader` *before* the helper method, because methods like `ctx.JSON` will call `Write` on the response, which implicitly sends a 200 if no status has been set yet.

2. Use Goster’s `NewHeaders` utility if available via `ctx.Response`. The `Response.NewHeaders(h map[string]string, status int)` function can set multiple headers and a status at once. For instance: 
   ```go
   ctx.Response.NewHeaders(map[string]string{
       "Content-Type": "application/json",
   }, http.StatusAccepted)
   ctx.Response.Write([]byte(`{"status": "accepted"}`))
   ```
   However, using `ctx.JSON` after `NewHeaders` isn’t advisable because `ctx.JSON` would try to set `Content-Type` again. Typically, stick to one approach: either manual or using helpers.

**No Response / Empty Response:** If your handler doesn’t need to send anything (for example, an endpoint that just consumes data and returns 204 No Content), you can do:
```go
ctx.Response.WriteHeader(http.StatusNoContent)
return nil
```
And not call any of the ctx helper methods to write a body. The client will get a 204 with an empty body.

**Errors in Handlers:** If you return an error from a handler, Goster will log it (to stdout or `g.Logs`) but it will not automatically send an error to the client. It’s up to your application design how to handle errors. A common pattern is to have middleware catch errors and format a response, or simply always respond within the handler and ensure you don’t return error unless you’ve handled it. For instance, you might do:

```go
g.Get("/data", func(ctx *goster.Ctx) error {
    data, err := loadData()
    if err != nil {
        ctx.Response.WriteHeader(http.StatusInternalServerError)
        ctx.Text("Internal Server Error")
        return nil  // we handled the error by responding to client
    }
    return ctx.JSON(data)  // happy path
})
```

Here we return `nil` after writing the error response so that upstream doesn’t attempt further handling. On the success path, we use `ctx.JSON`. This approach ensures the client always gets a response.

## Low-Level Access

Because `ctx.Response` embeds `http.ResponseWriter`, you can use all standard methods on it:
- `ctx.Response.Header().Set("X-Custom-Header", "value")` to set custom headers.
- `ctx.Response.Write(bytes)` to write raw bytes (the helper methods ultimately call this).
- `ctx.Response.WriteHeader(statusCode)` to set the HTTP status.

And `ctx.Request` is a normal `http.Request`, so you can use:
- `ctx.Request.URL.Path`, `ctx.Request.URL.Query()`, etc., if you prefer manual parsing.
- `io.ReadAll(ctx.Request.Body)` or streaming reads from the body for large payloads.
- `ctx.Request.Context()` if you need to observe cancellation (e.g., if the client disconnects, `ctx.Request.Context().Done()` will be signaled).

Goster’s context does **not** use Go’s `context.Context` for request values; instead, it offers the `Query` and `Path` maps. This is a deliberate design for simplicity. If you have existing code expecting `context.Context` with values (like from `net/http`), you can still retrieve that via `ctx.Request.Context()`.

## Examples

**Text response example:**

```go
g.Get("/ping", func(ctx *goster.Ctx) error {
    return ctx.Text("pong")
})
```
Client gets "pong" (Content-Type text/plain, Status 200).

**JSON response example:**

```go
g.Get("/api/time", func(ctx *goster.Ctx) error {
    now := time.Now()
    data := map[string]string{"time": now.Format(time.RFC3339)}
    return ctx.JSON(data)
})
```
Client gets `{"time": "2025-03-07T14:00:00Z"}` with `Content-Type: application/json`.

**HTML template example:**

```go
g.Get("/welcome/:name", func(ctx *goster.Ctx) error {
    name, _ := ctx.Path.Get("name")
    // Template "welcome.gohtml" expects a struct with Name field
    return ctx.Template("welcome.gohtml", struct{ Name string }{Name: name})
})
```
Client gets an HTML page with the name filled in.

**Setting a status code example:**

```go
g.Post("/items", func(ctx *goster.Ctx) error {
    // ... create item ...
    ctx.Response.WriteHeader(201)
    return ctx.Text("Created")
})
```
Client gets a "Created" message with HTTP 201 status.

## Conclusion

The `Ctx` gives you access to everything you need for request/response lifecycle:
- Use `ctx.Request` (and its Query and Path helpers) to **read input**.
- Use `ctx.Response` (and Text/JSON/Template helpers) to **write output**.

This design is similar to other Go frameworks (like Gin’s Context or Fiber’s Ctx) but stays close to the standard library, which makes it easy to integrate existing Go code. With this understanding, you can now handle virtually any kind of request in Goster – serving files, APIs, or web pages. Happy building!