package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)


func RollDice(request json.RawMessage, readCh chan <- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Roll Dice")
	var req models.Request

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}
	playerId := req.PlayerId

	diceVal := diceRoll()

	resp := models.RespDiceRoll{
		PlayerId: playerId,	
		GameId: req.GameId,
		DiceVal: diceVal,
	}
		
	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	go func()  {
		newPosReq := models.ReqMovePos{
			PlayerId: req.PlayerId,
			GameId: req.GameId,
			UpdateBy: diceVal,
		}

		newPosPayload, err := json.Marshal(newPosReq)
		if err != nil {
			logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
			return
		}

		movePos := models.WSMessage{
			Type: models.MOVEPOS,
			Payload: newPosPayload,
		}

		readCh <- movePos

	}()

	return
}


func MovePos(request json.RawMessage, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Info("Move Position")
	var req models.ReqMovePos

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "JSON Error", err)
		return
	}

	currPos, err := db.GetPlayerPos(req.PlayerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Current Position", currPos)

	newPos := updatePos(req.UpdateBy, currPos)
	logger.ZapLogger.Infow(models.MOVEPOS, "New Position", newPos)

	resp := models.RespMovePos{
		PlayerId: req.PlayerId,
		GameId: req.GameId,
		NewPos: newPos,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "JSON Error", err)
		return
	}

	return
}