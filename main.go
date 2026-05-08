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

	srv := &http.Server{Addr: ":" + port, Handler: handler.SecurityHeaders(mux)}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	log.Printf("Listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
