package gottp_client

import (
	"errors"
	"regexp"
	"strings"
)

// Pass in a url and see if there're parameteres in it. If there are parseParams will return a Params struct and will remove them from the url that's passed in, if there aren't any, parseParams will return an error describing the error that occured and will still clean the url
func parseParams(url *string) (Params, error) {
	p_params := regexp.MustCompile(`\?.+(\/)?`)
	p_url := regexp.MustCompile(`^(\/\w+)+(\?)?`)

	defer func() {
		*url = strings.Trim(p_url.FindString(*url), "?")
	}()

	params := p_params.FindString(*url)
	params = strings.Trim(params, "/?")

	if len(params) == 0 {
		return Params{}, errors.New("no query params")
	}

	c_params := make(map[string]string)

	for _, v := range strings.Split(params, "&") {
		query := strings.Split(v, "=")

		if len(query) == 1 {
			continue
		}

		c_params[query[0]] = query[1]
	}

	return Params{values: c_params}, nil
}

func matchDynamicRoute(full string, dyn string) (DynamicRoute, error) {
	dyn_pattern := regexp.MustCompile(`\:\w+`)
	full_pattern := regexp.MustCompile(`^\w+`)

	loc_identifier := dyn_pattern.FindStringIndex(dyn)
	identifier := strings.Trim(dyn_pattern.FindString(dyn), ":")

	var identifier_value string
	route := DynamicRoute{}
	if len(full) >= loc_identifier[0] {
		identifier_value = full_pattern.FindString(full[loc_identifier[0]:])
	} else {
		return route, errors.New("paths don't match")
	}

	replaced_dyn := dyn_pattern.ReplaceAllString(dyn, identifier_value)

	if replaced_dyn == full {
		route.DynPath = dyn
		route.FullPath = full
		route.Identifier = identifier
		route.IdentifierValue = identifier_value
		return route, nil
	} else {
		return route, errors.New("paths don't match")
	}
}

// Deprecated
func parsePath(url string) []string {
	r := strings.Split(url, "/")
	temp_r := make([]string, 0)

	for _, v := range r {
		if len(v) != 0 {
			temp_r = append(temp_r, v)
		}
	}

	return temp_r
}
