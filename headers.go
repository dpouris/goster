package gottp_server

import "net/http"

func DefaultHeader(h *http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Connection", "Keep-Alive")
	h.Set("Keep-Alive", "timeout=5, max=997")
}
