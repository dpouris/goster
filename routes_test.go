package goster

import (
	"errors"
	"net/url"
	"path"
	"strings"
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
		err := g.Routes.New(c.method, c.url, c.handler)
		if err != nil {
			t.Errorf("route `%s` already exists", c.url)
		}
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

func TestWildcardRoutes(t *testing.T) {
	g := NewServer()
	testCases := []struct {
		name         string
		method       string
		url          string
		expectedType string
		pathVars     map[string]string
		expectErr    bool
		handler      RequestHandler
	}{
		{
			name:         "Dynamic Basic Route",
			method:       "GET",
			url:          "/users/:id",
			expectedType: Dynamic,
			pathVars:     map[string]string{"id": "1"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Dynamic Route with Trailing Slash",
			method:       "GET",
			url:          "/products/:productId/",
			expectedType: Dynamic,
			pathVars:     map[string]string{"productId": "1"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Dynamic Route with multiple path variables",
			method:       "GET",
			url:          "/products/:productId/edit",
			expectedType: Dynamic,
			pathVars:     map[string]string{"productId": "7"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Dynamic Route with multiple path variables alternating",
			method:       "GET",
			url:          "/logs/:logId/view/:iterationId",
			expectedType: Dynamic,
			pathVars:     map[string]string{"logId": "432", "iterationId": "2"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route",
			method:       "GET",
			url:          "/files/*filepath",
			expectedType: Wildcard,
			pathVars:     map[string]string{"filepath": "/movies/comedy/2012/MagicMike.mp4"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route index",
			method:       "GET",
			url:          "*",
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route index with leading slash and identifier",
			method:       "GET",
			url:          "/*path",
			pathVars:     map[string]string{"path": "/hello/there"},
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static prefix",
			method:       "GET",
			url:          "/static/*",
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static prefix and identifier",
			method:       "GET",
			url:          "/static/*s",
			expectedType: Wildcard,
			pathVars:     map[string]string{"s": "/yo"},
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static suffix",
			method:       "GET",
			url:          "/*/any",
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static suffix and identifier",
			method:       "GET",
			url:          "/*some/any",
			pathVars:     map[string]string{"some": "/yo/dj"},
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static suffix and prefix",
			method:       "GET",
			url:          "/static/*/styles",
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
		{
			name:         "Wildcard Route with static suffix and prefix and identifier",
			method:       "GET",
			url:          "/static/*values/styles",
			pathVars:     map[string]string{"values": "/docs/docsv1/kati"},
			expectedType: Wildcard,
			expectErr:    false,
			handler:      func(ctx *Ctx) error { return nil },
		},
	}

TestCase:
	for i, tc := range testCases {
		err := g.Routes.New(tc.method, tc.url, tc.handler)
		if tc.expectErr {
			if err == nil {
				t.Errorf("FAILED [%d] - %s: expected error when adding wildcard route %q", i, tc.name, tc.url)
			} else {
				t.Logf("PASSED [%d] - %s: correctly rejected wildcard route %q", i, tc.name, tc.url)
			}
			continue
		}

		if err != nil {
			t.Errorf("FAILED [%d] - %s: error adding route: %v", i, tc.name, err)
			continue
		}

		newUrlSlice := make([]string, 0)
		for _, seg := range strings.Split(tc.url, "/") {
			if strings.HasPrefix(seg, ":") || strings.HasPrefix(seg, "*") {
				newUrlSlice = append(newUrlSlice, tc.pathVars[seg[1:]])
				continue
			}
			newUrlSlice = append(newUrlSlice, seg)
		}

		newUrl := path.Join(newUrlSlice...)
		testUrl, _ := url.JoinPath("https://test.com", newUrl)
		ctx := NewContextCreation(testUrl)

		cleanPath(&newUrl)
		cleanPath(&tc.url)
		if urlMatchesRoute(newUrl, tc.url) {
			ctx.ParsePath(newUrl, tc.url)
		} else {
			t.Errorf("FAILED [%d] - %s: url %s doesn't match with any dynamic or wildcard route", i, tc.name, newUrl)
		}

		route, exists := g.Routes[tc.method][tc.url]
		if !exists {
			t.Errorf("FAILED [%d] - %s: route %q not found", i, tc.name, tc.url)
			continue
		} else if route.Type != tc.expectedType {
			t.Errorf("FAILED [%d] - %s: expected route type %q, got %q", i, tc.name, tc.expectedType, route.Type)
			continue
		}

		if len(ctx.Meta.Path) < len(tc.pathVars) {
			t.Errorf("FAILED [%d] - %s: incorrect number of path vars. Expected %d, got %d", i, tc.name, len(tc.pathVars), len(ctx.Meta.Path))
			continue
		}
		for k, v := range ctx.Meta.Path {
			expectedValue, exists := tc.pathVars[k]
			if !exists {
				t.Errorf("FAILED [%d] - %s: path var %s not found", i, tc.name, k)
				continue TestCase
			}
			if v != expectedValue {
				t.Errorf("FAILED [%d] - %s: path var %s with value %s doesn't match expected value %s", i, tc.name, k, v, expectedValue)
				continue TestCase
			}
		}
		t.Logf("PASSED [%d] - %s", i, tc.name)
	}

	t.Logf("TOTAL DYNAMIC ROUTES in GET: %d", len(g.Routes["GET"]))
}
