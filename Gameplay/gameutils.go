package gameplay

import (
	models "Monopoly/Models"
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

func nextSeq(currSeq, count int) (nextSeq int) {
	nextSeq = currSeq + 1
	if nextSeq == count {
		nextSeq = 0
	}
	return
}

func utilityRentCalculation(baseRent, cardCount int) (rent int) {

	if cardCount < 2 {
		rent = baseRent

	} else if cardCount < 4 {
		rent = baseRent * 2

	} else  {
		rent = baseRent * 4

	}

	return
} 

func cityRentCalculation(baseRent int, status string) (rent int) {

	switch status {

	case models.BASE:
		rent = baseRent

	case models.COLOUR:
		rent = baseRent * 2

	case models.HOUSE_1:
		rent = baseRent * 2

	case models.HOUSE_2:
		rent = baseRent * 15

	case models.HOUSE_3:
		rent = baseRent * 35

	case models.HOUSE_4:
		rent = baseRent * 45

	case models.HOTEL:
		rent = baseRent * 50

	}

	return
}