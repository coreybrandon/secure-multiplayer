package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	hdr := rec.Header()
	assert.Equal(t, "nosniff", hdr.Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", hdr.Get("X-XSS-Protection"))
	assert.Equal(t, "no-store", hdr.Get("Surrogate-Control"))
	assert.Equal(t, "no-store, no-cache, must-revalidate, proxy-revalidate", hdr.Get("Cache-Control"))
	assert.Equal(t, "no-cache", hdr.Get("Pragma"))
	assert.Equal(t, "0", hdr.Get("Expires"))
	assert.Equal(t, "PHP 7.4.3", hdr.Get("X-Powered-By"))
}
