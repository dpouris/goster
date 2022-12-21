package goster

import (
	"log"
	"os"
	"sync"
)

type engine struct {
	Goster  *Goster
	startUp sync.Once
	config  any
}

// Init will run only once and set all the necessary fields for our one and only Goster instance
func (e *engine) Init() *Goster {
	initial := func() {
		logger := log.New(os.Stdout, "[SERVER] - ", log.LstdFlags)
		methods := make(map[string]map[string]Route)
		methods["GET"] = make(map[string]Route)
		methods["POST"] = make(map[string]Route)
		methods["PUT"] = make(map[string]Route)
		methods["PATCH"] = make(map[string]Route)
		methods["DELETE"] = make(map[string]Route)
		e.Goster = &Goster{Routes: methods, Middleware: make(map[string][]RequestHandler), Logger: logger}
	}

	// should set up config in here

	e.startUp.Do(initial)

	return e.Goster
}
