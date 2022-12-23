package goster

import (
	"errors"
	"regexp"
	"strings"
)

// Pass in a url and see if there're parameters in it. If there are parseParams will return a Params struct and will remove them from the url that's passed in, if there aren't any, parseParams will return an error describing the error that occurred and will still clean the url
func parseParams(url *string) (Params, error) {
	paramsPtrn := regexp.MustCompile(`\?.+(\/)?`)
	pathPtrn := regexp.MustCompile(`^(\/\w+)+(\?)?`)

	defer func() {
		*url = strings.Trim(pathPtrn.FindString(*url), "?")
	}()

	params := paramsPtrn.FindString(*url)
	params = strings.Trim(params, "/?")

	if len(params) == 0 {
		return Params{}, errors.New("no query params")
	}

	paramMap := make(map[string]string)

	for _, v := range strings.Split(params, "&") {
		query := strings.Split(v, "=")

		if len(query) == 1 {
			continue
		}

		paramMap[query[0]] = query[1]
	}

	return Params{values: paramMap}, nil
}

func matchDynamicRoute(full string, dyn string) (route DynamicRoute, err error) {
	dynPattern := regexp.MustCompile(`\:\w+`)
	fullPattern := regexp.MustCompile(`^\w+`)

	identifierLoc := dynPattern.FindStringIndex(dyn)
	identifier := strings.Trim(dynPattern.FindString(dyn), ":")

	var identifierValue string
	if len(full) >= identifierLoc[0] {
		identifierValue = fullPattern.FindString(full[identifierLoc[0]:])
	} else {
		err = errors.New("paths don't match")
		return
	}

	replacedDyn := dynPattern.ReplaceAllString(dyn, identifierValue)

	if replacedDyn == full {
		route.DynPath = dyn
		route.FullPath = full
		route.Identifier = identifier
		route.IdentifierValue = identifierValue
		return
	} else {
		err = errors.New("paths don't match")
		return
	}
}

// Adds basic headers
func DefaultHeader(c *Ctx) {
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

func parsePath(path *string) {
	if *path == "/" {
		*path = ""
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	if (*path)[len(*path)-1] == '/' {
		*path = (*path)[:len((*path))-1]
	}

}

func cleanEmptyBytes(b *[]byte) {
	cleaned := []byte{}

	for _, v := range *b {
		if v == 0 {
			break
		}
		cleaned = append(cleaned, v)
	}
	*b = cleaned
}
