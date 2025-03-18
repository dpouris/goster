package goster

import (
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Adds basic headers
func DefaultHeader(c *Ctx) {
	// TODO: refactor
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

// cleanPath sanatizes a URL path. It removes suffix '/' if any and adds prefix '/' if missing
func cleanPath(path *string) {
	if len(*path) == 0 {
		*path = "/"
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	if len(*path) != 1 {
		*path = strings.TrimSuffix(*path, "/")
	}
}

// urlMatchesRoute checks URL path `urlPath` matches a Route path
// `routePath`. A Dynamic Route is a path string that has the following format: "path/anotherPath/:variablePathname" where `:variablePathname`
// is a catch-all identifier that matches any route with the same structure up to that point.
//
// Ex:
//
//	var ctx = ...
//	var url = "path/anotherPath/andYetAnotherPath"
//	var dynamicPath = "path/anotherPath/:identifier"
//	if !urlMatchesRoute(&ctx, url, dynamicPath) {
//			panic(...)
//	}
//
// The above code will not panic as the urlMatchesRoute will evaluate to `true`
func urlMatchesRoute(urlPath string, routePath string) bool {
	cleanPath(&urlPath)
	cleanPath(&routePath)
	urlSlice := strings.Split(urlPath[1:], "/")
	routeSlice := strings.Split(routePath[1:], "/")
	pathElements := constructPathElements(routeSlice) // idx -> pathElement

	if len(routeSlice) > len(urlSlice) && !strings.Contains(routePath, "*") { // doesn't match, return
		return false
	}

	// fmt.Printf("\nFor url slice: %v\nPath Elements: %v\n", urlSlice, pathElements)
	skip := 0
	for i, v := range urlSlice {
		if skip > 0 {
			skip -= 1
			continue
		}

		currentElement, exists := pathElements[i]
		if !exists { // doesn't match, return
			return false
		}

		switch currentElement.t {
		case TypeDynamic:
			continue

		case TypeWildcard:
			nextElem := pathElements[i+1]
			if nextElem.t == TypeStatic { // go until static element
				wildcardEnd := slices.Index(urlSlice[i:], nextElem.v)
				if wildcardEnd == -1 { // doesn't match return
					return false
				}
				// skip i to wildcardEnd
				skip = wildcardEnd
				continue
			}
			return true

		case TypeStatic:
			if pathElements[i].v != v { // doesn't match, return
				return false
			}
		}
	}

	return true
}

func findPathValues(urlPath, routePath string) (pv []PathValues) {
	urlSlice := strings.Split(urlPath[1:], "/")
	routeSlice := strings.Split(routePath[1:], "/")
	pathElements := constructPathElements(routeSlice) // idx -> pathElement
	pv = make([]PathValues, 0)

	skip := 0
	for i, v := range urlSlice {
		if skip > 0 {
			skip -= 1
			continue
		}

		currentElement, exists := pathElements[i]
		if !exists { // doesn't match, return
			return
		}

		switch currentElement.t {
		case TypeDynamic:
			pv = append(pv, PathValues{strings.TrimPrefix(pathElements[i].v, ":"), strings.Split(v, "?")[0]})
			continue

		case TypeWildcard:
			// two cases:
			// 1. the wildcard match stops when a static element is hit
			// 2. the wildcard match spans the rest of urlPath

			// check if case 1 is true, theres a static element after the wildcard
			var wildcardPath []string
			nextElem := pathElements[i+1]
			if nextElem.t == TypeStatic { // go until static element
				staticElemIdx := (slices.Index(urlSlice[i:], nextElem.v)) + i
				if staticElemIdx <= 0 { // doesn't match return
					return
				}
				skip = staticElemIdx
				wildcardPath = urlSlice[i:staticElemIdx]
			} else { // spans the rest urlPath
				wildcardPath = urlSlice[i:]
			}

			if len(pathElements[i].v) == 1 { // no identifier, skip
				continue
			}
			// the wildcard has an identifier, assign it
			p := path.Join(wildcardPath...)
			cleanPath(&p)
			pv = append(pv, PathValues{strings.TrimPrefix(pathElements[i].v, "*"), p})
			continue

		case TypeStatic:
			if pathElements[i].v != v { // doesn't match, return
				return
			}
		}
	}

	return
}

func constructPathElements(routeSlice []string) map[int]struct{ v, t string } {
	// idx -> pathElement
	elems := make(map[int]struct{ v, t string }, 0)
	for i, v := range routeSlice {
		var t string
		if strings.HasPrefix(v, "*") {
			t = TypeWildcard
		} else if strings.HasPrefix(v, ":") {
			t = TypeDynamic
		} else {
			t = TypeStatic
		}
		elems[i] = struct{ v, t string }{v, t}
	}

	return elems
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func resolveAppPath(dir string) (string, error) {
	fileDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine working directory for static dir %s\n", dir)
		return dir, err
	}

	// construct the full path to the static directory
	return path.Join(fileDir, dir), nil
}

func pathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	var pos int
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			if pos < i {
				b = append(b, s[pos:i]...)
			}
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
			pos = i + 1
		}
	}
	if pos < len(s) {
		b = append(b, s[pos:]...)
	}
	return string(b)
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}

// func cleanEmptyBytes(b *[]byte) {
// 	cleaned := []byte{}

// 	for _, v := range *b {
// 		if v == 0 {
// 			break
// 		}
// 		cleaned = append(cleaned, v)
// 	}
// 	*b = cleaned
// }
