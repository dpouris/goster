# Template Rendering in Goster

Goster supports rendering HTML templates so you can serve dynamic content (e.g., HTML pages) in addition to raw text or JSON. It uses Go’s standard html/template or text/template packages under the hood. This document will guide you through setting up template rendering.

## Setting the Template Directory

First, organize your template files (often with `.html` or `.gohtml` extensions) in a directory. Common practice is to have a folder like `templates/` in your project.

Tell Goster where your templates live by using `TemplateDir`:

```go
g := goster.NewServer()
err := g.TemplateDir("templates")
if err != nil {
    log.Fatal("Failed to load templates:", err)
}
```

`TemplateDir` scans the specified directory (and subdirectories) for template files. It will record each template’s relative path in an internal map for quick access later. If the directory does not exist, Goster will create it (and log an info message)  – this is convenient if you run the app and the template folder is empty, but you plan to add files later.

**Note:** If a template file name appears more than once (e.g., two files with the same name in different subfolders), Goster may log a warning or error to avoid duplicates. Make sure template file names or relative paths are unique in the directory structure.

## Creating a Template

Templates are text files with placeholders for dynamic data. For example, create a file `templates/hello.gohtml`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Hello</title>
</head>
<body>
    <h1>Hello, {{.}}!</h1>
</body>
</html>
```

This is a simple template that expects a single value (denoted by `{{.}}`) which it will insert into the HTML.

## Rendering a Template in a Handler

Use `ctx.Template(name string, data interface{}) error` in your route handler to render a template. The `name` should match the file name (or relative path) of the template you want to render, and `data` is whatever you want to pass into the template (it becomes `.` inside the template).

For example:

```go
g.Get("/hello/:name", func(ctx *goster.Ctx) error {
    name, _ := ctx.Path.Get("name")
    return ctx.Template("hello.gohtml", name)
})
```

In this handler, when someone visits `/hello/John`, Goster will load the `hello.gohtml` template, execute it with the data `"John"`, and send the resulting HTML. The client will see an HTML page with “Hello, John!” as an `<h1>`.

Behind the scenes, when you call `ctx.Template("hello.gohtml", data)`, Goster will:
1. Find the compiled template associated with `"hello.gohtml"` (from the files loaded by `TemplateDir`).
2. Execute the template with the provided `data`.
3. Write the output to `ctx.Response` with `Content-Type: text/html`.

If the template name isn’t found in the map (e.g., you gave the wrong name or forgot to call `TemplateDir`), Goster will likely return an error or write a message to stderr. Ensure that the template was recognized at startup. The `TemplateDir` function prints to stderr if a directory is invalid, so check your logs if nothing is rendering.

## Template Data and Functions

The `data` you pass to `ctx.Template` can be any Go value:
- It can be a string (as in the example above).
- It can be a struct or map for more complex templates (e.g., passing a struct with multiple fields to use in the template).
- It can even be a slice or any other type; how you use it depends on your template content.

You can also define template functions, use template blocks, etc., by leveraging Go’s `html/template` capabilities. For instance, you might want to parse templates that include other files (partials) or define custom functions (like date formatting). Currently, Goster’s `TemplateDir` will load all files individually. If you need a more advanced setup (like combining templates or defining a base layout), you might manually use Go’s template package and incorporate that into Goster (for example, by storing a template instance in a global and writing to `ctx.Response`). However, for simple use cases, separate files per template as loaded by `TemplateDir` works fine.

## Example: Template with Struct Data

Imagine you have a template `templates/profile.gohtml`:

```html
<h1>User Profile</h1>
<p>Name: {{.Name}}</p>
<p>Age: {{.Age}}</p>
```

And a Go struct:

```go
type User struct {
    Name string
    Age  int
}
```

Your handler could be:

```go
g.Get("/users/:id/profile", func(ctx *goster.Ctx) error {
    // Fetch user data (this is just a stub example)
    user := User{Name: "Alice", Age: 30}
    return ctx.Template("profile.gohtml", user)
})
```

This will render the profile template, replacing `{{.Name}}` with "Alice" and `{{.Age}}` with 30.

## Template Reloading

Goster loads templates once at startup (when you call `TemplateDir`). If you edit template files while the server is running, Goster will not automatically reload them. You would need to restart the server to pick up changes. In development, this is fine (just restart on template edits). In production, you typically don’t change templates on the fly. If hot-reloading of templates is required, you’d have to implement a custom solution (which might involve calling `TemplateDir` again or managing your own template cache).

## Error Handling

If there’s an error during template execution (for example, a template syntax error, or the data doesn’t match what the template expects), `ctx.Template` will return an error. You should handle that error (perhaps by returning it from the handler, which could be picked up by middleware or result in a 500). When developing, check your console output for any errors that Goster’s template functions might log.

## Comparison to JSON/Text responses

Using `ctx.Template` is analogous to `ctx.JSON` or `ctx.Text`, except it’s for HTML content. Use it when you need to send structured HTML pages. For simple APIs that only send JSON or plain text, you won’t need templates. But if your microservice also serves a web UI or emails (in HTML format), templates can be very useful.

## Recap

- Use `g.TemplateDir("path/to/templates")` to load templates from disk.
- In handlers, call `ctx.Template("filename", data)` to render a template and send it to the client.
- Organize your templates and data so that they match (template placeholders correspond to fields in the data you pass).
- Remember to restart the server if template files change (in development).

With templates covered, you have a full spectrum of response options: plain text, JSON, and HTML. Next, you might want to read about [Logging](Logging.md) to see how to monitor your application’s behavior.
