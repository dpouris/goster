package goster

import (
	"testing"
)

type IsDynamicMatchCase struct {
	name           string
	url            string
	dynPath        string
	expectedResult bool
}

func TestIsDynRouteMatch(t *testing.T) {
	g := NewServer()

	testCases := []IsDynamicMatchCase{
		{
			name:           "Depth 1",
			dynPath:        "path/another_path/:var",
			url:            "path/another_path/something",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with '/' suffix on URL)",
			dynPath:        "path/another_path/:var",
			url:            "path/another_path/something/",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with '/' suffix on dynPath)",
			dynPath:        "path/another_path/:var/",
			url:            "path/another_path/something",
			expectedResult: true,
		},
		{
			name:           "Depth 1 (with Depth 2 URL)",
			dynPath:        "path/another_path/:var",
			url:            "path/another_path/something/something2",
			expectedResult: false,
		},
		{
			name:           "Depth 2",
			dynPath:        "path/another_path/:var/:var2",
			url:            "path/another_path/something/something2",
			expectedResult: true,
		},
		{
			name:           "Depth 2 (with Depth 1 URL)",
			dynPath:        "path/another_path/:var/:var2",
			url:            "path/another_path/something",
			expectedResult: false,
		},
		{
			name:           "Depth 2 (with '/' suffix on dynPath)",
			dynPath:        "path/another_path/:var/:var2/",
			url:            "path/another_path/something/something2",
			expectedResult: true,
		},
		{
			name:           "Depth 2 (with '/' suffix on URL)",
			dynPath:        "path/another_path/:var/:var2",
			url:            "path/another_path/something/something2/",
			expectedResult: true,
		},
	}

	failedCases := make(map[int]IsDynamicMatchCase, 0)
	for i, c := range testCases {
		if g.isDynRouteMatch(c.url, c.dynPath) != c.expectedResult {
			failedCases[i] = c
		} else {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.name)
		t.Errorf("Expected %t for '%s' and '%s'", c.expectedResult, c.url, c.dynPath)
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}
