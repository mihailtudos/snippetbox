package main

import (
	"log/slog"
	"net/http"
)

// secureHeaders middleware adds security related headers such as: CSP, Referrer-policy, etc.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// logRequest middleware would log the retails of each incoming request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.RequestURI
		)

		app.logger.Info("received request",
			slog.String("ip", ip),
			slog.String("proto", proto),
			slog.String("method", method),
			slog.String("uri", uri),
		)

		next.ServeHTTP(w, r)
	})
}
