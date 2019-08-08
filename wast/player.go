package wast

// Player represents a Worms player
type Player struct {
	ID         uint
	Nickname   string
	DiscordTag string
	Stats      []Stats
}

// NewPlayer returns a pointer to a new Player instance
func NewPlayer(nickename string, discordTag string) *Player {
	return &Player{Nickname: nickename, DiscordTag: discordTag}
}
