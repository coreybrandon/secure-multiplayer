package game

import (
	"fmt"
	"math/rand"
	"sync/atomic"
)

const (
	minCollectibleValue = 1
	maxCollectibleValue = 10
)

var collectibleID atomic.Uint64

type Collectible struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Value int     `json:"value"`
}

func NewCollectible() *Collectible {
	return &Collectible{
		ID:    fmt.Sprintf("c-%d", collectibleID.Add(1)),
		X:     rand.Float64() * (CanvasWidth - PlayerSize),
		Y:     rand.Float64() * (CanvasHeight - PlayerSize),
		Value: rand.Intn(maxCollectibleValue-minCollectibleValue+1) + minCollectibleValue,
	}
}
