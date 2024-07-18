package goster

import (
	"errors"
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
// If there're, ParseUrl will construct a Params struct and populate Meta.Query.Params
//
// If there aren't any, ParseUrl will return the error that occurred
//
// The `url` string reference that is passed in will have the parameters stripped in either case
func (m *Meta) ParseUrl(url string) (err error) {
	paramValues := make(map[string]string, 0)
	paramPattern := regexp.MustCompile(`\?.+(\/)?`)

	defer func() {
		m.Query = paramValues
	}()

	params := paramPattern.FindString(url)
	params = strings.Trim(params, "/?")

	if len(params) == 0 {
		err = errors.New("no query parameters")
		return
	}

	for _, v := range strings.Split(params, "&") {
		query := strings.Split(v, "=")

		if len(query) == 1 {
			continue
		}

		paramValues[query[0]] = query[1]
	}

	return
}

func (m *Meta) ParseDynPath(reqURL, dynPathURL string) (err error) {
	cleanPath(&reqURL)
	cleanPath(&dynPathURL)
	dynPaths, err := matchDynPathValue(dynPathURL, reqURL)

	if err != nil {
		return
	}

	for _, dynPath := range dynPaths {
		m.Path[dynPath.path] = dynPath.value
	}

	return
}
