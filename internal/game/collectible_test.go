package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCollectible_Fields(t *testing.T) {
	c := NewCollectible()
	assert.NotEmpty(t, c.ID)
	assert.GreaterOrEqual(t, c.X, 0.0)
	assert.Less(t, c.X, float64(CanvasWidth))
	assert.GreaterOrEqual(t, c.Y, 0.0)
	assert.Less(t, c.Y, float64(CanvasHeight))
	assert.Greater(t, c.Value, 0)
	assert.LessOrEqual(t, c.Value, 10)
}
