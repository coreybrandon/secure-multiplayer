// Package handler contains HTTP handlers and middleware.
package handler

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		h.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; frame-ancestors 'none'")
		h.Set("Surrogate-Control", "no-store")
		h.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		h.Set("Pragma", "no-cache")
		h.Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
