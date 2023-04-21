package handlers

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// cors - Stick this early in your handlers, sets server response headers in accoradance with cors.
func cors(v *viper.Viper, w http.ResponseWriter, r *http.Request) {
	for _, origin := range []string{
		"https://www.theneworleansseafoodconnection.com",
		"https://www.twilio.com",
	} {
		if origin == r.Header.Get("Origin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	if v.GetString("ENV") != "production" && (strings.Contains(r.Header.Get("Origin"), "localhost") || strings.Contains(r.Header.Get("Origin"), "127.0.0.1")) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, Referer")
	if r.Method == "OPTIONS" {
		return
	}
}
