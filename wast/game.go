package wast

import "time"

// Game represents a Worms game
type Game struct {
	ID        uint
	CreatedAt time.Time
	Stats     []Stats
}
