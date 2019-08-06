package wast

// Player represents a Worms player
type Player struct {
	ID         uint
	Nickname   string
	DiscordTag string
	Stats      []Stats
}
