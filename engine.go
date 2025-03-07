package goster

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Engine struct {
	Goster  *Goster
	startUp sync.Once
	Config  *Config
}

type Config struct {
	BaseTemplateDir string
	StaticDir       string
	TemplatePaths   map[string]string
	StaticFilePaths map[string]string
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
	e.Config = &Config{
		StaticDir:       "",
		BaseTemplateDir: "",
		TemplatePaths:   make(map[string]string, 0),
		StaticFilePaths: make(map[string]string, 0),
	}
}

// Sets the template directory to `d` relative to the path of the executable.
func (e *Engine) SetTemplateDir(path string) (err error) {
	templateDir, err := resolveAppPath(path)
	if err != nil {
		return
	}

	// if the given directory doesn't exist, create it and report it
	if ok, _ := pathExists(templateDir); !ok {
		os.Mkdir(templateDir, 0o711) // rwx--x--x (o+rwx) (g+x) (u+x)
		fmt.Printf("[ENGINE INFO] - given template path `%s` doesn't exist\n", path)
		fmt.Printf("[ENGINE INFO] - creating `%s`...\n", path)
	}
	e.Config.BaseTemplateDir = templateDir
	templatesMap := make(map[string]string)
	err = filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("cannot walk %s dir", templateDir)
		}
		if !d.IsDir() {
			relativePath, err := filepath.Rel(templateDir, path)
			if err != nil {
				return err
			}
			templatesMap[relativePath] = path
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not a valid directory\n", path)
		return
	}

	for templ := range templatesMap {
		fmt.Printf("[ENGINE INFO] - Recorded template `%s` -> %s\n", templ, templatesMap[templ])
		if !e.Config.AddTemplatePath(templ, templatesMap[templ]) {
			return fmt.Errorf("template `%s` already exists in `%s`", templ, e.Config.BaseTemplateDir)
		}
	}

	return
}

func (e *Engine) SetStaticDir(path string) (err error) {
	staticPath, err := resolveAppPath(path)
	if err != nil {
		return err
	}

	// if the given directory doesn't exist, create it and report it
	if ok, _ := pathExists(staticPath); !ok {
		os.Mkdir(staticPath, 0o711) // rwx--x--x (o+rwx) (g+x) (u+x)
		fmt.Printf("[ENGINE INFO] - given static path `%s` doesn't exist\n", path)
		fmt.Printf("[ENGINE INFO] - creating static path `%s`...\n", path)
	}

	e.Config.StaticDir = staticPath
	staticFileMap := make(map[string]string)
	// walk the directory and register a GET route for each file found
	err = filepath.WalkDir(staticPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// process only files (skip directories)
		if !d.IsDir() {
			// compute the route path relative to the static directory
			relPath, _ := filepath.Rel(staticPath, filePath)
			cleanPath(&relPath)

			staticFileMap[relPath] = filePath
		}

		return nil
	})

	for relPath := range staticFileMap {
		fmt.Printf("[ENGINE INFO] - Recorded static file `%s` -> %s\n", relPath, staticFileMap[relPath])
		if !e.Config.AddStaticFilePath(relPath, staticFileMap[relPath]) {
			return fmt.Errorf("static file `%s` already exists in `%s`", relPath, e.Config.StaticDir)
		}
	}

	return
}
func (c *Config) AddTemplatePath(relPath string, fullPath string) (added bool) {
	// check if path exists
	_, exists := c.TemplatePaths[relPath]

	if !exists {
		// add if it doesn't
		c.TemplatePaths[relPath] = fullPath
		added = true
		return
	}

	added = false
	return
}

func (c *Config) AddStaticFilePath(relPath string, fullPath string) (added bool) {
	// check if path exists
	_, exists := c.StaticFilePaths[relPath]

	if !exists {
		// add if it doesn't
		c.StaticFilePaths[relPath] = fullPath
		added = true
		return
	}

	added = false
	return
}
