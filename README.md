# GOTTP SERVER

Gottp is an HTTP server that can be used to serve static files and make simple API routes. It provides an abstraction on top of the built in http package to get up and running in no time.
-
<br>

# USAGE

### **New Server**
```go
g := Gottp.Server()
```

### **GET**
```go
g.Get("/path", func(r http.ResponseWriter, req *http.Request) error {
		// code
	})
```

### **POST**
```go
g.Post("/path", func(r http.ResponseWriter, req *http.Request) error {
		// code
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

## **EXAMPLE**

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	Gottp "github.com/dpouris/gottp-server"
)

func main() {
	g := Gottp.Server()

	g.AddGlobalMiddleware(func(r http.ResponseWriter, req *http.Request) error {
		fmt.Println("Executes for every request before the handler")
		return nil
	})

	g.Get("/path", func(r http.ResponseWriter, req *http.Request) error {
        // A map that is marshalled to JSON
		res := make(map[string]any, 10)
		res["foo"] = "bar"
		marsalled, _ := json.Marshal(res)

        // Write the JSON to the resonse body
		r.Write(marsalled)
		return nil
	})

	g.Post("/path", func(r http.ResponseWriter, req *http.Request) error {
        res := make([]byte, 20*1024)
        // Read the request body and write it to the res byte slice
		req.Body.Read(res)

		r.Write([]byte("Thank you for this tasty dump"))
		return nil
	})

	g.Get("/hey", func(r http.ResponseWriter, req *http.Request) error {
        // Loads the HTML file
		heyPage, err := ioutil.ReadFile("./examples/hey.html")

		if err != nil {
            Gottp.LogError(err.Error())
		}

        // Write the HTML to the response body
		r.Write(heyPage)
		return nil
	})

    // Listen on port 8088
	g.ListenAndServe(":8088")
}

```

## **LOGGING**
```go
import Gottp "github.com/dpouris/gottp-server"

func main() {
    // Logs to stdout
    Gottp.LogInfo("This is an info message")
    Gottp.LogWarning("This is an warning message")
    Gottp.LogError("This is an error message")
}
```
```shell
// OUTPUT

2022/06/07 11:45:40 INFO  - This is an info message
2022/06/07 11:45:40 WARN  - This is an warning message
2022/06/07 11:45:40 ERROR - This is an error message
```

## **INSTALLATION**

```shell
$ go get -u github.com/dpouris/gottp-server
```