package handler

import (
	"net/http"
)

const (
	indexHTML = "views/index.html"
	publicDir = "public"
)

func NewStaticHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, indexHTML)
	})

	fs := http.FileServer(http.Dir(publicDir))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	return mux
}
