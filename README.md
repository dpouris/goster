# Goster
[![GoDoc](https://godoc.org/github.com/gomarkdown/markdown?status.svg)](https://pkg.go.dev/github.com/dpouris/goster)
[![Go Report Card](https://goreportcard.com/badge/github.com/dpouris/goster)](https://goreportcard.com/report/github.com/dpouris/goster)
[![License](https://img.shields.io/github/license/dpouris/goster)](https://github.com/dpouris/goster/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/dpouris/goster)




Goster is a package that can be used to create servers and make API routes. It provides an abstraction on top of the built in http package to get up and running in no time.
-
<br>


## **INSTALLATION**

```shell
$ go get -u github.com/dpouris/goster
```
<br>

## **EXAMPLE**

```go
g := Goster.NewServer()

g.UseGlobal(func(ctx *Goster.Ctx) error {
	fmt.Println("global middleware")
	return nil
})

g.Get("/", func(ctx *Goster.Ctx) error {
	index, err := ioutil.ReadFile("./hey.html")

	if err != nil {
		fmt.Println(err)
	}
	ctx.ResponseWriter.Write(index)
	return nil
})

g.Get("db/", func(ctx *Goster.Ctx) error {
	// Query Param
	name, _ := ctx.Meta.Get("yourName")
	res := struct {
		Greet string `json:"greet"`
		Name string `json:"name"`
	}{
		Greet: "Hello",
		Name: name,
	}

	ctx.ResponseWriter.NewHeaders(map[string]string{
		"Content-Type": "application/json",
	}, 200)
	ctx.ResponseWriter.JSON(res)

	return nil
})

g.Get("fake_db/item/:id", func(ctx *Goster.Ctx) error {
	itemID, exists := ctx.Meta.Get("id")

	if !exists {
		itemID = ""
	}

	ctx.ResponseWriter.JSON(struct{
		itemID string `json:"itemID"`
	}{
		itemID
	})

	return nil
})

g.ListenAndServe(":8088")

```
<br>

# **USAGE**

### **New Server**
```go
g := Goster.NewServer()
```

### **GET**
```go
g.Get("/path", func(ctx *Goster.Ctx) error {
	// Handler logic
})
```

### **POST (Dynamic path)**
```go
g.Post("/path/:id", func(ctx *Goster.Ctx) error {
		// Handler logic
	})
```

### **ListenAndServe**
```go
g.ListenAndServe(":8000") //Pass in whatever port is free
```

### **Global Middleware**
```go
g.UseGlobal(func(ctx *Goster.Ctx) error {
    // middleware logic
	})
```
<br>

### **Path specific Middleware**
```go
g.Use("/path", func(ctx *Goster.Ctx) error {
    // middleware logic
	})
```
<br>

## **LOGGING**

By default Goster handles all incoming requests and Logs the info on the Logs field. On the example bellow I create a new instance of Goster server and supply `Goster.Logger` to the Log functions.
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

### **All logs**

You can access all the logs on the `Goster.Logs` field.

```go
g.Get("/logs", func(ctx *Goster.Ctx) error {
		log_map := make(map[int]any, len(g.Logs))

		for i, v := range g.Logs {
			log_map[i] = v
		}

		err := r.JSON(log_map)

		if err != nil {
			Goster.LogError(err.Error(), g.Logger)
		}
		return nil
	})
```

 - ### Sample Response

	```json
	{
		"0": "[GET] ON ROUTE /hey",
		"1": "[GET] ON ROUTE /logs"
	}	// Logs are stored in the Logs field of Goster instance
	```
