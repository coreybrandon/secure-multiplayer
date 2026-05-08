package game

import "testing"

func TestNewCollectible_Fields(t *testing.T) {
	c := NewCollectible()
	if c.ID == "" {
		t.Error("ID is empty")
	}
	if c.X < 0 || c.X >= CanvasWidth {
		t.Errorf("X = %v out of canvas bounds", c.X)
	}
	if c.Y < 0 || c.Y >= CanvasHeight {
		t.Errorf("Y = %v out of canvas bounds", c.Y)
	}
	if c.Value < minCollectibleValue || c.Value > maxCollectibleValue {
		t.Errorf("Value = %d, want between %d and %d", c.Value, minCollectibleValue, maxCollectibleValue)
	}
}
