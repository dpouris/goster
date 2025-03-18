package goster

import (
	"net/http/httptest"
	"strings"
)

func NewContextCreation(url string) Ctx {
	req := httptest.NewRequest("GET", "http://example.com/test", strings.NewReader("body"))
	rec := httptest.NewRecorder()

	return NewContext(req, rec)
}
