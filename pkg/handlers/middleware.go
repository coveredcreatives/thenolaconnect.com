package handlers

import (
	"net/http"
)

// cors - Stick this early in your handlers, sets server response headers in accoradance with cors.
// Not very secure at the moment.
func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, Referer")
	if r.Method == "OPTIONS" {
		return
	}
}
