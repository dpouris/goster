package goster

import (
	"maps"
	"testing"
)

type ParseUrlCase struct {
	name                string
	url                 string
	expectedQueryParams map[string]string
	shouldFail          bool
}

func TestParseUrl(t *testing.T) {
	testCases := []ParseUrlCase{
		{
			name: "1",
			url:  "/var/home?name=dimitris&age=24",
			expectedQueryParams: map[string]string{
				"name": "dimitris",
				"age":  "24",
			},
			shouldFail: false,
		},
		{
			name: "2",
			url:  "/var/home?name=dimitris&age=23",
			expectedQueryParams: map[string]string{
				"name": "dimitris",
				"age":  "24",
			},
			shouldFail: true,
		},
		{
			name: "3",
			url:  "/var/home?name=dimitris&name=gearge&age=24",
			expectedQueryParams: map[string]string{
				"name": "dimitris",
				"age":  "24",
			},
			shouldFail: true,
		},
	}

	failedCases := make(map[int]struct {
		Meta
		ParseUrlCase
	}, 0)
	for i, c := range testCases {
		meta := Meta{
			Query: make(map[string]string),
		}
		meta.ParseQueryParams(c.url)
		if (!maps.Equal(meta.Query, c.expectedQueryParams)) == !c.shouldFail {
			failedCases[i] = struct {
				Meta
				ParseUrlCase
			}{meta, c}
		} else {
			t.Logf("PASSED [%d] - %s\n", i, c.name)
		}
	}

	// Space
	t.Log("")

	for i, c := range failedCases {
		t.Errorf("FAILED [%d] - %s\n", i, c.ParseUrlCase.name)
		t.Errorf("Expected '%v' path, but got '%v'", c.expectedQueryParams, c.Query)
	}

	t.Logf("TOTAL CASES: %d\n", len(testCases))
	t.Logf("FAILED CASES: %d\n", len(failedCases))
}
