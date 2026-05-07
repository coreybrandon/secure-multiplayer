package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer_Fields(t *testing.T) {
	p := NewPlayer("test-id")
	assert.Equal(t, "test-id", p.ID)
	assert.GreaterOrEqual(t, p.X, 0.0)
	assert.Less(t, p.X, float64(CanvasWidth))
	assert.GreaterOrEqual(t, p.Y, 0.0)
	assert.Less(t, p.Y, float64(CanvasHeight))
	assert.Equal(t, 0, p.Score)
}

func TestMovePlayer_Up(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	p.MovePlayer("up", 5)
	assert.Equal(t, 95.0, p.Y)
	assert.Equal(t, 100.0, p.X)
}

func TestMovePlayer_Down(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	p.MovePlayer("down", 5)
	assert.Equal(t, 105.0, p.Y)
}

func TestMovePlayer_Left(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	p.MovePlayer("left", 5)
	assert.Equal(t, 95.0, p.X)
}

func TestMovePlayer_Right(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	p.MovePlayer("right", 5)
	assert.Equal(t, 105.0, p.X)
}

func TestMovePlayer_ClampsTop(t *testing.T) {
	p := &Player{X: 100, Y: 5}
	p.MovePlayer("up", 20)
	assert.Equal(t, 0.0, p.Y)
}

func TestMovePlayer_ClampsBottom(t *testing.T) {
	p := &Player{X: 100, Y: float64(CanvasHeight - PlayerSize - 5)}
	p.MovePlayer("down", 20)
	assert.Equal(t, float64(CanvasHeight-PlayerSize), p.Y)
}

func TestMovePlayer_ClampsLeft(t *testing.T) {
	p := &Player{X: 5, Y: 100}
	p.MovePlayer("left", 20)
	assert.Equal(t, 0.0, p.X)
}

func TestMovePlayer_ClampsRight(t *testing.T) {
	p := &Player{X: float64(CanvasWidth - PlayerSize - 5), Y: 100}
	p.MovePlayer("right", 20)
	assert.Equal(t, float64(CanvasWidth-PlayerSize), p.X)
}

func TestMovePlayer_ClampsMaxDist(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	p.MovePlayer("up", 100)
	assert.Equal(t, 80.0, p.Y) // capped at dist=20
}

func TestCollision_True(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	c := &Collectible{X: 100, Y: 100}
	assert.True(t, p.Collision(c))
}

func TestCollision_Overlapping(t *testing.T) {
	p := &Player{X: 100, Y: 100}
	c := &Collectible{X: 115, Y: 115}
	assert.True(t, p.Collision(c))
}

func TestCollision_False(t *testing.T) {
	p := &Player{X: 0, Y: 0}
	c := &Collectible{X: 500, Y: 500}
	assert.False(t, p.Collision(c))
}

func TestCalculateRank_First(t *testing.T) {
	p1 := &Player{ID: "a", Score: 100}
	p2 := &Player{ID: "b", Score: 50}
	p3 := &Player{ID: "c", Score: 25}
	assert.Equal(t, "Rank: 1 / 3", p1.CalculateRank([]*Player{p1, p2, p3}))
}

func TestCalculateRank_Last(t *testing.T) {
	p1 := &Player{ID: "a", Score: 100}
	p2 := &Player{ID: "b", Score: 50}
	p3 := &Player{ID: "c", Score: 25}
	assert.Equal(t, "Rank: 3 / 3", p3.CalculateRank([]*Player{p1, p2, p3}))
}

func TestCalculateRank_Solo(t *testing.T) {
	p := &Player{ID: "a", Score: 0}
	assert.Equal(t, "Rank: 1 / 1", p.CalculateRank([]*Player{p}))
}
