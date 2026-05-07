package handler

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Surrogate-Control", "no-store")
		h.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		h.Set("Pragma", "no-cache")
		h.Set("Expires", "0")
		h.Set("X-Powered-By", "PHP 7.4.3")
		next.ServeHTTP(w, r)
	})
}
