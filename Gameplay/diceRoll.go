package gameplay

import "math/rand/v2"


func DiceRoll() (diceVal int) {
	diceVal = rand.IntN(6) + 1 + rand.IntN(6) + 1
	return
}