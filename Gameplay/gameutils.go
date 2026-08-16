package gameplay

import (
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
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

func houseHotelFine(statusList []string) (fine int) {
	fine = 0

	for _, status := range statusList {
		
		switch status {

		case models.HOUSE_1:
			fine += 400 * 1

		case models.HOUSE_2:
			fine += 400 * 2

		case models.HOUSE_3:
			fine += 400 * 3

		case models.HOUSE_4:
			fine += 400 * 4

		case models.HOTEL:
			fine += 1150

		default:
			fine += 0
			
		}
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

func callChangePlayer(client *models.Client, readCh chan<- models.WSMessage) {

	changePlayerReq := models.Request{}
	payload, err := json.Marshal(changePlayerReq)
	if err != nil {
		return
	}

	changePlayer := models.WSMessage{
		Type: models.CHANGEPLAYER,
		Client: client,
		Payload: payload,
	}

	readCh <- changePlayer

}

func checkPlayerCash(cash, fine int) (available bool) {
	if cash >= fine {
		return true
	} else {
		return false
	}
}

func callCalculateRent(req *models.ReqCalculateRent, client *models.Client, readCh chan <- models.WSMessage) {
	reqPayload, err := json.Marshal(req)
	if err != nil {
		return
	}

	calRent := models.WSMessage{
		Type: models.CALCULATERENT,
		Client: client,
		Payload: reqPayload,
	}

	readCh <- calRent
}

func oldcallMovePos(dice int, client *models.Client, readCh chan <- models.WSMessage) {
		newPosReq := models.ReqMovePos{
			UpdateBy: dice,
		}

		newPosPayload, err := json.Marshal(newPosReq)
		if err != nil {
			return
		}

		movePos := models.WSMessage{
			Type:    models.MOVEPOS,
			Client:  client,
			Payload: newPosPayload,
		}

		readCh <- movePos

}

func callMovePos(resp []byte, client *gameroom.NewClient, param map[string]string) {
	var diceValResp models.RespDiceRoll
	json.Unmarshal(resp, &diceValResp)

	newPosReq := models.ReqMovePos{
		UpdateBy: diceValResp.DiceVal,
	}

	req, err := json.Marshal(newPosReq)
	if err != nil {
		logger.ZapLogger.Infow(models.ROLLDICE, "Json Error", err)
		return
	}

	err = client.Server.Write(models.MOVEPOS, req, param, nil)
	if err != nil {
		logger.ZapLogger.Errorw("WS Message Router", "Error", err)
		
	}
}

func getCardNo() (value int) {
	value = rand.IntN(6) + 1
	return

}


func callActionCard(client *models.Client, readCh chan <- models.WSMessage) {

}