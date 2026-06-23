package gameplay

import (
	"math/rand/v2"
)


func diceRoll() (diceVal int) {
	diceVal = rand.IntN(6) + 1 + rand.IntN(6) + 1
	return
}

func updatePos(diceVal, currPos int) (newPos int) {
	movingTo := currPos + diceVal
	if movingTo >= 40 {
		newPos = movingTo - 40
	} else {
	newPos = movingTo
	}

	return
}

func ownershipConfirm(owenerId, playerId string) (bool) {
	if owenerId == playerId {
		return true
	} else {
		return  false
	}
}
