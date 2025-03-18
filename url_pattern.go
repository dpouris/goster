package goster

import (
	"strings"
)

// SegmentType distinguishes different segment types
type SegmentType int

const (
	StaticSegment SegmentType = iota
	DynamicSegment
	WildcardSegment
)

// RouteSegment represents one part of a route
type RouteSegment struct {
	Raw  string      // The original string (e.g. "users", ":id", "*filepath")
	Type SegmentType // Whether it’s static, dynamic, or wildcard
	// For dynamic and wildcard segments, Name holds the identifier (without ':' or '*')
	Name string
}

// NewRouteSegment parses a single segment string
func NewRouteSegment(s string) RouteSegment {
	if len(s) > 0 {
		switch s[0] {
		case ':':
			return RouteSegment{Raw: s, Type: DynamicSegment, Name: s[1:]}
		case '*':
			// For a bare "*" we treat Name as empty and do not capture any parameter
			return RouteSegment{Raw: s, Type: WildcardSegment, Name: s[1:]}
		}
	}
	return RouteSegment{Raw: s, Type: StaticSegment}
}

// RoutePattern holds the parsed route
type RoutePattern struct {
	Segments []RouteSegment
}

// NewRoutePattern pre-parses a route string (e.g. "/users/:id" or "/static/*") into a RoutePattern
func NewRoutePattern(route string) RoutePattern {
	cleanPath(&route)
	route = strings.TrimPrefix(route, "/")
	parts := strings.Split(route, "/")
	segments := make([]RouteSegment, len(parts))
	for i, part := range parts {
		segments[i] = NewRouteSegment(part)
	}
	return RoutePattern{Segments: segments}
}

// Match checks whether the provided URL matches the pattern
//
// If matched, it returns a map of captured parameters (for dynamic or wildcard segments)
// but wildcard segments with an empty identifier are not captured
func (rp RoutePattern) Match(url string) (bool, map[string]string) {
	cleanPath(&url)
	url = strings.TrimPrefix(url, "/")
	urlParts := strings.Split(url, "/")
	params := make(map[string]string)

	i, j := 0, 0 // i for urlParts, j for rp.Segments
	for j < len(rp.Segments) {
		seg := rp.Segments[j]
		switch seg.Type {
		case StaticSegment:
			// if no URL part or it doesnt match then no match
			if i >= len(urlParts) || urlParts[i] != seg.Raw {
				return false, nil
			}
			i++
			j++
		case DynamicSegment:
			// a dynamic segment matches exactly one non-empty URL part
			if i >= len(urlParts) || urlParts[i] == "" {
				return false, nil
			}
			params[seg.Name] = urlParts[i]
			i++
			j++
		case WildcardSegment:
			// a wildcard can match zero or more segments
			// If it is the last segment, capture the remainder (if an identifier is provided)
			if j == len(rp.Segments)-1 {
				captured := strings.Join(urlParts[i:], "/")
				if captured != "" && captured[0] != '/' {
					captured = "/" + captured
				}
				if seg.Name != "" {
					params[seg.Name] = captured
				}
				i = len(urlParts)
				j++
				break
			}
			// ff not last capture until the next static segment is found
			nextSeg := rp.Segments[j+1]
			found := false
			var k int
			for k = i; k < len(urlParts); k++ {
				if nextSeg.Type == StaticSegment && urlParts[k] == nextSeg.Raw {
					captured := strings.Join(urlParts[i:k], "/")
					if captured != "" && captured[0] != '/' {
						captured = "/" + captured
					}
					if seg.Name != "" {
						params[seg.Name] = captured
					}
					i = k // move the url pointer to the boundary
					found = true
					break
				}
			}
			if !found {
				// capture the rest if no boundary is found
				captured := strings.Join(urlParts[i:], "/")
				if captured != "" && captured[0] != '/' {
					captured = "/" + captured
				}
				if seg.Name != "" {
					params[seg.Name] = captured
				}
				i = len(urlParts)
				j++
				break
			}
			j++
		}
	}
	// ensure all url parts were consumed
	if i != len(urlParts) {
		return false, nil
	}
	return true, params
}

// urlMatchesRoutePattern is a simple helper that returns whether urlPath matches routePath
func urlMatchesRoutePattern(urlPath string, routePath string) bool {
	rp := NewRoutePattern(routePath)
	matched, _ := rp.Match(urlPath)
	return matched
}

// constructElementsPattern extracts captured dynamic/wildcard parameters as a slice of PathValues
func constructElementsPattern(urlPath, routePath string) (pv []PathValues) {
	rp := NewRoutePattern(routePath)
	matched, params := rp.Match(urlPath)
	if !matched {
		return nil
	}
	for k, v := range params {
		pv = append(pv, PathValues{Key: k, Value: v})
	}
	return pv
}
