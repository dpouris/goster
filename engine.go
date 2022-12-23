package goster

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)

type Engine struct {
	Goster  *Goster
	startUp sync.Once
	Config  *Config
}

type Config struct {
	BaseStaticDir string
	FilePaths     []string
}

var engine = Engine{}

// init will run only once and set all the necessary fields for our one and only Goster instance
func (e *Engine) init() *Goster {
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
	e.DefaultConfig()

	e.startUp.Do(initial)

	return e.Goster
}

// Set the default config settings for the engine
func (e *Engine) DefaultConfig() {
	e.Config = &Config{}

	err := e.SetTemplateDir("templates")

	if err != nil {
		cwd, _ := os.Getwd()
		e.Config.BaseStaticDir = cwd
		return
	}

}

func (e *Engine) SetTemplateDir(d string) (err error) {
	var staticDir string

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	staticDir = path.Join(cwd, d)

	e.Config.BaseStaticDir = staticDir

	files, err := os.ReadDir(staticDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not a valid directory\n", d)
		return
	}

	for _, de := range files {
		e.Config.FilePaths = append(e.Config.FilePaths, path.Join(e.Config.BaseStaticDir, de.Name()))
	}

	return
}
