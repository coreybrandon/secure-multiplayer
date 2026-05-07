package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Collectible struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Value int     `json:"value"`
}

func NewCollectible() *Collectible {
	return &Collectible{
		ID:    fmt.Sprintf("c-%d", time.Now().UnixNano()),
		X:     rand.Float64() * (CanvasWidth - PlayerSize),
		Y:     rand.Float64() * (CanvasHeight - PlayerSize),
		Value: rand.Intn(10) + 1,
	}
}
