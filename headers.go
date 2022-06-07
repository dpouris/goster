package gottp_server

import "net/http"

func MakeSuccessHeader(h *http.Header) {
	h.Add("Access-Control-Allow-Origin", "*")
	// h.Add("Content-Type", "application/json")
	h.Add("Connection", "Keep-Alive")
	h.Add("Keep-Alive", "timeout=5, max=997")
}
