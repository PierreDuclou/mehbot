package wast

// Player represents a Worms player
type Player struct {
	ID       string
	Nickname string
	Stats    []Stats
}

// NewPlayer returns a pointer to a new Player instance
func NewPlayer(ID string, nickename string) *Player {
	return &Player{ID: ID, Nickname: nickename}
}
