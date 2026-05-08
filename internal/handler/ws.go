package handler

import (
	"log"
	"net/http"

	"github.com/coreybrandon/secure-multi-game/internal/hub"
	"github.com/gorilla/websocket"
)

// upgrader accepts connections from any origin; this is intentional for a toy game with no auth.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWSHandler(h *hub.Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}
		h.Register(conn)
	})
}
