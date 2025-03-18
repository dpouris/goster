package goster

import (
	"testing"
)

// equalMaps is a helper function to compare two maps.
func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func TestRoutePattern_Match(t *testing.T) {
	tests := []struct {
		route  string
		url    string
		match  bool
		params map[string]string
	}{
		{
			route:  "/users/:id",
			url:    "/users/123",
			match:  true,
			params: map[string]string{"id": "123"},
		},
		{
			route:  "/users/:id",
			url:    "/users/",
			match:  false,
			params: nil,
		},
		{
			route:  "/users/:id/profile",
			url:    "/users/42/profile",
			match:  true,
			params: map[string]string{"id": "42"},
		},
		{
			route:  "/files/*filepath",
			url:    "/files/movies/comedy/2012/MagicMike.mp4",
			match:  true,
			params: map[string]string{"filepath": "/movies/comedy/2012/MagicMike.mp4"},
		},
		{
			route:  "/files/*filepath",
			url:    "/files",
			match:  true,
			params: map[string]string{"filepath": ""},
		},
		{
			// Wildcard with no identifier: no parameter should be captured.
			route:  "/*",
			url:    "/hello",
			match:  true,
			params: map[string]string{},
		},
		{
			// Mix of dynamic and wildcard with identifier.
			route:  "/:a/*b",
			url:    "/foo/bar/baz",
			match:  true,
			params: map[string]string{"a": "foo", "b": "/bar/baz"},
		},
		{
			// Route with a trailing static segment after wildcard.
			route:  "/static/*filepath",
			url:    "/static",
			match:  true,
			params: map[string]string{"filepath": ""},
		},
		{
			route:  "/static/*filepath",
			url:    "/static/",
			match:  true,
			params: map[string]string{"filepath": ""},
		},
		{
			// Root route.
			route:  "/",
			url:    "/",
			match:  true,
			params: map[string]string{},
		},
		{
			route:  "/:id",
			url:    "/123",
			match:  true,
			params: map[string]string{"id": "123"},
		},
		{
			route:  "/:id",
			url:    "/",
			match:  false,
			params: nil,
		},
	}

	for _, tt := range tests {
		rp := NewRoutePattern(tt.route)
		gotMatch, gotParams := rp.Match(tt.url)
		if gotMatch != tt.match {
			t.Errorf("RoutePattern(%q).Match(%q) expected match %v, got %v", tt.route, tt.url, tt.match, gotMatch)
		}
		if tt.match && !equalMaps(gotParams, tt.params) {
			t.Errorf("RoutePattern(%q).Match(%q) expected params %v, got %v", tt.route, tt.url, tt.params, gotParams)
		}
	}
}

func TestUrlMatchesRoute(t *testing.T) {
	tests := []struct {
		route string
		url   string
		match bool
	}{
		{route: "/users/:id", url: "/users/123", match: true},
		{route: "/users/:id", url: "/users/", match: false},
		{route: "/files/*filepath", url: "/files/movies/comedy/2012/MagicMike.mp4", match: true},
		{route: "/static", url: "/static", match: true},
		{route: "/static", url: "/static/", match: true},
	}
	for _, tt := range tests {
		got := urlMatchesRoute(tt.url, tt.route)
		if got != tt.match {
			t.Errorf("urlMatchesRoute(%q, %q) expected %v, got %v", tt.url, tt.route, tt.match, got)
		}
	}
}

func TestConstructElements(t *testing.T) {
	tests := []struct {
		route  string
		url    string
		params map[string]string
	}{
		{
			route:  "/users/:id",
			url:    "/users/123",
			params: map[string]string{"id": "123"},
		},
		{
			route:  "/files/*filepath",
			url:    "/files/movies/comedy/2012/MagicMike.mp4",
			params: map[string]string{"filepath": "/movies/comedy/2012/MagicMike.mp4"},
		},
		{
			route:  "/*",
			url:    "/hello",
			params: map[string]string{},
		},
		{
			route:  "/:a/*b",
			url:    "/foo/bar/baz",
			params: map[string]string{"a": "foo", "b": "/bar/baz"},
		},
	}

	for _, tt := range tests {
		pv := constructElements(tt.url, tt.route)
		gotParams := make(map[string]string)
		for _, v := range pv {
			gotParams[v.Key] = v.Value
		}
		if !equalMaps(gotParams, tt.params) {
			t.Errorf("constructElements(%q, %q) expected params %v, got %v", tt.url, tt.route, tt.params, gotParams)
		}
	}
}
