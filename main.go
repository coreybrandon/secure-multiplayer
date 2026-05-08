package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/coreybrandon/secure-multi-game/internal/handler"
	"github.com/coreybrandon/secure-multi-game/internal/hub"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	h := hub.NewHub()
	go h.Run(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", handler.NewStaticHandler())
	mux.Handle("/ws", handler.NewWSHandler(h))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler.SecurityHeaders(mux)))
}
