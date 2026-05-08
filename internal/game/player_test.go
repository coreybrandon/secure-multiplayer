package game

import "testing"

func TestNewPlayer_Fields(t *testing.T) {
	p := NewPlayer("test-id")
	if p.ID != "test-id" {
		t.Errorf("ID = %q, want %q", p.ID, "test-id")
	}
	if p.X < 0 || p.X >= CanvasWidth {
		t.Errorf("X = %v out of canvas bounds", p.X)
	}
	if p.Y < 0 || p.Y >= CanvasHeight {
		t.Errorf("Y = %v out of canvas bounds", p.Y)
	}
	if p.Score != 0 {
		t.Errorf("Score = %d, want 0", p.Score)
	}
}

func TestMovePlayer(t *testing.T) {
	tests := []struct {
		name   string
		startX float64
		startY float64
		dir    string
		dist   float64
		wantX  float64
		wantY  float64
	}{
		{"up", 100, 100, "up", 5, 100, 95},
		{"down", 100, 100, "down", 5, 100, 105},
		{"left", 100, 100, "left", 5, 95, 100},
		{"right", 100, 100, "right", 5, 105, 100},
		{"clamps top", 100, 5, "up", 20, 100, 0},
		{"clamps bottom", 100, CanvasHeight - PlayerSize - 5, "down", 20, 100, CanvasHeight - PlayerSize},
		{"clamps left", 5, 100, "left", 20, 0, 100},
		{"clamps right", CanvasWidth - PlayerSize - 5, 100, "right", 20, CanvasWidth - PlayerSize, 100},
		{"clamps max dist", 100, 100, "up", 100, 100, 80},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &Player{X: tc.startX, Y: tc.startY}
			p.MovePlayer(tc.dir, tc.dist)
			if p.X != tc.wantX {
				t.Errorf("X = %v, want %v", p.X, tc.wantX)
			}
			if p.Y != tc.wantY {
				t.Errorf("Y = %v, want %v", p.Y, tc.wantY)
			}
		})
	}
}

func TestCollision(t *testing.T) {
	tests := []struct {
		name    string
		player  Player
		collect Collectible
		want    bool
	}{
		{"exact overlap", Player{X: 100, Y: 100}, Collectible{X: 100, Y: 100}, true},
		{"partial overlap", Player{X: 100, Y: 100}, Collectible{X: 115, Y: 115}, true},
		{"no overlap", Player{X: 0, Y: 0}, Collectible{X: 500, Y: 500}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.player.Collision(&tc.collect)
			if got != tc.want {
				t.Errorf("Collision() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRank(t *testing.T) {
	p1 := &Player{ID: "a", Score: 100}
	p2 := &Player{ID: "b", Score: 50}
	p3 := &Player{ID: "c", Score: 25}
	solo := &Player{ID: "solo", Score: 0}

	tests := []struct {
		name      string
		player    *Player
		players   []*Player
		wantRank  int
		wantTotal int
	}{
		{"first place", p1, []*Player{p1, p2, p3}, 1, 3},
		{"last place", p3, []*Player{p1, p2, p3}, 3, 3},
		{"solo", solo, []*Player{solo}, 1, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rank, total := tc.player.Rank(tc.players)
			if rank != tc.wantRank {
				t.Errorf("rank = %d, want %d", rank, tc.wantRank)
			}
			if total != tc.wantTotal {
				t.Errorf("total = %d, want %d", total, tc.wantTotal)
			}
		})
	}
}
