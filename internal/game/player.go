// Package game contains the core game entities and rules.
package game

import (
	"math"
	"math/rand"
	"sort"
)

const (
	CanvasWidth  = 640
	CanvasHeight = 480
	PlayerSize   = 30
	MaxMoveDist  = 20
)

const (
	dirUp    = "up"
	dirDown  = "down"
	dirLeft  = "left"
	dirRight = "right"
)

type Player struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
}

func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
		X:  rand.Float64() * (CanvasWidth - PlayerSize),
		Y:  rand.Float64() * (CanvasHeight - PlayerSize),
	}
}

func (p *Player) MovePlayer(dir string, dist float64) {
	if dist > MaxMoveDist {
		dist = MaxMoveDist
	}
	switch dir {
	case dirUp:
		p.Y = math.Max(0, p.Y-dist)
	case dirDown:
		p.Y = math.Min(CanvasHeight-PlayerSize, p.Y+dist)
	case dirLeft:
		p.X = math.Max(0, p.X-dist)
	case dirRight:
		p.X = math.Min(CanvasWidth-PlayerSize, p.X+dist)
	}
}

func (p *Player) Collision(c *Collectible) bool {
	return math.Abs(p.X-c.X) < PlayerSize && math.Abs(p.Y-c.Y) < PlayerSize
}

func (p *Player) Rank(players []*Player) (rank, total int) {
	sorted := make([]*Player, len(players))
	copy(sorted, players)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})
	rank = len(players)
	for i, pl := range sorted {
		if pl.ID == p.ID {
			rank = i + 1
			break
		}
	}
	return rank, len(players)
}
