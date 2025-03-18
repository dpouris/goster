package goster

import (
	"path"
	"slices"
	"strings"
)

type element struct {
	Value, Type string
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

		switch currentElement.Type {
		case Dynamic:
			continue

		case Wildcard:
			nextElem := pathElements[i+1]
			if nextElem.Type == Static { // go until static element
				wildcardEnd := slices.Index(urlSlice[i:], nextElem.Value)
				if wildcardEnd == -1 { // doesn't match return
					return false
				}
				// skip i to wildcardEnd
				skip = wildcardEnd
				continue
			}
			return true

		case Static:
			if pathElements[i].Value != v { // doesn't match, return
				return false
			}
		}
	}

	return true
}

func constructElements(urlPath, routePath string) (pv []PathValues) {
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

		switch currentElement.Type {
		case Dynamic:
			pv = append(pv, PathValues{strings.TrimPrefix(pathElements[i].Value, ":"), strings.Split(v, "?")[0]})
			continue

		case Wildcard:
			// two cases:
			// 1. the wildcard match stops when a static element is hit
			// 2. the wildcard match spans the rest of urlPath

			// check if case 1 is true, theres a static element after the wildcard
			var wildcardPath []string
			nextElem := pathElements[i+1]
			if nextElem.Type == Static { // go until static element
				staticElemIdx := (slices.Index(urlSlice[i:], nextElem.Value)) + i
				if staticElemIdx <= 0 { // doesn't match return
					return
				}
				skip = staticElemIdx
				wildcardPath = urlSlice[i:staticElemIdx]
			} else { // spans the rest urlPath
				wildcardPath = urlSlice[i:]
			}

			if len(pathElements[i].Value) == 1 { // no identifier, skip
				continue
			}
			// the wildcard has an identifier, assign it
			p := path.Join(wildcardPath...)
			cleanPath(&p)
			pv = append(pv, PathValues{strings.TrimPrefix(pathElements[i].Value, "*"), p})
			continue

		case Static:
			if pathElements[i].Value != v { // doesn't match, return
				return
			}
		}
	}

	return
}

func constructPathElements(routeSlice []string) map[int]element {
	// idx -> pathElement
	elems := make(map[int]element, 0)
	for i, v := range routeSlice {
		var t string
		if strings.HasPrefix(v, "*") {
			t = Wildcard
		} else if strings.HasPrefix(v, ":") {
			t = Dynamic
		} else {
			t = Static
		}
		elems[i] = element{v, t}
	}

	return elems
}
