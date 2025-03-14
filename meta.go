package goster

import (
	"regexp"
	"strings"
)

type Meta struct {
	Query Params
	Path  Path
}

type DynamicPath struct {
	path  string
	value string
}

type Params map[string]string

type Path map[string]string

// Get tries to find if `id` is in the URL's Query Params
//
// If the specified `id` isn't found `exists` will be false
func (p *Params) Get(id string) (value string, exists bool) {
	value, exists = (*p)[id]
	return
}

// Get tries to find if `id` is in the URL's as a Dynamic Path Identifier
//
// If the specified `id` isn't found `exists` will be false
func (p *Path) Get(id string) (value string, exists bool) {
	value, exists = (*p)[id]
	return
}

// Pass in a `url` and see if there're parameters in it
//
// If there're, ParseQueryParams will construct a Params struct and populate Meta.Query.Params
//
// If there aren't any, ParseQueryParams will return
//
// The `url` string reference that is passed in will have the parameters stripped in either case
func (m *Meta) ParseQueryParams(url string) {
	paramValues := make(map[string]string, 0)
	paramPattern := regexp.MustCompile(`\?.+(\/)?`)

	defer func() {
		m.Query = paramValues
	}()

	params := paramPattern.FindString(url)
	params = strings.Trim(params, "/?")

	if len(params) == 0 {
		return
	}

	for _, v := range strings.Split(params, "&") {
		query := strings.Split(v, "=")

		if len(query) == 1 {
			paramValues[query[0]] = ""
			continue
		}

		paramValues[query[0]] = query[1]
	}
}

func (m *Meta) ParseDynamicPath(url, urlPath string) {
	cleanPath(&url)
	cleanPath(&urlPath)
	dynamicPaths, isDynamic := matchDynamicPath(url, urlPath)

	if !isDynamic {
		return
	}

	for _, dynamicPath := range dynamicPaths {
		m.Path[dynamicPath.path] = dynamicPath.value
	}
}
