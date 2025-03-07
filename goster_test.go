package goster

import (
	"fmt"
	"testing"
)

type IsDynamicMatchCase struct {
	name           string
	url            string
	dynamicPath    string
	expectedResult bool
}

func TestIsDynRouteMatch(t *testing.T) {
	g := NewServer()

	testCases := []IsDynamicMatchCase{
		{
			name:           "Depth 1",
			dynamicPath:    "path/another_path/:var",
			url:            "path/another_path/something",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with '/' suffix on URL)",
			dynamicPath:    "path/another_path/:var",
			url:            "path/another_path/something/",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with '/' suffix on dynPath)",
			dynamicPath:    "path/another_path/:var/",
			url:            "path/another_path/something",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with Depth 2 URL)",
			dynamicPath:    "path/another_path/:var",
			url:            "path/another_path/something/something2",
			expectedResult: false,
		},
		{
			name:           "Depth 2",
			dynamicPath:    "path/another_path/:var/:var2",
			url:            "path/another_path/something/something2",
			expectedResult: true,
		},
		{
			name:           "Depth 2 (with Depth 1 URL)",
			dynamicPath:    "path/another_path/:var/:var2",
			url:            "path/another_path/something",
			expectedResult: false,
		},
		{
			name:           "Depth 2 (with '/' suffix on dynPath)",
			dynamicPath:    "path/another_path/:var/:var2/",
			url:            "path/another_path/something/something2",
			expectedResult: true,
		},
		{
			name:           "Depth 2 (with '/' suffix on URL)",
			dynamicPath:    "path/another_path/:var/:var2",
			url:            "path/another_path/something/something2/",
			expectedResult: true,
		},
	}

	failedCases := make(map[int]IsDynamicMatchCase, 0)
	for i, c := range testCases {
		if g.isDynamicRouteMatch(c.url, c.dynamicPath) != c.expectedResult {
			failedCases[i] = c
		} else {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.name)
		t.Errorf("Expected %t for '%s' and '%s'", c.expectedResult, c.url, c.dynamicPath)
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}

type TemplateDirMatch struct {
	name        string
	givenPath   string
	exectedPath string
}

func TestTemplateDir(t *testing.T) {
	g := NewServer()

	testCases := []TemplateDirMatch{
		{
			name:        "Test 1",
			givenPath:   "/templates",
			exectedPath: "",
		},
		// {
		// 	name:        "Test 2",
		// 	givenPath:   "/static/templates/",
		// 	exectedPath: "",
		// },
	}

	failedCases := make(map[int]TemplateDirMatch, 0)
	for _, tmpl := range testCases {
		err := g.TemplateDir(tmpl.givenPath)
		if err != nil {
			t.Error(err)
		}
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}

func TestStaticDir(t *testing.T) {
	g := NewServer()

	testCases := []TemplateDirMatch{
		{
			name:        "Test 1",
			givenPath:   "/static",
			exectedPath: "",
		},
	}

	failedCases := make(map[int]TemplateDirMatch, 0)
	for _, tmpl := range testCases {
		err := g.StaticDir(tmpl.givenPath)
		if err != nil {
			t.Error("could not set templates dir")
		}

		for route := range g.Routes["GET"] {
			fmt.Printf("Route %s\n", route)
		}
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}
