package wast

// MaxScore is the maximum score one player can obtain during one game
var MaxScore = 40

// damageScore returns a score depending on the given damage number
func damageScore(damage int) int {
	switch {
	case damage < 500:
		return 0
	case damage < 800:
		return 1
	case damage < 1000:
		return 2
	case damage < 1200:
		return 4
	case damage < 1400:
		return 6
	case damage < 1600:
		return 8
	case damage < 1800:
		return 10
	case damage < 2000:
		return 12
	case damage < 2200:
		return 14
	case damage < 2400:
		return 16
	case damage < 2600:
		return 18
	default:
		return 20
	}
}

// killingScore returns a score depending on the given number of kills
func killingScore(kills int) int {
	switch kills {
	case 0:
		return 0
	case 1:
		return 2
	case 2:
		return 4
	case 3:
		return 6
	case 4:
		return 7
	case 5:
		return 8
	case 6:
		return 10
	case 7, 8:
		return 11
	case 9:
		return 12
	case 10:
		return 13
	case 11:
		return 14
	default:
		return 15
	}
}

// victoryScore returns a score depending on the given bool
func victoryScore(winner bool) int {
	if winner {
		return 5
	}

	return 0
}
