package goster

import "testing"

type CleanPathCase struct {
	name         string
	path         string
	expectedPath string
}

func TestCleanPath(t *testing.T) {
	testCases := []CleanPathCase{
		{
			name:         "Same path",
			path:         "home/var",
			expectedPath: "/home/var",
		},
		{
			name:         "No path",
			path:         "",
			expectedPath: "",
		},
		{
			name:         "Single slash",
			path:         "/",
			expectedPath: "",
		},
		{
			name:         "A lot of slashes",
			path:         "/////",
			expectedPath: "////",
		},
		{
			name:         "Trailing slash",
			path:         "home/var/",
			expectedPath: "/home/var",
		},
		{
			name:         "A lot of trailing slashes",
			path:         "home/var///////",
			expectedPath: "/home/var//////",
		},
	}

	failedCases := make(map[int]CleanPathCase, 0)
	for i, c := range testCases {
		cleanPath(&c.path)
		if c.path != c.expectedPath {
			failedCases[i] = c
		} else {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.name)
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}
