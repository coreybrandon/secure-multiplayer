package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/coreybrandon/secure-multi-game/internal/game"
	"github.com/gorilla/websocket"
)

const sendBufSize = 256

var nextID atomic.Uint64

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   string
}

type moveMsg struct {
	client *Client
	dir    string
	dist   float64
}

type InboundMsg struct {
	Type string  `json:"type"`
	Dir  string  `json:"dir,omitempty"`
	Dist float64 `json:"dist,omitempty"`
}

type OutboundMsg struct {
	Type        string            `json:"type"`
	Player      *game.Player      `json:"player,omitempty"`
	Players     []*game.Player    `json:"players,omitempty"`
	Collectible *game.Collectible `json:"collectible,omitempty"`
	ID          string            `json:"id,omitempty"`
}

type Hub struct {
	clients     map[string]*Client
	players     map[string]*game.Player
	collectible *game.Collectible
	register    chan *Client
	unregister  chan *Client
	move        chan moveMsg
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		players:     make(map[string]*game.Player),
		collectible: game.NewCollectible(),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		move:        make(chan moveMsg),
	}
}

// Register upgrades a connection and joins the game. Safe to call from any goroutine.
func (h *Hub) Register(conn *websocket.Conn) {
	id := fmt.Sprintf("p-%d", nextID.Add(1))
	c := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, sendBufSize),
		id:   id,
	}
	h.register <- c
	go c.writePump()
	go c.readPump()
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c.id] = c
			h.players[c.id] = game.NewPlayer(c.id)
			players := h.playerList()
			h.sendTo(c, OutboundMsg{
				Type:        "init",
				Player:      h.players[c.id],
				Players:     players,
				Collectible: h.collectible,
			})
			h.broadcastExcept(c.id, OutboundMsg{
				Type:        "state",
				Players:     players,
				Collectible: h.collectible,
			})

		case c := <-h.unregister:
			if _, ok := h.clients[c.id]; ok {
				delete(h.clients, c.id)
				delete(h.players, c.id)
				close(c.send)
			}
			h.broadcastAll(OutboundMsg{
				Type: "player-left",
				ID:   c.id,
			})

		case msg := <-h.move:
			p, ok := h.players[msg.client.id]
			if !ok {
				continue
			}
			p.MovePlayer(msg.dir, msg.dist)
			if p.Collision(h.collectible) {
				p.Score += h.collectible.Value
				h.collectible = game.NewCollectible()
			}
			h.broadcastAll(OutboundMsg{
				Type:        "state",
				Players:     h.playerList(),
				Collectible: h.collectible,
			})
		}
	}
}

func (h *Hub) playerList() []*game.Player {
	list := make([]*game.Player, 0, len(h.players))
	for _, p := range h.players {
		list = append(list, p)
	}
	return list
}

func (h *Hub) sendTo(c *Client, msg OutboundMsg) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (h *Hub) broadcastAll(msg OutboundMsg) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	for _, c := range h.clients {
		select {
		case c.send <- data:
		default:
		}
	}
}

func (h *Hub) broadcastExcept(excludeID string, msg OutboundMsg) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	for id, c := range h.clients {
		if id == excludeID {
			continue
		}
		select {
		case c.send <- data:
		default:
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg InboundMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		if msg.Type == "move" {
			c.hub.move <- moveMsg{client: c, dir: msg.Dir, dist: msg.Dist}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
