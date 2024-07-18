package goster

import (
	"testing"
)

type ParseUrlCase struct {
	name        string
	url         string
	expectedUrl string
	shouldFail  bool
}

func TestParseUrl(t *testing.T) {
	meta := Meta{
		Query: Params{
			values: make(map[string]string),
		},
	}
	testCases := []ParseUrlCase{
		{
			name:        "1",
			url:         "/home/var/",
			expectedUrl: "/home/var",
			shouldFail:  false,
		},
		{
			name:        "2",
			url:         "/home/var/",
			expectedUrl: "/home/var/",
			shouldFail:  true,
		},
		{
			name:        "3",
			url:         "/home/var///",
			expectedUrl: "/home/var///",
			shouldFail:  true,
		},
		{
			name:        "4",
			url:         "/home/var///",
			expectedUrl: "/home/var//",
			shouldFail:  false,
		},
		{
			name:        "5",
			url:         "////",
			expectedUrl: "///",
			shouldFail:  false,
		},
	}

	failedCases := make(map[int]ParseUrlCase, 0)
	for i, c := range testCases {
		meta.ParseUrl(&c.url)
		if (c.url != c.expectedUrl) == !c.shouldFail {
			failedCases[i] = c
		} else {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.name)
		t.Errorf("Expected '%s' path, but got '%s'", c.expectedUrl, c.url)
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}
