package handlers

import (
	"net/http"
)

// cors - Stick this early in your handlers, sets server response headers in accoradance with cors.
func cors(w http.ResponseWriter, r *http.Request) {
	for _, origin := range []string{
		"https://www.theneworleansseafoodconnection.com",
		"https://www.twilio.com",
	} {
		if origin == r.Header.Get("Origin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, Referer")
	if r.Method == "OPTIONS" {
		return
	}
}
