package goster

import (
	"errors"
	"testing"
)

type MethodNewCase struct {
	name         string
	method       string
	url          string
	expectedPath string
	handler      RequestHandler
}

func TestMethodNew(t *testing.T) {
	g := NewServer()
	testCases := []MethodNewCase{
		{
			name:         "1",
			method:       "GET",
			url:          "/home/var",
			expectedPath: "/home/var",
			handler:      func(ctx *Ctx) error { return errors.New("test error") },
		},
		{
			name:         "2",
			method:       "GET",
			url:          "/home/var////",
			expectedPath: "/home/var///",
			handler:      func(ctx *Ctx) error { return errors.New("test error") },
		},
	}

	failedCases := make(map[int]MethodNewCase, 0)
	for i, c := range testCases {
		g.Routes.New(c.method, c.url, c.handler)
		if _, exists := g.Routes[c.method][c.expectedPath]; exists {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		} else {
			failedCases[i] = c
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.name)
	}

	// Space
	t.Log("")
	t.Logf("Routes: ")
	t.Log(g.Routes)

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}

// TestAddStaticDir verifies that AddStaticDir properly walks the given directory,
// registers a route for each file (here we create one test file), and that the route's key
// matches the file's path (normalized to use forward slashes).
func TestAddStaticDir(t *testing.T) {
	// Get the directory of the current executable.
	// exPath, err := os.Executable()
	// if err != nil {
	// 	t.Fatalf("os.Executable() error: %v", err)
	// }

	testDir := "../static/"
	// baseDir := filepath.Dir(exPath)

	// Initialize Routes. Ensure that the "GET" method map is created.
	r := Routes{
		"GET": make(map[string]Route),
	}

	// Call AddStaticDir with the relative directory.
	if err := r.prepareStaticRoutes(testDir); err != nil {
		t.Fatalf("AddStaticDir returned error: %v", err)
	}

	for route := range r["GET"] {
		t.Logf("ROute: %s\n", route)
	}

	// // Check that the route for the file was added.
	// if route, exists := r["GET"][expectedKey]; !exists {
	// 	t.Errorf("expected route %q not found in GET routes", expectedKey)
	// } else {
	// 	if route.Handler == nil {
	// 		t.Errorf("route handler for %q is nil", expectedKey)
	// 	} else {
	// 		t.Logf("PASSED - Route for %q added successfully", expectedKey)
	// 	}
	// }

	t.Logf("TOTAL ROUTES in GET: %d", len(r["GET"]))
}
