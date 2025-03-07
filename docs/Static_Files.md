# Serving Static Files with Goster

Many web services need to serve static assets such as images, CSS files, JavaScript files, or downloadable content. Goster provides a convenient way to serve files from a directory on your filesystem through your web server.

## Setting a Static Directory

To tell Goster about your static files, use the `StaticDir` method on your server. This method takes the path to a directory on your system that contains static files.

```go
g := goster.NewServer()
err := g.StaticDir("static")
if err != nil {
    log.Fatal("Failed to set static dir:", err)
}
```

In this example, Goster will serve files from the `static` directory (relative to where your program is running). Suppose the directory structure is:

```
static/
├── css/
│   └── style.css
├── images/
│   └── logo.png
└── index.html
```

After calling `g.StaticDir("static")`:

- A request to `GET /static/index.html` will return the `static/index.html` file.
- A request to `GET /static/css/style.css` will return the `static/css/style.css` file.
- A request to `GET /static/images/logo.png` will return the `static/images/logo.png` file.

In general, `g.StaticDir("<dirname>")` makes the contents of that directory available under the URL path `/<dirname>/*`.

**Behind the scenes:** When you call `StaticDir`, Goster scans the directory and automatically registers routes for each fil e. Each file gets a corresponding route (usually a GET route) that serves the file’s content. Goster also attempts to set the appropriate Content-Type based on file extension (using an internal utility function `getContentType`).

If `StaticDir` returns an error, it likely means the directory doesn’t exist or couldn’t be read. Goster will print an error to stderr if, for example, the directory path is wrong or files can’t be opened. Ensure the path is correct and that your program has read access to the files.

## Accessing Static Files

Once `StaticDir` is set up, clients can retrieve the files by making requests to the corresponding URL. Typically, you’ll use this to serve front-end assets. For example, if you have an `index.html` as a single-page app entry point, you might have:

```go
g.StaticDir("static")
g.Get("/", func(ctx *goster.Ctx) error {
    // Redirect root to the main index page
    return ctx.Template("index.html", nil)  // or ctx.Text/ctx.File if you prefer
})
```

However, note that `ctx.Template` is for server-side templates, not static HTML files. If you just want to serve `index.html` as a static file, you could also place it in the static directory and let users access it via `/static/index.html`. Alternatively, use `ctx.Response` to manually serve files (but StaticDir handles this for you).

### Example

If your app’s HTML, CSS, JS are in a folder named **web**:

```go
g := goster.NewServer()
g.StaticDir("web")
g.ListenAndServe(":8080")
```

Now:
- `http://localhost:8080/web/` will list the files or require a specific file (depending on the exact behavior, typically you’d request an explicit file).
- `http://localhost:8080/web/app.js` serves the file `web/app.js` if it exists.
- You might configure your front-end build to output to the **web** folder, so all static assets are served by Goster.

## Security Considerations

Goster’s static file serving will expose all files in the directory you specify and its subdirectories. Be careful not to include sensitive files in that directory. For instance, do not point `StaticDir` at a directory that contains configuration files or private data. It’s best to keep a dedicated folder for public assets.

Also, Goster’s static serving is intended for convenience. For high-throughput static file serving or serving very large files, a dedicated static file server or CDN might be more appropriate. But for many applications (especially APIs that just need to serve a few static files for a frontend), Goster’s approach is sufficient.

## Disabling Static Serving

If you no longer want to serve static files, you could stop calling `StaticDir` or remove those routes. Currently, Goster doesn’t provide a method to *remove* a static directory once added. If needed, you would have to manage that logic (for example, by not calling `StaticDir` based on some config). Typically, though, you set it once at startup and leave it.

## Combining with Templates

It’s common to use static file serving alongside template rendering. Static files are for assets, whereas templates are for dynamic HTML generation. You can use both in Goster:

```go
g.StaticDir("static")       // for CSS, JS, images
g.TemplateDir("templates")  // for dynamic HTML templates

g.Get("/", func(ctx *goster.Ctx) error {
    return ctx.Template("home.gohtml", someData)
})
```

In this scenario:
- `/static/*` URLs serve files.
- The root `/` (and other non-static routes) can render templates that likely include links to those static assets (for example, `<link rel="stylesheet" href="/static/styles.css">` in the HTML).

## Conclusion

Serving static files in Goster is straightforward: just point `StaticDir` at your folder of assets. This allows you to build a simple web server that not only provides JSON APIs but also serves a web interface or documentation files, etc., without an additional server. For a more detailed understanding of how Goster registers static file routes, you can refer to the implementation in the source (Goster reads files and registers routes for each file).

Continue to [Templates](Templates.md) to learn how to serve dynamic content using Goster’s template functionality.