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
