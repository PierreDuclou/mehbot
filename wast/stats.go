package wast

// Stats represents the stats of a player for a certain game
type Stats struct {
	ID       uint
	Kills    int
	Deaths   int
	Damage   int
	Winner   bool
	Player   Player
	Game     Game
	PlayerID string
	GameID   uint
}

// score returns the final player score for a certain game
func (s Stats) score() int {
	return (damageScore(s.Damage) + killingScore(s.Kills) + victoryScore(s.Winner)) - MaxScore/2
}
