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

// NewStats returns a pointer to a new Stats instance
func NewStats(kills int, deaths int, damage int, winner bool, playerID string, gameID uint) *Stats {
	return &Stats{
		Kills:    kills,
		Deaths:   deaths,
		Damage:   damage,
		Winner:   winner,
		PlayerID: playerID,
		GameID:   gameID,
	}
}

// score returns the final player score for a certain game
func (s Stats) score() int {
	return (damageScore(s.Damage) + killingScore(s.Kills) + victoryScore(s.Winner)) - MaxScore/2
}
